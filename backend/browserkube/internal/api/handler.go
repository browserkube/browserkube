package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	netwebsocket "golang.org/x/net/websocket"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/browserkube/browserkube/browserkube/internal/provision"
	"github.com/browserkube/browserkube/browserkube/internal/snippet"
	browserkubev1 "github.com/browserkube/browserkube/operator/api/v1"
	browserkubehttp "github.com/browserkube/browserkube/pkg/http"
	"github.com/browserkube/browserkube/pkg/opentelemetry"
	"github.com/browserkube/browserkube/pkg/session"
	"github.com/browserkube/browserkube/pkg/sessionresult"
	browserkubeutil "github.com/browserkube/browserkube/pkg/util"
	"github.com/browserkube/browserkube/pkg/util/broadcast"
	"github.com/browserkube/browserkube/pkg/wd/wdproto"
	"github.com/browserkube/browserkube/storage"
)

const (
	defaultBatchFrameDuration = 3 * time.Second
	defaultPageSize           = 20
	keySessionID              = "sessionID"
	screenshotID              = "screenshotID"
	CommandsPath              = "/commands"
	ScreenshotsPath           = "/screenshots"
)

var Module = fx.Options(
	fx.Provide(newHandler),
	fx.Invoke(initRoutes),
)

func initRoutes(mux chi.Router, h *handler) {
	mux.Group(func(r chi.Router) {
		if h.provider != nil {
			r.Use(opentelemetry.HTTPMiddleware(h.provider))
		} else {
			r.Use(opentelemetry.NewMetricsMiddleware("api"))
		}

		r.Route("/sessions", func(r chi.Router) {
			r.Get("/", browserkubehttp.Handler(h.sessions))
			r.Delete("/{sessionID}", browserkubehttp.Handler(h.deleteSessionResult))

			r.Get("/{sessionID}/commands", browserkubehttp.Handler(h.getListCommands))

			r.Post("/{sessionID}/screenshots", browserkubehttp.Handler(h.createScreenshot))
			r.Get("/{sessionID}/screenshots/{screenshotID}", browserkubehttp.Handler(h.getScreenshotByID))

			r.Get("/{sessionID}/screenshots", browserkubehttp.Handler(h.getListScreenshots))
			r.Get("/{sessionID}/files/*", browserkubehttp.Handler(h.getSessionFile))
		})

		r.Route("/results", func(r chi.Router) {
			r.Get("/", browserkubehttp.Handler(h.results))
			r.Get("/{sessionID}", browserkubehttp.Handler(h.resultByID))
		})
		r.Get("/status", browserkubehttp.Handler(h.status))
		// Deprecated
		r.Get("/browsers", browserkubehttp.Handler(h.browsers))
		//-
		r.HandleFunc("/events", h.events)

		r.Handle("/logs/{sessionID}", h.logs())
		r.Handle("/vnc/{sessionID}", h.vnc())
		r.HandleFunc("/devtools/{sessionID}", h.reverseProxy(func(s *session.Session) string {
			return s.Browser.Status.PortConfig.DevTools
		}))
		r.HandleFunc("/download/{sessionID}", h.reverseProxy(func(s *session.Session) string {
			return s.Browser.Status.PortConfig.FileServer
		}))
		r.HandleFunc("/clipboard/{sessionID}", h.reverseProxy(func(s *session.Session) string {
			return s.Browser.Status.PortConfig.Clipboard
		}))

		r.Get("/snippet", browserkubehttp.Handler(h.snippet))
	})
}

type handler struct {
	sessionRepo           session.Repository
	sessionResultsRepo    sessionresult.Repository
	provisioner           provision.Provisioner
	upgrader              websocket.Upgrader
	logger                *zap.SugaredLogger
	startTime             time.Time
	provider              *sdktrace.TracerProvider
	sessionStorage        storage.BlobSessionStorage
	archiveSessionStorage storage.BlobSessionArchiveStorage
}

func newHandler(
	logger *zap.SugaredLogger,
	sessionRepo session.Repository,
	sessionResultsRepo sessionresult.Repository,
	provisioner provision.Provisioner,
	sessionStorage storage.BlobSessionStorage,
) *handler {
	provider, err := opentelemetry.InitProvider("api")
	if err != nil {
		logger.Error("Failed to initialize provider, error: ", err)
	}

	return &handler{
		logger:             logger,
		sessionRepo:        sessionRepo,
		sessionResultsRepo: sessionResultsRepo,
		provisioner:        provisioner,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		startTime:      time.Now(),
		provider:       provider,
		sessionStorage: sessionStorage,
	}
}

