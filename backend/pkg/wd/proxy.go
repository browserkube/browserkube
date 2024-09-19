package wd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"dario.cat/mergo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	v1 "github.com/browserkube/browserkube/operator/api/v1"
	"github.com/browserkube/browserkube/pkg/session"
	browserkubeutil "github.com/browserkube/browserkube/pkg/util"
	revuuid "github.com/browserkube/browserkube/pkg/util/uuid"
	"github.com/browserkube/browserkube/pkg/wd/wdproto"
	"github.com/browserkube/browserkube/pkg/websocketproxy"
)

const (
	downloadsBasePath = "/downloads/"
	keySessionID      = "sessionID"
)

type Context struct {
	context.Context //nolint: containedctx
}

func (c *Context) WithContext(ctx context.Context) *Context {
	c.Context = ctx
	return c
}

func (c *Context) WithValue(key, val any) *Context {
	c.Context = context.WithValue(c.Context, key, val)
	return c
}

type (
	PluginOpt  func(*ProxyBuilder)
	PluginOpts struct {
		// Weight is a plugin initialization priority parameter.
		// Higher weight means earlier initialization. Valid range is from 0 to 255
		// Note: two plugins with equal weight may be in uncertain order.
		Weight uint8
		Opts   []PluginOpt
	}
)

type (
	OnBeforeSessionStart = func(*Context, *httputil.ProxyRequest, *wdproto.NewSessionRQ, string) error
	OnAfterSessionStart  = func(*Context, *http.Response, string) error
	OnBeforeCommand      = func(*Context, *httputil.ProxyRequest, *session.Session) error
	OnAfterCommand       = func(*Context, *http.Response, *session.Session, string) error
	OnSessionQuit        = func(*Context, *session.Session) error
)

func WithBeforeSessionCreated(f func(hook OnBeforeSessionStart) OnBeforeSessionStart) PluginOpt {
	return func(p *ProxyBuilder) {
		p.beforeSessionHooks = append(p.beforeSessionHooks, f)
	}
}

func WithAfterSessionCreated(f func(OnAfterSessionStart) OnAfterSessionStart) PluginOpt {
	return func(p *ProxyBuilder) {
		p.afterSessionHooks = append(p.afterSessionHooks, f) //nolint:bodyclose
	}
}

func WithQuitSession(f func(OnSessionQuit) OnSessionQuit) PluginOpt {
	return func(p *ProxyBuilder) {
		p.quitSessionHooks = append(p.quitSessionHooks, f)
	}
}

func WithBeforeCommand(f func(OnBeforeCommand) OnBeforeCommand) PluginOpt {
	return func(p *ProxyBuilder) {
		p.beforeCommandHooks = append(p.beforeCommandHooks, f)
		p.beforeCommandHooks = append(p.beforeCommandHooks, func(cmd OnBeforeCommand) OnBeforeCommand {
			return func(ctx *Context, rq *httputil.ProxyRequest, s *session.Session) error {
				adjustUploadPath(rq)
				return nil
			}
		})
	}
}

func WithAfterCommand(f func(OnAfterCommand) OnAfterCommand) PluginOpt {
	return func(p *ProxyBuilder) {
		p.afterCommandHooks = append(p.afterCommandHooks, f) //nolint:bodyclose
	}
}

type ProxyBuilder struct {
	beforeSessionHooks []func(OnBeforeSessionStart) OnBeforeSessionStart
	afterSessionHooks  []func(OnAfterSessionStart) OnAfterSessionStart
	beforeCommandHooks []func(OnBeforeCommand) OnBeforeCommand
	afterCommandHooks  []func(OnAfterCommand) OnAfterCommand
	quitSessionHooks   []func(quit OnSessionQuit) OnSessionQuit
}

func (pb *ProxyBuilder) Build(sessionRepo session.Repository) *ProxyManager {
	dummyOnBeforeSession := func(ctx *Context, request *httputil.ProxyRequest, rq *wdproto.NewSessionRQ, sessionID string) error {
		return nil
	}
	dummyOnAfterSession := func(ctx *Context, response *http.Response, sessionID string) error {
		return nil
	}
	dummyOnBeforeCommand := func(ctx *Context, rq *httputil.ProxyRequest, sess *session.Session) error {
		return nil
	}
	dummyOnAfterCommand := func(ctx *Context, rs *http.Response, sess *session.Session, command string) error {
		return nil
	}
	dummyOnQuitCommand := func(ctx *Context, sess *session.Session) error {
		return nil
	}

	return &ProxyManager{
		sessionRepo:       sessionRepo,
		beforeSessionHook: chain[OnBeforeSessionStart](pb.beforeSessionHooks, dummyOnBeforeSession),
		afterSessionHook:  chain[OnAfterSessionStart](pb.afterSessionHooks, dummyOnAfterSession), //nolint:bodyclose
		beforeCommandHook: chain[OnBeforeCommand](pb.beforeCommandHooks, dummyOnBeforeCommand),
		afterCommandHook:  chain[OnAfterCommand](pb.afterCommandHooks, dummyOnAfterCommand), //nolint:bodyclose
		quitSessionHook:   chain[OnSessionQuit](pb.quitSessionHooks, dummyOnQuitCommand),
		log:               zap.S(),
	}
}

