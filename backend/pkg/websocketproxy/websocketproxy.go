// Package websocketproxy is a reverse proxy for WebSocket connections.
package websocketproxy

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	browserkubeutil "github.com/browserkube/browserkube/pkg/util"
)

var (
	// DefaultUpgrader specifies the parameters for upgrading an HTTP
	// connection to a WebSocket connection.
	DefaultUpgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// DefaultDialer is a dialer with all fields set to the default zero values.
	DefaultDialer = websocket.DefaultDialer

	ErrDoNotSend = errors.New("do not send a message")
)

// WebsocketProxy is an HTTP Handler that takes an incoming WebSocket
// connection and proxies it to another server.
type WebsocketProxy struct {
	// Director, if non-nil, is a function that may copy additional request
	// headers from the incoming WebSocket connection into the output headers
	// which will be forwarded to another server.
	Director func(incoming *http.Request, out http.Header)

	// Backend returns the backend URL which the proxy uses to reverse proxy
	// the incoming WebSocket connection. Request is the initial incoming and
	// unmodified request.
	Backend func(*http.Request) *url.URL

	// Upgrader specifies the parameters for upgrading a incoming HTTP
	// connection to a WebSocket connection. If nil, DefaultUpgrader is used.
	Upgrader *websocket.Upgrader

	//  Dialer contains options for connecting to the backend WebSocket server.
	//  If nil, DefaultDialer is used.
	Dialer *websocket.Dialer

	WSConn *websocket.Conn

	onIncomingMessageF []OnMessageFunc
	onOutgoingMessageF []OnMessageFunc

	log *zap.SugaredLogger
}

type (
	ProxyOpt      func(*WebsocketProxy) error
	OnMessageFunc func(msg *Message) error
)

func WithIncomingMiddleware(f OnMessageFunc) ProxyOpt {
	return func(proxy *WebsocketProxy) error {
		proxy.onIncomingMessageF = append(proxy.onIncomingMessageF, f)
		return nil
	}
}

func WithOutgoingMiddleware(f OnMessageFunc) ProxyOpt {
	return func(proxy *WebsocketProxy) error {
		proxy.onOutgoingMessageF = append(proxy.onOutgoingMessageF, f)
		return nil
	}
}

// ProxyHandler returns a new http.Handler interface that reverse proxies the
// request to the given target.
// func ProxyHandler(target *url.URL) http.Handler { return NewProxy(target, nil, nil) }

