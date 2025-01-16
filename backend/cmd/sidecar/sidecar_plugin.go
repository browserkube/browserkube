package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/fx"
	"go.uber.org/zap"

	browserkubeutil "github.com/browserkube/browserkube/pkg/util"
	"github.com/browserkube/browserkube/pkg/wd"
	"github.com/browserkube/browserkube/pkg/wd/wdproto"
	"github.com/browserkube/browserkube/pkg/websocketproxy"
)

const (
	wsURLFieldKey  = "webSocketUrl"
	cdpURLFieldKey = "se:cdp"
)

type session struct {
	ID string
}

type bidiConfig struct {
	BiDiURL string
	CDPURL  string
}

type wdProxy struct {
	sessionTimeout *idleTimer
	idleTimer      *idleTimer
	sessionID      *browserkubeutil.TypedAtomic[*session]
	bidiURL        *browserkubeutil.TypedAtomic[*bidiConfig]
	proxyURL       *url.URL
	browserHomeDir string
	logger         *zap.SugaredLogger
	counter        atomic.Int32
}

func newWDProxy(c *conf, quit fx.Shutdowner) *wdProxy {
	logger := zap.S()
	sessionID := browserkubeutil.NewTypedAtomic[*session]()
	bidiURL := browserkubeutil.NewTypedAtomic[*bidiConfig]()
	proxyURL := c.proxyURL
	timeoutCloseFunc := func() {
		logger.Warn("Closing session due to timeout")
		sessID := sessionID.Load()
		if sessID == nil {
			return
		}
		if err := wdproto.NewWebDriver(proxyURL.String(), sessID.ID).Quit(context.Background()); err != nil {
			logger.Errorf("unable to quit webdriver: %s", err.Error())
		}
	}
	idleTimer := newIdleTimer(logger, quit, c.idleTimeout, timeoutCloseFunc)
	sessionTimer := newIdleTimer(logger, quit, c.sessionTimeout, timeoutCloseFunc)
	return &wdProxy{
		idleTimer:      idleTimer,
		sessionTimeout: sessionTimer,
		proxyURL:       proxyURL,
		sessionID:      sessionID,
		bidiURL:        bidiURL,
		logger:         logger,
		browserHomeDir: c.browserHomeDir,
		counter:        atomic.Int32{},
	}
}

// StartSessionHandler starts a session
func (p *wdProxy) StartSessionHandler(w http.ResponseWriter, rq *http.Request) {
	(&httputil.ReverseProxy{
		Rewrite: func(prq *httputil.ProxyRequest) {
			if s := p.sessionID.Load(); s != nil && s.ID != "" {
				wdproto.BadGatewayError(w, errors.New("only single session is supported"))
				return
			}
			prq.Out.URL.Scheme = p.proxyURL.Scheme
			prq.Out.URL.Host = p.proxyURL.Host
			prq.Out.URL.Path = path.Clean(path.Join(p.proxyURL.Path, wd.RemoveBase(rq.URL.Path)))

			// delete host and origin since drivers
			// may support local connections only
			prq.Out.Host = ""
			prq.Out.Header.Del("Origin")
			prq.Out.Header.Del("Host")

			p.logger.Infof("Proxying [%s] request to [%s]", prq.Out.Method, prq.Out.URL.String())
			p.logger.Info(prq.Out.Header)
		},
		ModifyResponse: func(rs *http.Response) error {
			if rs.StatusCode != http.StatusOK {
				return nil
			}
			payload := &bytes.Buffer{}
			pReader := io.TeeReader(rs.Body, payload)
			var sessionRS wdproto.NewSessionRS
			if err := json.NewDecoder(pReader).Decode(&sessionRS); err != nil {
				return errors.WithStack(err)
			}

			// SessionID given by webdriver
			oldSessionID := sessionRS.Value.SessionID
			// SessionID given by the browserkube
			newSessionID := rq.Header.Get("sessionID")
			bidiConf := &bidiConfig{}
			// Bi-directional capability
			webdriverWsURL := replaceWebsocketURL(rq, &sessionRS, newSessionID, wsURLFieldKey)
			if webdriverWsURL != "" {
				p.logger.Infof("SessionID: %s, wsURL: %s", newSessionID, webdriverWsURL)
				bidiConf.BiDiURL = webdriverWsURL
			}
			// Chrome Dev Tools capability
			cdpWsURL := replaceWebsocketURL(rq, &sessionRS, newSessionID, cdpURLFieldKey)
			if cdpWsURL != "" {
				p.logger.Infof("SessionID: %s, cdpURL: %s", newSessionID, cdpWsURL)
				bidiConf.CDPURL = cdpWsURL
			}
			if webdriverWsURL != "" || cdpWsURL != "" {
				p.bidiURL.Set(bidiConf)
			}

			replaced, _ := json.Marshal(sessionRS)
			rs.Header["Content-Length"] = []string{fmt.Sprint(len(replaced))}
			rs.ContentLength = int64(len(replaced))
			rs.Body = io.NopCloser(bytes.NewReader(replaced))

			p.sessionID.Set(&session{ID: oldSessionID})
			p.logger.With("session", oldSessionID).Infof("Session Creation: %s", rs.Status)

			return nil
		},
	}).ServeHTTP(w, rq)
}