func NewProxyBuilder(opts ...PluginOpt) *ProxyBuilder {
	p := &ProxyBuilder{}
	for _, o := range opts {
		o(p)
	}
	return p
}

type ProxyManager struct {
	beforeSessionHook OnBeforeSessionStart
	afterSessionHook  OnAfterSessionStart
	beforeCommandHook OnBeforeCommand
	afterCommandHook  OnAfterCommand
	quitSessionHook   OnSessionQuit

	sessionRepo session.Repository
	log         *zap.SugaredLogger
}

// CreateWDSession godoc
//
//	@Summary		createWDSession
//	@Description	create webdriver session
//	@Tags			browsers
//	@Accept			json
//	@Produce		json
//	@Param			request	body		wdproto.CreateBrowserRequest	true	"browser request"
//	@Success		200		{object}	v1.Browser
//	@Failure		400		{string}	Bad	request
//	@Failure		502		{string}	Bad	Gateway	Error
//	@Router			/api/browsers [post]
//
// CreateWDSession Wraps the StartSessionHandler to change the input from manual creation request to generic selenium creation request.
func (p *ProxyManager) CreateWDSession(w http.ResponseWriter, rq *http.Request) {
	req := &wdproto.CreateBrowserRequest{}
	if err := json.NewDecoder(rq.Body).Decode(req); err != nil {
		wdproto.BadGatewayError(w, err)
		return
	}
	browserkubeutil.CloseQuietly(rq.Body)
	newSessionRQ := &wdproto.NewSessionRQ{
		W3CCapabilities: wdproto.W3CCapabilities{
			Capabilities: session.Capabilities{
				BrowserName:    req.BrowserName,
				BrowserVersion: req.BrowserVersion,
				BrowserKubeOpts: session.BrowserKubeOpts{
					Type:             v1.TypeWebDriver,
					EnableVideo:      req.RecordVideo,
					EnableVNC:        true,
					Name:             req.SessionName,
					Manual:           true,
					ScreenResolution: req.Resolution,
				},
			},
		},
	}
	modified := &bytes.Buffer{}
	err := json.NewEncoder(modified).Encode(newSessionRQ)
	if err != nil {
		p.log.Errorf("error while marshaling. err: %w", err)
		wdproto.BadGatewayError(w, err)
		return
	}
	rq.Body = io.NopCloser(modified)

	rq.ContentLength = int64(modified.Len())
	rq.URL.Path = "/wd/hub/session"

	p.StartSessionHandler(w, rq)
}

// DeleteWDSession godoc
//
//	@Summary		deleteWDSession
//	@Description	delete webdriver session
//	@Tags			browsers
//	@Param			sessionID	path		string	true	"session ID"
//	@Success		200			{string}	ok
//	@Failure		502			{string}	Bad	Gateway	Error
//	@Router			/api/browsers/{sessionID} [delete]
//
// DeleteWDSession Wraps the ProxySessionHandler func to delete the ongoing session for manual sessions.
func (p *ProxyManager) DeleteWDSession(w http.ResponseWriter, rq *http.Request) {
	if rq.Method != http.MethodDelete {
		wdproto.BadGatewayError(w, fmt.Errorf("incorrect method"))
		return
	}
	// replace the api/browsers with selenium url
	rq.URL.Path = strings.Replace(rq.URL.Path, "api/browsers", "wd/hub/session", 1)
	p.ProxySessionHandler(w, rq)
}

