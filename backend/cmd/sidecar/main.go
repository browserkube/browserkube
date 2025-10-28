package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"go.uber.org/fx"

	browserkubeapp "github.com/browserkube/browserkube/pkg/app"
	browserkubehttp "github.com/browserkube/browserkube/pkg/http"
	"github.com/browserkube/browserkube/pkg/sessionresult"
	browserkubeutil "github.com/browserkube/browserkube/pkg/util"
)

func main() {
	browserkubeapp.Run(
		browserkubehttp.Module,
		fx.Provide(
			provideConfig,
			newWDProxy,
		),
		fx.Invoke(initHandlers),
	)
}

type conf struct {
	proxyURL       *url.URL
	recorderURL    *url.URL
	idleTimeout    time.Duration
	sessionTimeout time.Duration
	browserHomeDir string
}

func provideConfig() (*conf, error) {
	proxyURL := browserkubeutil.FirstNonEmpty(os.Getenv("PROXY_URL"), "http://localhost:4444")
	u, err := url.Parse(proxyURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	recorderURL := browserkubeutil.FirstNonEmpty(os.Getenv("RECORDER_URL"), "http://localhost:5555")
	recURL, err := url.Parse(recorderURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	browserHomeDir := browserkubeutil.FirstNonEmpty(os.Getenv("BROWSER_HOME_DIR"), "/home/user")

	iTimeout, err := time.ParseDuration(browserkubeutil.FirstNonEmpty(os.Getenv("IDLE_TIMEOUT"), "10m"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	sTimeout, err := time.ParseDuration(browserkubeutil.FirstNonEmpty(os.Getenv("SESSION_TIMEOUT"), "1h"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &conf{
		proxyURL:       u,
		recorderURL:    recURL,
		idleTimeout:    iTimeout,
		sessionTimeout: sTimeout,
		browserHomeDir: browserHomeDir,
	}, nil
}

func initHandlers(mux chi.Router, c *conf, proxy *wdProxy) {
	// downloads server
	mux.Handle("/downloads/*", http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		rq.URL.Path = path.Join("/", chi.URLParam(rq, "*"))
		http.FileServer(http.Dir(filepath.Join(c.browserHomeDir, "/Downloads"))).ServeHTTP(w, rq)
	}))
	// videos server
	mux.Handle(path.Join(sessionresult.VideosPath, "*"), http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		rq.URL.Path = path.Join("/", chi.URLParam(rq, "*"))
		http.FileServer(http.Dir(filepath.Join(c.browserHomeDir, "/videos"))).ServeHTTP(w, rq)
	}))

	mux.HandleFunc("/recorder/stop", func(w http.ResponseWriter, rq *http.Request) {
		(&httputil.ReverseProxy{
			Rewrite: func(prq *httputil.ProxyRequest) {
				u := *prq.In.URL
				u.Scheme, u.Host = c.recorderURL.Scheme, c.recorderURL.Host
				prq.SetXForwarded()
				prq.Out.URL = &u
			},
		}).ServeHTTP(w, rq)
	})

	mux.HandleFunc("/wd/hub/session", proxy.StartSessionHandler)
	mux.HandleFunc("/wd/hub/session/*", proxy.ProxySessionHandler)
	mux.HandleFunc("/wd/hub/bidi/{sessionID}", proxy.ProxyBidirectionalSession)
	mux.HandleFunc("/wd/hub/cdp/{sessionID}", proxy.ProxyCDPSession)

	mux.HandleFunc("/wd/hub/status", func(w http.ResponseWriter, rq *http.Request) {
		proxy := &httputil.ReverseProxy{
			Rewrite: func(rq *httputil.ProxyRequest) {
				u := *rq.In.URL
				u.Scheme, u.Host = c.proxyURL.Scheme, c.proxyURL.Host
				u.Path = path.Join(c.proxyURL.Path, "/status")
				rq.SetXForwarded()
				rq.Out.URL = &u
				rq.Out.Host = u.Host
			},
		}
		proxy.ServeHTTP(w, rq)
	})
}