func (h *handler) events(w http.ResponseWriter, rq *http.Request) {
	ws, err := h.upgrader.Upgrade(w, rq, nil)
	if err != nil {
		var wsErr websocket.HandshakeError
		if errors.As(err, &wsErr) {
			h.logger.Errorf("Handshake error: %+v", wsErr)
		}
		return
	}

	var barchFrame time.Duration
	batchFrameStr := rq.URL.Query().Get("batch")
	if batchFrameStr == "" {
		barchFrame = defaultBatchFrameDuration
	} else {
		barchFrame, err = time.ParseDuration(batchFrameStr)
		if err != nil || barchFrame <= 0 {
			var msg string
			if err != nil {
				msg = err.Error()
			} else {
				msg = "incorrect duration interval"
			}
			if wErr := browserkubehttp.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": msg}); wErr != nil {
				h.logger.Error(wErr)
			}
		}
	}

	go func() {
		ctx, cancelFunc := context.WithCancel(context.Background())
		defer cancelFunc()
		defer browserkubeutil.CloseQuietly(ws)

		sessions := h.sessionRepo.Watch(ctx)
		batcher := broadcast.NewBatcher[*session.Session](barchFrame, false)

		bErr := batcher.Batch(ctx, sessions, func(batch []*session.Session) error {
			// Events
			deduplicated := browserkubeutil.ReverseDeduplicateBY[*session.Session, string](batch, func(s *session.Session) string {
				return s.ID
			})
			sess := browserkubeutil.Map[*session.Session, *Session](deduplicated, h.toSession)
			wErr := ws.WriteJSON(NewWSMessage("session", sess))
			if wErr != nil {
				return errors.WithStack(wErr)
			}

			// Stats
			return h.wsStatus(ws)
		})
		if bErr != nil {
			h.logger.Error(bErr)
		}
	}()
}

// sessions godoc
//
//	@Summary		sessions
//	@Description	get all sessions
//	@Tags			browsers
//	@Produce		json
//	@Success		200	{object}	[]Session
//	@Failure		500	{string}	Internal	Server	Error
//	@Router			/sessions [get]
func (h *handler) sessions(w http.ResponseWriter, _ *http.Request) error {
	sessions, err := h.sessionRepo.FindAll()
	if err != nil {
		return browserkubehttp.NewHTTPErr(http.StatusInternalServerError, errors.WithStack(err))
	}
	sort.Slice(sessions, func(i, j int) bool {
		ts1 := sessions[i].Browser.CreationTimestamp
		ts2 := sessions[j].Browser.CreationTimestamp
		return ts1.After(ts2.Time)
	})

	bffSess := browserkubeutil.Map[*session.Session, *Session](sessions, h.toSession)
	return errors.WithStack(browserkubehttp.WriteJSON(w, http.StatusOK, bffSess))
}

// results godoc
//
//	@Summary		results
//	@Description	get results of sessions
//	@Tags			browsers
//	@Produce		json
//	@Param			pageSize	path		string	true	"page size"
//	@Param			pageToken	path		string	true	"page token"
//	@Success		200			{object}	browserkubeutil.Page[SessionResult]
//	@Failure		400			{string}	Bad			request
//	@Failure		500			{string}	Internal	Server	Error
//	@Router			/results [get]
func (h *handler) results(w http.ResponseWriter, rq *http.Request) error {
	var (
		pageSize  int
		pageToken string
		err       error
	)
	pageSizeStr := chi.URLParam(rq, "pageSize")
	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			return browserkubehttp.NewHTTPErr(http.StatusBadRequest, errors.Wrap(err, "incorrect limit format"))
		}
	} else {
		pageSize = defaultPageSize
	}

	pageToken = chi.URLParam(rq, "pageToken")

	sessResults, err := h.sessionResultsRepo.FindAll(rq.Context(), pageSize, pageToken)
	if err != nil {
		return browserkubehttp.NewHTTPErr(http.StatusInternalServerError, errors.WithStack(err))
	}

	results, err := browserkubeutil.MapPageErr[*sessionresult.Result, *SessionResult](sessResults, h.toSessionResult)
	if err != nil {
		return browserkubehttp.NewHTTPErr(http.StatusInternalServerError, errors.WithStack(err))
	}
	return errors.WithStack(browserkubehttp.WriteJSON(w, http.StatusOK, results))
}