// StartSessionHandler starts a session
func (p *ProxyManager) StartSessionHandler(w http.ResponseWriter, rq *http.Request) {
	innerCtx, cancel := context.WithCancel(
		trace.ContextWithSpanContext(context.Background(),
			trace.SpanContextFromContext(rq.Context())))
	defer cancel()
	ctx := &Context{Context: innerCtx}

	// we use modified uuid version here
	// so the objects are sorted in Kubernetes/etcd in descending order
	sessionID := uuid.Must(revuuid.NewV7Reverse()).String()
	(&httputil.ReverseProxy{
		Rewrite: func(prq *httputil.ProxyRequest) {
			payload := &bytes.Buffer{}
			pReader := io.TeeReader(prq.In.Body, payload)
			var startSessionRQ wdproto.NewSessionRQ
			if err := json.NewDecoder(pReader).Decode(&startSessionRQ); err != nil {
				wdproto.BadGatewayError(w, err)
				return
			}
			prq.Out.Body = io.NopCloser(payload)
			p.log.Info("capabilities", payload.String())

			prq.Out.Header.Set("sessionID", sessionID)

			if err := adjustCapabilities(&startSessionRQ); err != nil {
				wdproto.BadGatewayError(w, err)
				return
			}

			if err := p.beforeSessionHook(ctx, prq, &startSessionRQ, sessionID); err != nil {
				wdproto.BadGatewayError(w, err)
				return
			}
		},
		ModifyResponse: func(rs *http.Response) error {
			if rs.StatusCode != http.StatusOK {
				return nil
			}
			pRSPayload := &wdproto.NewSessionRS{}
			if err := json.NewDecoder(rs.Body).Decode(pRSPayload); err != nil {
				return errors.WithStack(err)
			}

			originalSession := pRSPayload.Value.SessionID
			pRSPayload.Value.SessionID = sessionID

			replaced, _ := json.Marshal(pRSPayload)

			rs.Header["Content-Length"] = []string{fmt.Sprint(len(replaced))}
			rs.ContentLength = int64(len(replaced))
			rs.Body = io.NopCloser(bytes.NewReader(replaced))

			return p.afterSessionHook(ctx, rs, originalSession)
		},
	}).ServeHTTP(w, rq)
}

func (p *ProxyManager) ProxySessionHandler(w http.ResponseWriter, rq *http.Request) {
	log := p.log

	log.Infof("Parsing session ID from %s", rq.URL.Path)

	sID, command, err := ParseSessionPath(rq.URL.Path)
	if err != nil {
		wdproto.BadGatewayError(w, err)
		return
	}
	log = log.With("session", sID, "method", rq.Method, "command", command)
	log.Info("Execute command")

	// load the session
	sess, err := p.sessionRepo.FindByID(sID)
	if err != nil {
		log.Error("unable to find session")
		wdproto.BadGatewayError(w, err)
		return
	}

	if sess == nil {
		log.Error("unable to find session")
		wdproto.BadGatewayError(w, errors.New("session is nil"))
		return
	}

	/*	if sess.Caps.BrowserKubeOpts.Manual {
		return
	}*/

	parentCtx := otel.GetTextMapPropagator().Extract(rq.Context(), propagation.MapCarrier(sess.Browser.Annotations))
	innerCtx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	ctx := &Context{Context: innerCtx}

	payload := &bytes.Buffer{}
	if _, err = io.Copy(payload, rq.Body); err != nil {
		log.Error("unable to copy rq.Body: %v", err)
	}
	rq.Body = io.NopCloser(payload)

	(&httputil.ReverseProxy{
		Rewrite: func(prq *httputil.ProxyRequest) {
			if err = p.beforeCommandHook(ctx, prq, sess); err != nil {
				wdproto.BadGatewayError(w, err)
				return
			}
			pURL, pErr := url.Parse(sess.Browser.Status.SeleniumURL)
			if pErr != nil {
				wdproto.BadGatewayError(w, pErr)
				return
			}

			adjustUploadPath(prq)
			pURL.Path = path.Clean(path.Join(pURL.Path, RemoveBase(rq.URL.Path)))
			prq.Out.URL = pURL
			p.cleanupOriginHeaders(prq.Out)

			log.Info("Proxying request to ", rq.URL.String())
		},
		ModifyResponse: func(rs *http.Response) error {
			rs.Request.Body = rq.Body

			// delete the pod if this is a QUIT request
			aErr := p.afterCommandHook(ctx, rs, sess, command)
			if aErr != nil {
				log.Error("Command hook error", aErr)
			}

			if p.isQuit(rs.Request.Method, command) {
				if qErr := p.quitSessionHook(ctx, sess); qErr != nil {
					log.Error("Session quit error", qErr)
				}
				return nil
			}
			return aErr
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			// if this is session termination, the session still need to be deleted
			if p.isQuit(r.Method, command) {
				if qErr := p.quitSessionHook(ctx, sess); qErr != nil {
					log.Error("Session quit error", qErr)
				}
			}
		},
	}).ServeHTTP(w, rq)
}