func (p *wdProxy) ProxySessionHandler(w http.ResponseWriter, rq *http.Request) {
	(&httputil.ReverseProxy{
		Rewrite: func(prq *httputil.ProxyRequest) {
			sess := p.sessionID.Load()
			if sess == nil {
				wdproto.BadGatewayError(w, errors.New("there is no active session"))
				return
			}
			p.logger.Infof("Replacing %s with new session %s", rq.URL.Path, sess.ID)

			p.idleTimer.Reset()
			commandPath, err := wd.ReplaceSession(wd.RemoveBase(rq.URL.Path), sess.ID)
			if err != nil {
				wdproto.BadGatewayError(w, err)
				return
			}

			prq.Out.URL.Scheme = p.proxyURL.Scheme
			prq.Out.URL.Host = p.proxyURL.Host
			prq.Out.URL.Path = path.Clean(path.Join(p.proxyURL.Path, commandPath))

			// delete host and origin since drivers
			// may support local connections only
			prq.Out.Host = ""
			prq.Out.Header.Del("Origin")
			prq.Out.Header.Del("Host")

			p.logger.With("session", p.sessionID.Load()).Infof("Proxying [%s] request to [%s]", prq.Out.Method, prq.Out.URL.String())
		},
		ModifyResponse: func(rs *http.Response) error {
			newCommand := p.counter.Add(1)
			rs.Header.Set("commandID", strconv.FormatInt(int64(newCommand), 10))
			return nil
		},
	}).ServeHTTP(w, rq)
}

func (p *wdProxy) ProxyBidirectionalSession(w http.ResponseWriter, rq *http.Request) {
	wdWebSocketURL := p.bidiURL.Load().BiDiURL
	if wdWebSocketURL == "" {
		p.logger.Error("web socket url is empty")
	}

	u, err := url.Parse(wdWebSocketURL)
	if err != nil {
		p.logger.With("request", fmt.Sprintf("%s %s", rq.Method, rq.URL.Path)).Error("error parsing url")
		return
	}

	p.logger.Infof("Proxying Webdriver websocket to %s", u)
	proxy, err := websocketproxy.NewProxy(u)
	if err != nil {
		p.logger.Errorf("error while proxying from sidecar: %w", err)
	}

	proxy.ServeHTTP(w, rq)
}

func (p *wdProxy) ProxyCDPSession(w http.ResponseWriter, rq *http.Request) {
	cdpWebsocketURL := p.bidiURL.Load().CDPURL
	if cdpWebsocketURL == "" {
		p.logger.Error("cdp websocket url is empty")
	}

	u, err := url.Parse(cdpWebsocketURL)
	if err != nil {
		p.logger.With("request", fmt.Sprintf("%s %s", rq.Method, rq.URL.Path)).Error("error parsing url")
		return
	}

	p.logger.Infof("Proxying CDP websocket to %s", u)
	proxy, err := websocketproxy.NewProxy(u)
	if err != nil {
		p.logger.Errorf("error while proxying from sidecar: %w", err)
	}

	proxy.ServeHTTP(w, rq)
}

func newIdleTimer(logger *zap.SugaredLogger, quit fx.Shutdowner, timeout time.Duration, callback func()) *idleTimer {
	return &idleTimer{
		timeout: timeout,
		timer: time.AfterFunc(timeout, func() {
			callback()
			if err := quit.Shutdown(); err != nil {
				logger.Errorf("shutdown error: %+v", err)
			}
		}),
	}
}

// replaceWebsocketUrl replaces selenium generated websocket URL (capabilities in response)
// by new one that points to browserkube backend proxy endpoint and returns selenium generated URL.
func replaceWebsocketURL(rq *http.Request, rsPayload *wdproto.NewSessionRS, newSessionID, fieldKey string) string {
	v, found := rsPayload.Value.Capabilities[fieldKey]
	if !found {
		return ""
	}
	wsURL, ok := v.(string)
	if !ok {
		return ""
	}

	switch fieldKey {
	case wsURLFieldKey:
		rsPayload.Value.Capabilities[fieldKey] = composeWebdriverWsURL(rq, newSessionID)
	case cdpURLFieldKey:
		rsPayload.Value.Capabilities[fieldKey] = composeCDPWsURL(rq, newSessionID)
	}

	return wsURL
}

func composeWebdriverWsURL(rq *http.Request, sessionID string) string {
	// "http://browserkube.fqdn/browserkube/wd/hub/session" -> "http://browserkube.fqdn/browserkube/wd/hub/bidi/abcef12345"
	newURL := &url.URL{
		Scheme: "ws",
		Host:   rq.Host,
		Path:   rq.Header.Get("X-Forwarded-Prefix") + strings.ReplaceAll(rq.URL.Path, "/session", "/bidi/") + sessionID,
	}

	return newURL.String()
}

func composeCDPWsURL(rq *http.Request, sessionID string) string {
	// "http://browserkube.fqdn/browserkube/wd/hub/session" -> "http://browserkube.fqdn/browserkube/wd/hub/cdp/abcef12345"
	newURL := &url.URL{
		Scheme: "ws",
		Host:   rq.Host,
		Path:   rq.Header.Get("X-Forwarded-Prefix") + strings.ReplaceAll(rq.URL.Path, "/session", "/cdp/") + sessionID,
	}

	return newURL.String()
}

type idleTimer struct {
	timer   *time.Timer
	timeout time.Duration
}

func (i *idleTimer) Stop() {
	i.timer.Stop()
}

func (i *idleTimer) Reset() {
	i.timer.Reset(i.timeout)
}