// resultByID godoc
//
//	@Summary		resultByID
//	@Description	get result of sessions by ID
//	@Tags			browsers
//	@Produce		json
//	@Param			sessionID	path		string	true	"sessionID"
//	@Success		200			{object}	SessionResult
//	@Failure		400			{string}	Bad	request
//	@Failure		404			{string}	NotFound
//	@Failure		500			{string}	Internal	Server	Error
//	@Router			/results/{sessionID} [get]
func (h *handler) resultByID(w http.ResponseWriter, rq *http.Request) error {
	sessionID := chi.URLParam(rq, keySessionID)
	if sessionID == "" {
		msg := "session id not found"
		h.logger.Error(msg)
		return browserkubehttp.NewHTTPErr(http.StatusBadRequest, errors.New(msg))
	}

	sessRes, err := h.sessionResultsRepo.FindByID(rq.Context(), sessionID)
	if err != nil {
		return browserkubehttp.NewHTTPErr(http.StatusNotFound, errors.WithStack(err))
	}
	result, err := h.toSessionResult(sessRes)
	if err != nil {
		return browserkubehttp.NewHTTPErr(http.StatusInternalServerError, errors.Wrap(err, "Incorrect session result"))
	}
	return errors.WithStack(browserkubehttp.WriteJSON(w, http.StatusOK, result))
}

// status godoc
//
//	@Summary		status
//	@Description	get status of sessions
//	@Tags			browsers
//	@Produce		json
//	@Success		200	{object}	Status
//	@Failure		500	{string}	Internal	Server	Error
//	@Router			/status [get]
func (h *handler) status(w http.ResponseWriter, _ *http.Request) error {
	qCurrent, qMax, err := h.sessionRepo.Quota()
	if err != nil {
		return browserkubehttp.NewHTTPErr(http.StatusInternalServerError, errors.WithStack(err))
	}
	var qConnecting, qQueued, qRunning int
	sessions, err := h.sessionRepo.FindAll()
	if err != nil {
		return browserkubehttp.NewHTTPErr(http.StatusInternalServerError, errors.WithStack(err))
	}
	for _, sess := range sessions {
		switch sess.State {
		case "pending":
			if qCurrent > qMax {
				qQueued++
				break
			}
			qConnecting++
		case "running":
			qRunning++
		}
	}

	return errors.WithStack(browserkubehttp.WriteJSON(w, http.StatusOK,
		&Status{
			QuotesLimit: qMax,
			MaxTimeout:  time.Minute,
			Stats:       StatusStats{All: qCurrent, Connecting: qConnecting, Queued: qQueued, Running: qRunning},
		}))
}

// browsers godoc
//
//	@Summary		listBrowsers
//	@Description	get a list of browsers
//	@Tags			browsers
//	@Produce		json
//	@Param			manual	query		boolean	false	"manual session"
//	@Success		200		{array}		Browser
//	@Failure		400		{string}	Bad			request
//	@Failure		500		{string}	Internal	Server	Error
//	@Router			/browsers [get]
func (h *handler) browsers(w http.ResponseWriter, rq *http.Request) error {
	manualOnly, _ := strconv.ParseBool(rq.URL.Query().Get("manual"))

	var mapping browserkubev1.BrowserSet
	mappingList, err := h.provisioner.Available(rq.Context())
	if err != nil {
		return errors.WithStack(err)
	}
	if len(mappingList.Items) > 0 {
		mapping = mappingList.Items[0]
	}

	var browsers []Browser
	for browserName, versions := range mapping.Spec.WebDriver {
		for version, browserConfig := range versions.Versions {
			browsers = append(browsers, Browser{
				Name:        browserName,
				Platform:    provision.PlatformLinux,
				Version:     version,
				Image:       browserConfig.Image,
				Type:        browserkubev1.TypeWebDriver,
				Resolutions: defaultResolutions,
			})
		}
	}

	if !manualOnly {
		for browserName, versions := range mapping.Spec.Playwright {
			for version, browserConfig := range versions.Versions {
				browsers = append(browsers, Browser{
					Name:        browserName,
					Platform:    provision.PlatformLinux,
					Version:     version,
					Image:       browserConfig.Image,
					Type:        browserkubev1.TypePlaywright,
					Resolutions: defaultResolutions,
				})
			}
		}
	}
	h.sortBrowsers(browsers)

	return errors.WithStack(browserkubehttp.WriteJSON(w, http.StatusOK, browsers))
}