func (p *ProxyManager) ProxyBidirectionalSession(w http.ResponseWriter, rq *http.Request) {
	sessionID := chi.URLParam(rq, keySessionID)
	if sessionID == "" {
		p.log.With("request", fmt.Sprintf("%s %s", rq.Method, rq.URL.Path)).Error("can't parse session ID")
		return
	}

	sess, err := p.sessionRepo.FindByID(sessionID)
	if err != nil {
		p.log.With("request", fmt.Sprintf("%s %s", rq.Method, rq.URL.Path)).Error("session ID not found")
		return
	}

	// sidecar must know the URL generated by selenium webdriver
	rq.Header.Set(keySessionID, sessionID)

	u := &url.URL{
		Scheme: "ws",
		Host:   net.JoinHostPort(sess.Browser.Status.Host, sess.Browser.Status.PortConfig.Sidecar),
		Path:   "/wd/hub/bidi/" + sessionID,
	}

	proxy, err := websocketproxy.NewProxy(u)
	if err != nil {
		p.log.Errorf("error while creating proxy: %v", err)
	}
	proxy.ServeHTTP(w, rq)
}

func (p *ProxyManager) ProxyCDPSession(w http.ResponseWriter, rq *http.Request) {
	sessionID := chi.URLParam(rq, keySessionID)
	if sessionID == "" {
		p.log.With("request", fmt.Sprintf("%s %s", rq.Method, rq.URL.Path)).Error("can't parse session ID")
		return
	}
	sess, err := p.sessionRepo.FindByID(sessionID)
	if err != nil {
		p.log.With("request", fmt.Sprintf("%s %s", rq.Method, rq.URL.Path)).Error("session ID not found")
		return
	}

	// sidecar must know the URL generated by selenium webdriver
	rq.Header.Set(keySessionID, sessionID)

	u := &url.URL{
		Scheme: "ws",
		Host:   net.JoinHostPort(sess.Browser.Status.Host, sess.Browser.Status.PortConfig.Sidecar),
		Path:   "/wd/hub/cdp/" + sessionID,
	}

	proxy, err := websocketproxy.NewProxy(u)
	if err != nil {
		p.log.Errorf("error while creating proxy: %w", err)
	}
	proxy.ServeHTTP(w, rq)
}

func (p *ProxyManager) ProxyDownloads(w http.ResponseWriter, rq *http.Request) error {
	sessionID := chi.URLParam(rq, keySessionID)
	if sessionID == "" {
		return errors.New("unable to parse session ID")
	}
	log := p.log.With("session", sessionID)
	// load the session
	sess, err := p.sessionRepo.FindByID(sessionID)
	if err != nil || sess == nil {
		return errors.Wrapf(err, "unable to find session: %s", sessionID)
	}

	filePath := chi.URLParam(rq, "*")
	(&httputil.ReverseProxy{
		Director: func(rq *http.Request) {
			u := &url.URL{
				Scheme: "http",
				Host:   net.JoinHostPort(sess.Browser.Status.Host, sess.Browser.Status.PortConfig.Sidecar),
			}
			if filePath == "" {
				u.Path = downloadsBasePath
			} else {
				u.Path = path.Join(downloadsBasePath, filePath)
			}
			rq.URL = u
			log.Info("Proxying request to ", rq.URL.String())
		},
	}).ServeHTTP(w, rq)
	return nil
}

func (p *ProxyManager) isQuit(method, command string) bool {
	return method == http.MethodDelete && command == ""
}

func (p *ProxyManager) cleanupOriginHeaders(out *http.Request) {
	out.Header.Del("Origin")
	for name := range out.Header {
		if strings.HasPrefix(name, "Access-Control") {
			out.Header.Del(name)
		}
	}
}

func adjustCapabilities(rq *wdproto.NewSessionRQ) error {
	if rq.Capabilities.BrowserName == "" {
		var validCaps session.Capabilities
		if rq.W3CCapabilities.Capabilities.BrowserName != "" {
			validCaps = rq.W3CCapabilities.Capabilities
		} else {
			for _, caps := range rq.W3CCapabilities.FirstMatch {
				if caps.BrowserName != "" {
					validCaps = *caps
					break
				}
			}
		}
		if err := mergo.Merge(&rq.Capabilities, &validCaps); err != nil {
			return errors.WithStack(err)
		}
	}
	if rq.Capabilities.BrowserKubeOpts.Type == "" {
		rq.Capabilities.BrowserKubeOpts.Type = v1.TypeWebDriver
	}
	return nil
}

func adjustUploadPath(rq *httputil.ProxyRequest) {
	seUploadPath, uploadPath := "/se/file", "/file"
	if strings.HasSuffix(rq.In.URL.Path, seUploadPath) {
		rq.Out.URL.Path = strings.TrimSuffix(rq.In.URL.Path, seUploadPath) + uploadPath
	}
}

// chain builds a handler composed of an inline middleware stack and root
// handler in the order they are passed.
func chain[H any](middlewares []func(H) H, root H) H {
	// Return ahead of time if there aren't any middlewares for the chain
	if len(middlewares) == 0 {
		return root
	}

	// Wrap the end handler with the middleware chain
	h := middlewares[len(middlewares)-1](root)
	for i := len(middlewares) - 2; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}