// NewProxy returns a new Websocket reverse proxy that rewrites the
// URL's to the scheme, host and base path provider in target.
func NewProxy(target *url.URL, opts ...ProxyOpt) (*WebsocketProxy, error) {
	backend := func(r *http.Request) *url.URL {
		// Shallow copy
		u := *target
		u.Fragment = r.URL.Fragment
		// Do not take path since in multi phase proxies(Ex:backend -> sidecar -> webdriver) it can point to a wrong address
		// u.Path = r.URL.Path
		u.RawQuery = r.URL.RawQuery
		return &u
	}
	p := &WebsocketProxy{Backend: backend, log: zap.S()}
	for _, opt := range opts {
		err := opt(p)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

// ServeHTTP implements the http.Handler that proxies WebSocket connections.
func (w *WebsocketProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) { //nolint:all
	if w.Backend == nil {
		log.Println("websocketproxy: backend function is not defined")
		http.Error(rw, "internal server error (code: 1)", http.StatusInternalServerError)
		return
	}

	backendURL := w.Backend(req)
	if backendURL == nil {
		log.Println("websocketproxy: backend URL is nil")
		http.Error(rw, "internal server error (code: 2)", http.StatusInternalServerError)
		return
	}

	dialer := w.Dialer
	if w.Dialer == nil {
		dialer = DefaultDialer
	}

	// Pass headers from the incoming request to the dialer to forward them to
	// the final destinations.
	requestHeader := http.Header{}
	if origin := req.Header.Get("Origin"); origin != "" {
		requestHeader.Add("Origin", origin)
	}
	for _, prot := range req.Header[http.CanonicalHeaderKey("Sec-WebSocket-Protocol")] {
		requestHeader.Add("Sec-WebSocket-Protocol", prot)
	}
	for _, cookie := range req.Header[http.CanonicalHeaderKey("Cookie")] {
		requestHeader.Add("Cookie", cookie)
	}
	if req.Host != "" {
		requestHeader.Set("Host", req.Host)
	}

	// Pass X-Forwarded-For headers too, code below is a part of
	// httputil.ReverseProxy. See http://en.wikipedia.org/wiki/X-Forwarded-For
	// for more information
	// TODO: use RFC7239 http://tools.ietf.org/html/rfc7239
	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		// If we aren't the first proxy retain prior
		// X-Forwarded-For information as a comma+space
		// separated list and fold multiple headers into one.
		if prior, ok := req.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		requestHeader.Set("X-Forwarded-For", clientIP)
	}

	// Set the originating protocol of the incoming HTTP request. The SSL might
	// be terminated on our site and because we doing proxy adding this would
	// be helpful for applications on the backend.
	requestHeader.Set("X-Forwarded-Proto", "http")
	if req.TLS != nil {
		requestHeader.Set("X-Forwarded-Proto", "https")
	}

	// Enable the director to copy any additional headers it desires for
	// forwarding to the remote server.
	if w.Director != nil {
		w.Director(req, requestHeader)
	}

	// Connect to the backend URL, also pass the headers we get from the requst
	// together with the Forwarded headers we prepared above.
	// TODO: support multiplexing on the same backend connection instead of
	// opening a new TCP connection time for each request. This should be
	// optional:
	// http://tools.ietf.org/html/draft-ietf-hybi-websocket-multiplexing-01
	w.log.Infof("Proxy: Proxying to %s", backendURL.String())
	connBackend, resp, err := dialer.Dial(backendURL.String(), requestHeader) //nolint:bodyclose
	if err != nil {
		log.Printf("websocketproxy: couldn't dial to remote backend url %s", err)
		if resp != nil {
			log.Printf("Proxy: error status: %d", resp.StatusCode)
			// If the WebSocket handshake fails, ErrBadHandshake is returned
			// along with a non-nil *http.Response so that callers can handle
			// redirects, authentication, etcetera.
			if cErr := copyResponse(rw, resp); cErr != nil {
				log.Printf("websocketproxy: couldn't write response after failed remote backend handshake: %s", cErr)
			}
		} else {
			http.Error(rw, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		}
		return
	}
	defer browserkubeutil.CloseQuietly(connBackend)

	w.WSConn = connBackend

	upgrader := w.Upgrader
	if w.Upgrader == nil {
		upgrader = DefaultUpgrader
	}

	// Only pass those headers to the upgrader.
	upgradeHeader := http.Header{}
	if hdr := resp.Header.Get("Sec-Websocket-Protocol"); hdr != "" {
		upgradeHeader.Set("Sec-Websocket-Protocol", hdr)
	}
	if hdr := resp.Header.Get("Set-Cookie"); hdr != "" {
		upgradeHeader.Set("Set-Cookie", hdr)
	}

	// Now upgrade the existing incoming request to a WebSocket connection.
	// Also pass the header that we gathered from the Dial handshake.
	connPub, err := upgrader.Upgrade(rw, req, upgradeHeader)
	if err != nil {
		log.Printf("websocketproxy: couldn't upgrade %s", err)
		return
	}
	defer browserkubeutil.CloseQuietly(connPub)

	errClient := make(chan error, 1)
	errBackend := make(chan error, 1)
	replicateWebsocketConn := func(dst, src *websocket.Conn, middleware []OnMessageFunc, errc chan error) {
		for {
			msgType, msg, rErr := src.ReadMessage()
			if rErr != nil {
				m := websocket.FormatCloseMessage(websocket.CloseNormalClosure, fmt.Sprintf("%v", rErr))
				var closeErr *websocket.CloseError
				if errors.As(rErr, &closeErr) {
					if closeErr.Code != websocket.CloseNoStatusReceived {
						m = websocket.FormatCloseMessage(closeErr.Code, closeErr.Text)
					}
				}
				errc <- rErr
				if wErr := dst.WriteMessage(websocket.CloseMessage, m); wErr != nil {
					w.log.Error(wErr)
				}
				break
			}

			if len(middleware) > 0 {
				runMware := func() ([]byte, error) {
					var serMsg Message
					if serErr := json.Unmarshal(msg, &serMsg); serErr != nil {
						zap.S().Error(serErr)
						return nil, serErr
					}
					for _, mware := range middleware {
						if mwareErr := mware(&serMsg); mwareErr != nil {
							zap.S().Error(mwareErr)
							return nil, mwareErr
						}
					}
					newMsg, serErr := json.Marshal(&serMsg)
					if serErr != nil {
						zap.S().Error(serErr)
						return nil, serErr
					}
					return newMsg, nil
				}
				newMsg, mErr := runMware()
				if errors.Is(mErr, ErrDoNotSend) {
					break
				}
				if mErr == nil {
					msg = newMsg
				}
			}
			err = dst.WriteMessage(msgType, msg)
			if err != nil {
				errc <- err
				break
			}
		}
	}

	go replicateWebsocketConn(connPub, connBackend, w.onOutgoingMessageF, errClient)
	go replicateWebsocketConn(connBackend, connPub, w.onIncomingMessageF, errBackend)

	var message string
	select {
	case err = <-errClient:
		message = "websocketproxy: Error when copying from backend to client: %v"
	case err = <-errBackend:
		message = "websocketproxy: Error when copying from client to backend: %v"
	}

	var closeErr *websocket.CloseError
	if !errors.As(err, &closeErr) || closeErr.Code == websocket.CloseAbnormalClosure {
		w.log.Errorw(message, "error", err)
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func copyResponse(rw http.ResponseWriter, resp *http.Response) error {
	copyHeader(rw.Header(), resp.Header)
	rw.WriteHeader(resp.StatusCode)
	defer browserkubeutil.CloseQuietly(resp.Body)

	_, err := io.Copy(rw, resp.Body)
	return err
}