// deleteSessionResult godoc
//
//	@Summary		deleteSessionResult
//	@Description	delete session result
//	@Tags			browsers
//	@Produce		json
//	@Param			sessionID	path		string	true	"sessionID"
//	@Success		204			{string}	ok
//	@Failure		400			{string}	Bad			request
//	@Failure		500			{string}	Internal	Server	Error
//	@Router			/sessions/{sessionID} [delete]
func (h *handler) deleteSessionResult(w http.ResponseWriter, rq *http.Request) error {
	logger := h.logger.With("session", rq.URL.String())

	sessionID := chi.URLParam(rq, keySessionID)
	if sessionID == "" {
		logger.Error("session id not found")
		return browserkubehttp.NewHTTPErr(http.StatusBadRequest, fmt.Errorf("session id not found"))
	}

	if err := h.sessionResultsRepo.Delete(rq.Context(), sessionID, metav1.DeleteOptions{}); err != nil {
		logger.Error("Session has not been deleted: ", err)
		return browserkubehttp.NewHTTPErr(http.StatusInternalServerError, errors.WithStack(err))
	}

	logger.Info("Session has been deleted")
	return errors.WithStack(browserkubehttp.WriteJSON(w, http.StatusNoContent, nil))
}

// createScreenshot godoc
//
//	@Summary		createScreenshot
//	@Description	create a screenshot
//	@Tags			browsers
//	@Param			sessionID	path		string	true	"sessionID"
//	@Success		200			{string}	ok
//	@Failure		400			{string}	Bad	request
//	@Failure		404			{string}	NotFound
//	@Failure		500			{string}	Internal	Server	Error
//	@Router			/sessions/{sessionID}/screenshots [post]
func (h *handler) createScreenshot(w http.ResponseWriter, rq *http.Request) error {
	sessionID := chi.URLParam(rq, keySessionID)
	if sessionID == "" {
		h.logger.Error("session id not found")
		return browserkubehttp.NewHTTPErr(http.StatusBadRequest, fmt.Errorf("session id not found"))
	}

	logger := h.logger.With("session_id", sessionID)

	sess, err := h.sessionRepo.FindByID(sessionID)
	if err != nil {
		logger.Errorf("unable to find session: %v", err)
		return browserkubehttp.NewHTTPErr(http.StatusNotFound, errors.Wrap(err, "unable to find session"))
	}

	ctx := rq.Context()

	screenshotBytes, err := wdproto.NewWebDriver(sess.Browser.Status.SeleniumURL, sess.ID).TakeScreenshot(ctx)
	if err != nil {
		logger.Errorf("unable to take a screenshot: %s", err)
		return browserkubehttp.NewHTTPErr(http.StatusInternalServerError, errors.WithStack(err))
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(screenshotBytes); err != nil {
		logger.Errorf("failed to encode screenshotBytes: %v", err)
		return browserkubehttp.NewHTTPErr(http.StatusInternalServerError, errors.WithStack(err))
	}

	fileName := time.Now().UTC().Format("2006-01-02-15-04-05") + "-screenshot.png"

	if err := h.sessionStorage.SaveFile(ctx, sess.ID, ScreenshotsPath, &storage.BlobFile{
		FileName:    fileName,
		ContentType: "image/png",
		Content:     &buf,
	}); err != nil {
		logger.Errorf("failed to save screenshot to storage: %v", err)
		return browserkubehttp.NewHTTPErr(http.StatusInternalServerError, errors.WithStack(err))
	}

	return errors.WithStack(browserkubehttp.WriteJSON(w, http.StatusOK, fileName))
}

// getScreenshotByID godoc
//
//	@Summary		get screenshot by ID
//	@Description	get a screenshot
//	@Tags			browsers
//	@Param			sessionID		path		string	true	"sessionID"
//	@Param			screenshotID	path		string	true	"screenshotID"
//	@Success		200				{string}	string	ok
//	@Failure		400				{string}	Bad		request
//	@Failure		404				{string}	NotFound
//	@Failure		500				{string}	Internal	Server	Error
//	@Router			/sessions/{sessionID}/screenshots/{screenshotID} [get]
func (h *handler) getScreenshotByID(w http.ResponseWriter, rq *http.Request) error {
	sessionID := chi.URLParam(rq, keySessionID)
	if sessionID == "" {
		h.logger.Error("session id not found")
		return browserkubehttp.NewHTTPErr(http.StatusBadRequest, fmt.Errorf("session id not found"))
	}

	logger := h.logger.With("session_id", sessionID)

	fileName := chi.URLParam(rq, screenshotID)
	if sessionID == "" {
		logger.Error("screenshot name not found")
		return browserkubehttp.NewHTTPErr(http.StatusBadRequest, fmt.Errorf("screenshot name not found"))
	}

	sess, err := h.sessionRepo.FindByID(sessionID)
	if err != nil {
		logger.Errorf("unable to find session: %v", err)
		return browserkubehttp.NewHTTPErr(http.StatusNotFound, errors.Wrap(err, "unable to find session"))
	}

	screenshotsPath := filepath.Join(sess.ID, ScreenshotsPath)

	screenshot, err := h.sessionStorage.GetFile(rq.Context(), screenshotsPath, fileName)
	if err != nil {
		logger.Errorf("failed to get screenshot: %v", err)
		return browserkubehttp.NewHTTPErr(http.StatusInternalServerError, errors.WithStack(err))
	}

	payload := &bytes.Buffer{}
	if _, err = io.Copy(payload, screenshot.Content); err != nil {
		logger.Errorf("failed to copy screenshot.Content: %v", err)
		return browserkubehttp.NewHTTPErr(http.StatusInternalServerError, errors.WithStack(err))
	}

	return errors.WithStack(browserkubehttp.WriteJSON(w, http.StatusOK, payload.Bytes()))
}

// getListScreenshots godoc
//
//	@Summary	get list screenshots
//	@Tags		browsers
//	@Param		sessionID	path		string	true	"sessionID"
//	@Success	200			{array}		string
//	@Failure	400			{string}	Bad	request
//	@Failure	404			{string}	NotFound
//	@Failure	500			{string}	Internal	Server	Error
//	@Router		/sessions/{sessionID}/screenshots [get]
func (h *handler) getListScreenshots(w http.ResponseWriter, rq *http.Request) error {
	sessionID := chi.URLParam(rq, keySessionID)
	if sessionID == "" {
		return browserkubehttp.NewHTTPErr(http.StatusBadRequest, fmt.Errorf("session id is required"))
	}

	logger := h.logger.With("session_id", sessionID)

	fileNames, err := h.sessionStorage.ListFileNames(rq.Context(), sessionID, ScreenshotsPath)
	if err != nil {
		logger.Errorf("failed to get list of fileNames for screenshots: %v", err)
		return browserkubehttp.NewHTTPErr(http.StatusInternalServerError, errors.WithStack(err))
	}

	scr := struct {
		Screenshots [][]byte `json:"screenshots"`
	}{}

	for _, fileName := range fileNames {
		screenshot, err := h.sessionStorage.GetFile(rq.Context(), sessionID, fileName)
		if err != nil {
			logger.Errorf("failed to get screenshot: %v", err)
			return browserkubehttp.NewHTTPErr(http.StatusInternalServerError, errors.WithStack(err))
		}

		payload := &bytes.Buffer{}
		if _, err = io.Copy(payload, screenshot.Content); err != nil {
			logger.Errorf("failed to copy screenshot.Content: %v", err)
			return browserkubehttp.NewHTTPErr(http.StatusInternalServerError, errors.WithStack(err))
		}

		scr.Screenshots = append(scr.Screenshots, payload.Bytes())
	}
	// return an empty array instead of nil
	if len(scr.Screenshots) == 0 {
		scr.Screenshots = make([][]byte, 0)
	}

	return errors.WithStack(browserkubehttp.WriteJSON(w, http.StatusOK, scr))
}

// getListCommands godoc
//
//	@Summary	get list commands
//	@Tags		browsers
//	@Param		sessionID	path		string	true	"sessionID"
//	@Param		pageToken	query		string	true	"for the first request always indicate - first page"
//	@Param		pageSize	query		int		true	"pageSize"
//	@Success	200			{object}	CommandLogResponse
//	@Failure	400			{string}	Bad	request
//	@Failure	404			{string}	NotFound
//	@Failure	500			{string}	Internal	Server	Error
//	@Router		/sessions/{sessionID}/commands [get]
func (h *handler) getListCommands(w http.ResponseWriter, rq *http.Request) error {
	sessionID := chi.URLParam(rq, keySessionID)
	if sessionID == "" {
		h.logger.Error("session id not found")
		return browserkubehttp.NewHTTPErr(http.StatusBadRequest, fmt.Errorf("session id not found"))
	}

	logger := h.logger.With("session_id", sessionID)

	qParams := rq.URL.Query()
	pageToken := qParams.Get("pageToken")

	pageSize, err := strconv.Atoi(qParams.Get("pageSize"))
	if err != nil {
		logger.Errorf("unable to convert qParams pageSize: %v", err)
		return browserkubehttp.NewHTTPErr(http.StatusBadRequest, errors.Wrap(err, "unable to convert pageSize: provide the correct value"))
	}

	sessResult, err := h.sessionResultsRepo.FindByID(rq.Context(), sessionID)
	if err != nil {
		logger.Errorf("unable to find session: %v", err)
		return browserkubehttp.NewHTTPErr(http.StatusNotFound, errors.Wrap(err, "unable to find session"))
	}

	fileNames, newPageToken, err := h.sessionStorage.ListPage(rq.Context(), sessResult.Name, CommandsPath, pageToken, pageSize)
	if err != nil {
		logger.Errorf("unable to get list of command files: %v", err)
		return browserkubehttp.NewHTTPErr(http.StatusNotFound, errors.WithStack(err))
	}

	commandLogResponse := CommandLogResponse{
		NewPageToken: string(newPageToken),
	}

	var command CommandLog

	for _, fileName := range fileNames {
		commandRecord, err := h.sessionStorage.GetFile(rq.Context(), sessResult.Name, fileName)
		if err != nil {
			logger.Errorf("unable to get command record: %v", err)
			return browserkubehttp.NewHTTPErr(http.StatusInternalServerError, errors.WithStack(err))
		}

		if err := json.NewDecoder(commandRecord.Content).Decode(&command); err != nil {
			logger.Error("unable to decode command record: ", err)
			return browserkubehttp.NewHTTPErr(http.StatusInternalServerError, errors.WithStack(err))
		}

		commandLogResponse.Commands = append(commandLogResponse.Commands, command)
	}

	return errors.WithStack(browserkubehttp.WriteJSON(w, http.StatusOK, commandLogResponse))
}

// getSessionFile godoc
//
//	@Summary	get session file commands
//	@Tags		browsers
//	@Param		sessionID	path		string	true	"sessionID"
//	@Param		filePath	path		string	true	"filePath"
//	@Success	200			{string}	string	ok
//	@Failure	400			{string}	Bad	request
//	@Failure	404			{string}	NotFound
//	@Failure	500			{string}	Internal	Server	Error
//	@Router		/sessions/{sessionID}/files/{filePath} [get]
func (h *handler) getSessionFile(w http.ResponseWriter, rq *http.Request) error {
	fPath := chi.URLParam(rq, "*")
	sessionID := chi.URLParam(rq, "sessionID")
	file, err := h.sessionStorage.GetFile(rq.Context(), sessionID, fPath)
	if err != nil {
		return browserkubehttp.NewHTTPErr(http.StatusInternalServerError, errors.WithStack(err))
	}
	w.Header().Set("Content-Type", file.ContentType)
	if _, err = io.Copy(w, file.Content); err != nil {
		zap.S().Errorf("failed to copy screenshot.Content: %v", err)
		return browserkubehttp.NewHTTPErr(http.StatusInternalServerError, errors.WithStack(err))
	}
	return nil
}

func (h *handler) logs() netwebsocket.Handler {
	return func(wsconn *netwebsocket.Conn) {
		defer browserkubeutil.CloseQuietly(wsconn)

		logger := h.logger.With(
			"request", fmt.Sprintf("%s %s", wsconn.Request().Method, wsconn.Request().URL.Path),
		)

		sessionID := chi.URLParam(wsconn.Request(), keySessionID)
		if sessionID == "" {
			logger.Error("session id not found")
			return
		}

		logger = logger.With("session_id", sessionID)

		sess, err := h.sessionRepo.FindByID(sessionID)
		if err != nil {
			logger.Errorf("unable to find session: %v", err)
			return
		}

		logs, err := h.provisioner.Logs(wsconn.Request().Context(), sess.Browser.Status.PodName, true)
		if err != nil {
			logger.Errorf("stream logs error: %v", err)
			return
		}
		defer browserkubeutil.CloseQuietly(logs)

		wsconn.PayloadType = netwebsocket.BinaryFrame
		_, err = io.Copy(wsconn, logs)
		logger.With("error", err).Infof("stream logs disconnected")
	}
}

func (h *handler) vnc() netwebsocket.Handler {
	return func(wsconn *netwebsocket.Conn) {
		defer browserkubeutil.CloseQuietly(wsconn)

		sessionID := chi.URLParam(wsconn.Request(), keySessionID)
		if sessionID == "" {
			h.logger.With("request",
				fmt.Sprintf("%s %s", wsconn.Request().Method, wsconn.Request().URL.Path)).Error("can't parse session ID")
			return
		}

		sess, err := h.sessionRepo.FindByID(sessionID)
		if err != nil {
			h.logger.With("request",
				fmt.Sprintf("%s %s", wsconn.Request().Method, wsconn.Request().URL.Path)).Error("session ID not found")
			return
		}

		// host := browserkubeutil.BuildHostPort(sessionID, "browserkube", strconv.Itoa(session.Remote.PortConfig.VNCPort))
		host := net.JoinHostPort(sess.Browser.Status.Host, sess.Browser.Status.PortConfig.VNC)
		logger := h.logger.With(
			"request_id", uuid.New(),
			"session_id", sessionID,
			"request", fmt.Sprintf("%s %s", wsconn.Request().Method, wsconn.Request().URL.Path),
		)
		logger.Infof("vnc request: %s", host)

		var dialer net.Dialer
		conn, err := dialer.DialContext(wsconn.Request().Context(), "tcp", host)
		if err != nil {
			logger.Errorf("vnc connection error: %v", err)
			return
		}
		defer browserkubeutil.CloseQuietly(conn)

		wsconn.PayloadType = netwebsocket.BinaryFrame
		go func() {
			if _, err := io.Copy(wsconn, conn); err != nil {
				logger.Errorf("ws connection error: %s", err)
			}
			logger.Warnf("vnc connection closed")
		}()
		if _, err := io.Copy(conn, wsconn); err != nil {
			logger.Errorf("ws connection error: %s", err)
		}
		logger.Infof("vnc client disconnected")
	}
}

func (h *handler) ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(struct {
		Uptime         string `json:"uptime"`
		LastReloadTime string `json:"lastReloadTime"`
		NumRequests    uint64 `json:"numRequests"`
		Version        string `json:"version"`
	}{time.Since(h.startTime).String(), h.startTime.Format(time.RFC3339), 0, "0.0.1"})
	if err != nil {
		h.logger.Error(err)
	}
}

