package wd

import (
	"net/http"

	"github.com/browserkube/browserkube/cmd/browserkube/internal/wd/plugin/reportcommand"
	"github.com/browserkube/browserkube/cmd/browserkube/internal/wd/plugin/reportlog"
	"github.com/browserkube/browserkube/cmd/browserkube/internal/wd/plugin/reportportal"
	"github.com/browserkube/browserkube/cmd/browserkube/internal/wd/plugin/sessionresult"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"

	"github.com/browserkube/browserkube/cmd/browserkube/internal/provision"
	"github.com/browserkube/browserkube/cmd/browserkube/internal/wd/plugin/reportvideo"
	browserkubehttp "github.com/browserkube/browserkube/pkg/http"
	"github.com/browserkube/browserkube/pkg/opentelemetry"
	"github.com/browserkube/browserkube/pkg/session"
	pkgsessionresult "github.com/browserkube/browserkube/pkg/sessionresult"
	"github.com/browserkube/browserkube/pkg/storage"
	"github.com/browserkube/browserkube/pkg/wd"
)

var Module = fx.Options(
	fx.Provide(func(
		serviceProvider provision.Provisioner,
		sessionResultsRepo pkgsessionresult.Repository,
		store storage.BlobSessionStorage,
		client *http.Client,
		clientset *kubernetes.Clientset,
		envConfig *provision.Config,
	) wd.PluginOpts {
		var plugins []wd.PluginOpt

		plugins = append(plugins, reportvideo.NewReportVideoPlugin(client, store))              // weight 251
		plugins = append(plugins, reportportal.NewReportPortalPlugins(clientset, envConfig)...) // weight 250
		plugins = append(plugins,
			reportlog.NewReportLogPlugin(serviceProvider, store),            // weight 250
			reportcommand.NewReportCommandPlugin(store),                     // weight 250
			opentelemetry.NewMetricsProxyPlugin(),                           // weight 1
			sessionresult.NewSessionResultPlugin(sessionResultsRepo, store), // weight 1
		)
		plugins = append(plugins, NewK8SProxyPlugins(serviceProvider)...) // weight 1

		return plugins
	},
	),
	fx.Invoke(initRoutes),
)

func initRoutes(params inputParams) {
	provider, err := opentelemetry.InitProvider("proxy")
	if err != nil {
		zap.S().Error("failed to initialize provider, error: ", err)
	}

	proxy := wd.NewProxyBuilder(params.PluginOpts...).Build(params.SessionRepo)
	params.Mux.Group(func(r chi.Router) {
		if provider != nil {
			r.Use(opentelemetry.HTTPMiddleware(provider))
		} else {
			r.Use(opentelemetry.NewMetricsMiddleware("proxy"))
		}

		r.HandleFunc("/api/browsers/*", proxy.DeleteWDSession)
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
	PluginOpts  wd.PluginOpts
}
