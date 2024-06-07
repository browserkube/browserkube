package wd

import (
	"sort"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"

	browserkubehttp "github.com/browserkube/browserkube/pkg/http"
	"github.com/browserkube/browserkube/pkg/opentelemetry"
	"github.com/browserkube/browserkube/pkg/session"
	"github.com/browserkube/browserkube/pkg/wd"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			provideK8SProxyPlugin,
			fx.ResultTags(`group:"wd-extensions"`),
		),
	),
	fx.Invoke(initRoutes),
)

func initRoutes(params inputParams) {
	provider, err := opentelemetry.InitProvider("proxy")
	if err != nil {
		zap.S().Error("failed to initialize provider, error: ", err)
	}

	var opts []wd.PluginOpt

	// Sort PluginOpts by their weight. Plugin with the highest weight applies first.
	sort.Slice(params.PluginOpts, func(i, j int) bool {
		return params.PluginOpts[i].Weight > params.PluginOpts[j].Weight
	})
	for _, opt := range params.PluginOpts {
		opts = append(opts, opt.Opts...)
	}
	proxy := wd.NewProxyBuilder(opts...).Build(params.SessionRepo)
	params.Mux.Group(func(r chi.Router) {
		if provider != nil {
			r.Use(opentelemetry.HTTPMiddleware(provider))
		} else {
			r.Use(opentelemetry.NewMetricsMiddleware("proxy"))
		}

		// V2
		r.HandleFunc("/api/browsers", proxy.CreateWDSession)
		r.HandleFunc("/api/browsers/*", proxy.DeleteWDSession)
		//
		r.HandleFunc("/wd/hub/session", proxy.StartSessionHandler)
		r.HandleFunc("/wd/hub/session/*", proxy.ProxySessionHandler)
		r.HandleFunc("/wd/hub/bidi/{sessionID}", proxy.ProxyBidirectionalSession)
		r.HandleFunc("/wd/hub/cdp/{sessionID}", proxy.ProxyCDPSession)
	})

	// downloads server
	params.Mux.Group(func(r chi.Router) {
		if provider != nil {
			r.Use(opentelemetry.HTTPMiddleware(provider))
		} else {
			r.Use(opentelemetry.NewMetricsMiddleware("downloads"))
		}

		params.Mux.Handle("/wd/hub/session/{sessionID}/browserkube/downloads/*", browserkubehttp.Handler(proxy.ProxyDownloads))
	})
}

type inputParams struct {
	fx.In
	Mux         chi.Router
	SessionRepo session.Repository
	PluginOpts  []wd.PluginOpts `group:"wd-extensions"`
}