func (h *handler) reverseProxy(portF func(s *session.Session) string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := chi.URLParam(r, keySessionID)
		logger := h.logger.With("session", sessionID)

		// load the session
		sess, err := h.sessionRepo.FindByID(sessionID)
		if err != nil || sess == nil {
			logger.Error("unable to find session")
			r.URL.Path = "/error"
			return
		}

		(&httputil.ReverseProxy{
			Director: func(rq *http.Request) {
				scheme := browserkubeutil.FirstNonEmpty(rq.URL.Scheme, "http")
				rq.URL.Scheme = scheme
				rq.URL.Host = net.JoinHostPort(sess.Browser.Status.Host, portF(sess))
				_, rqPath, _ := strings.Cut(rq.URL.Path, sessionID)
				rq.URL.Path = rqPath
				logger.Info("Proxying request to ", rq.URL.String())
			},
			ErrorHandler: func(w http.ResponseWriter, rq *http.Request, err error) {
				w.WriteHeader(http.StatusBadGateway)
			},
		}).ServeHTTP(w, r)
	}
}

func (h *handler) toSession(sess *session.Session) *Session {
	return &Session{
		ID:               sess.ID,
		Name:             sess.Caps.BrowserKubeOpts.Name,
		State:            sess.State,
		Platform:         provision.PlatformLinux,
		Manual:           sess.Caps.BrowserKubeOpts.Manual,
		ScreenResolution: sess.Caps.BrowserKubeOpts.ScreenResolution,
		Browser:          sess.Browser.Spec.BrowserName,
		BrowserVersion:   sess.Browser.Spec.BrowserVersion,
		Image:            sess.Browser.Status.Image,
		LogsOn:           true,
		CreatedAt:        Timestamp(sess.Browser.CreationTimestamp.Time),
		Type:             sess.Caps.BrowserKubeOpts.Type,
		VideoRecOn:       sess.Caps.BrowserKubeOpts.EnableVideo,
		VncOn:            sess.Caps.BrowserKubeOpts.EnableVNC,
		VncPsw:           sess.Browser.Status.VncPass,
	}
}

