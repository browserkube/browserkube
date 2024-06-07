package playwright

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/browserkube/browserkube/browserkube/internal/provision"
	browserkubev1 "github.com/browserkube/browserkube/operator/api/v1"
	browserkubehttp "github.com/browserkube/browserkube/pkg/http"
	"github.com/browserkube/browserkube/pkg/opentelemetry"
	"github.com/browserkube/browserkube/pkg/session"
	"github.com/browserkube/browserkube/pkg/sessionresult"
	revuuid "github.com/browserkube/browserkube/pkg/util/uuid"
	"github.com/browserkube/browserkube/storage"
)

var Module = fx.Options(
	fx.Provide(newPlaywrightProxy),
	fx.Invoke(initHandlers),
)

type playwrightProxy struct {
	logger             *zap.SugaredLogger
	manager            provision.Provisioner
	sessionRecord      bool
	provider           *sdktrace.TracerProvider
	sessionResultsRepo sessionresult.Repository
	sessionRecorder    storage.BlobSessionStorage
}

func newPlaywrightProxy(
	logger *zap.SugaredLogger,
	manager provision.Provisioner,
	sessionResultsRepo sessionresult.Repository,
	sessionRecorder storage.BlobSessionStorage,
) *playwrightProxy {
	provider, err := opentelemetry.InitProvider("playwrightProxy")
	if err != nil {
		logger.Error("Failed to initialize provider, error: ", err)
	}
	return &playwrightProxy{
		manager:            manager,
		logger:             logger,
		sessionRecord:      true,
		provider:           provider,
		sessionResultsRepo: sessionResultsRepo,
		sessionRecorder:    sessionRecorder,
	}
}

func (g *playwrightProxy) start(w http.ResponseWriter, rq *http.Request) error {
	// we use modified uuid version here
	// so the objects are sorted in Kubernetes/etcd in descending order
	uid := uuid.Must(revuuid.NewV7Reverse()).String()
	logger := g.logger.With("session", uid)
	browser := chi.URLParam(rq, "browser")
	if browser == "" {
		return browserkubehttp.NewHTTPErr(http.StatusBadRequest, errors.New("browser must be provided"))
	}

	var browserkubeOpts session.BrowserKubeOpts

	decoder := schema.NewDecoder()
	err := decoder.Decode(&browserkubeOpts, rq.URL.Query())
	if err != nil {
		logger.Error("fail to decode query")
		return browserkubehttp.NewHTTPErr(http.StatusBadRequest, err)
	}

	remote, err := g.manager.Provision(rq.Context(), uid, &session.Capabilities{
		Platform:    provision.PlatformLinux,
		BrowserName: browser,
		BrowserKubeOpts: session.BrowserKubeOpts{
			EnableVNC:        browserkubeOpts.EnableVNC,
			Type:             browserkubev1.TypePlaywright,
			Manual:           false,
			ScreenResolution: browserkubeOpts.ScreenResolution,
			EnableVideo:      browserkubeOpts.EnableVideo,
		},
	})
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		if dErr := g.manager.Delete(context.Background(), remote.Name); dErr != nil {
			logger.Errorf("unable to delete browser: %+v", dErr)
		}
	}()

	u := &url.URL{
		Scheme:   "ws",
		Host:     net.JoinHostPort(remote.Status.Host, remote.Status.PortConfig.Browser),
		RawQuery: rq.URL.RawQuery,
	}
	g.logger.Debugf("Proxying Playwright to %s", u.String())

	proxy, err := NewProxy(rq.Context(), u, uid, g.sessionRecorder)
	if err != nil {
		logger.Errorf("error while initializing proxy. err:%w", err)
		return err
	}
	proxy.ServeHTTP(w, rq)

	browserLogs, err := g.manager.Logs(rq.Context(), remote.Status.PodName, false)
	if err != nil {
		logger.Errorf("error while getting browser logs: %w", err)
	}
	// After we are done with this request, save the session...
	ctx, pCancel := context.WithTimeout(rq.Context(), 30*time.Second)
	defer pCancel()
	if g.sessionRecord {
		if err = proxy.SaveSessionRecord(ctx, uid); err != nil {
			logger.Errorf("error while save session record: %w", err)
			return err
		}
		if err = proxy.SaveScreenshotRecord(ctx, uid); err != nil {
			logger.Errorf("error while save screenshot record: %w", err)
			return err
		}
		if err = proxy.SaveBrowserLogRecord(ctx, uid, browserLogs); err != nil {
			logger.Errorf("error while save browser logs record: %w", err)
			return err
		}
	}

	if g.sessionResultsRepo != nil {
		err = proxy.SaveSessionResult(rq.Context(), uid, remote, g.sessionResultsRepo)
		if err != nil {
			logger.Error("error while trying to save the session result: %w", err)
		}
	}
	return nil
}

func initHandlers(mux chi.Router, pp *playwrightProxy) {
	mux.Group(func(r chi.Router) {
		if pp.provider != nil {
			r.Use(opentelemetry.HTTPMiddleware(pp.provider))
		} else {
			r.Use(opentelemetry.NewMetricsMiddleware("playwrightProxy"))
		}
		r.HandleFunc("/playwright/{browser}", browserkubehttp.Handler(pp.start))
	})
}
