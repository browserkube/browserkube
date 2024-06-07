package swagger

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"

	browserkubeutil "github.com/browserkube/browserkube/pkg/util"
)

var Module = fx.Options(
	fx.Invoke(initSwagger),
)

func initSwagger(mux chi.Router) {
	mux.Group(func(r chi.Router) {
		// Swagger Docs
		const headerXForwarded = "X-Forwarded-Prefix"
		mux.Get("/swagger/", func(w http.ResponseWriter, rq *http.Request) {
			if prefix := rq.Header.Get(headerXForwarded); prefix != "" {
				rq.RequestURI, _ = url.JoinPath(prefix, rq.RequestURI)
			}
			p, _ := url.JoinPath(rq.RequestURI, "/index.html")
			http.Redirect(w, rq, p, http.StatusMovedPermanently)
		})
		mux.Get("/swagger/*", func(w http.ResponseWriter, rq *http.Request) {
			docPath := "/swagger/doc.json"
			if prefix := rq.Header.Get(headerXForwarded); prefix != "" {
				docPath, _ = url.JoinPath(prefix, docPath)
			}

			scheme := browserkubeutil.FirstNonEmpty(rq.URL.Scheme, "https")
			u := url.URL{
				Scheme: scheme,
				Host:   rq.Host,
				Path:   docPath,
			}
			zap.S().Info("Request URI")
			zap.S().Info(rq.RequestURI)
			httpSwagger.Handler(
				httpSwagger.URL(u.String()), // The url pointing to API definition
			).ServeHTTP(w, rq)
		})
	})
}