func (h *handler) toSessionResult(sess *sessionresult.Result) (*SessionResult, error) {
	caps, err := sess.GetCaps()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	sr := &SessionResult{
		Session: Session{
			State:            "Terminated",
			ID:               sess.Name,
			Name:             caps.BrowserKubeOpts.Name,
			Platform:         provision.PlatformLinux,
			Manual:           caps.BrowserKubeOpts.Manual,
			ScreenResolution: caps.BrowserKubeOpts.ScreenResolution,
			Browser:          sess.Spec.Browser.BrowserName,
			BrowserVersion:   sess.Spec.Browser.BrowserVersion,
			Image:            sess.Spec.BrowserImage,
			VncOn:            sess.Spec.Browser.EnableVNC,
			LogsOn:           true,
			LogsRefAddr:      path.Join(sess.Name, sessionresult.BrowserLogFileName),
			CreatedAt:        Timestamp(sess.Spec.StartedAt.Time),
			Type:             caps.BrowserKubeOpts.Type,
		},
	}

	if caps.BrowserKubeOpts.EnableVideo {
		sr.Session.VideoRefAddr = path.Join(sess.Name, sessionresult.VideoFileName)
	}

	return sr, nil
}

func (h *handler) wsStatus(ws *websocket.Conn) error {
	qCurrent, qMax, err := h.sessionRepo.Quota()
	if err != nil {
		return errors.WithStack(err)
	}
	var qConnecting, qQueued, qRunning int
	sessions, err := h.sessionRepo.FindAll()
	if err != nil {
		return errors.WithStack(err)
	}
	for _, sess := range sessions {
		switch sess.State {
		case "pending":
			if qCurrent > qMax {
				qQueued++
				break
			}
			qConnecting++
		case "running":
			qRunning++
		}
	}
	wErr := ws.WriteJSON(NewWSMessage("status", &Status{
		QuotesLimit: qMax,
		MaxTimeout:  time.Minute,
		Stats:       StatusStats{All: qCurrent, Connecting: qConnecting, Queued: qQueued, Running: qRunning},
	}))
	return errors.WithStack(wErr)
}

func (h *handler) sortBrowsers(browsers []Browser) {
	slices.SortFunc(browsers, func(a, b Browser) int {
		typeCompare := strings.Compare(a.Type, b.Type)
		// sort by type desc
		if typeCompare != 0 {
			return -1 * typeCompare
		}
		platformCmp := strings.Compare(a.Name, b.Name)
		if platformCmp != 0 {
			return platformCmp
		}
		return -1 * strings.Compare(a.Version, b.Version)
	})
}

// snippet godoc
//
//	@Summary		getSnippet
//	@Description	get code snippet
//	@Tags			browsers
//	@Produce		json
//	@Param			manual	query		boolean	false	"manual session"
//	@Success		200		{array}		Browser
//	@Failure		400		{string}	Bad			request
//	@Failure		500		{string}	Internal	Server	Error
//	@Router			/snippet [get]
func (h *handler) snippet(w http.ResponseWriter, rq *http.Request) error {
	qParams := rq.URL.Query()

	tp := qParams.Get("type")
	lang := qParams.Get("language")
	browserName := qParams.Get("browserName")
	browserVersion := qParams.Get("browserVersion")

	snippet, err := snippet.GetSnippet(tp, lang, browserName, browserVersion)
	if err != nil {
		return browserkubehttp.NewHTTPErr(http.StatusBadRequest, errors.Wrap(err, "unsupported"))
	}

	return errors.WithStack(browserkubehttp.WritePlainText(w, http.StatusOK, snippet))
}
