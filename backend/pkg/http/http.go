package browserkubehttp

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/browserkube/browserkube/pkg/buildinfo"
	browserkubeutil "github.com/browserkube/browserkube/pkg/util"
)

type Port string

var Module = fx.Options(
	fx.Provide(
		provideMux,
		func() Port {
			return Port(browserkubeutil.FirstNonEmpty(os.Getenv("PORT"), "4444"))
		},
	),
	fx.Invoke(startSRV),
)

func provideMux() chi.Router {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Heartbeat("/health"))
	mux.Use(func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if (r.Method == http.MethodGet) && strings.EqualFold(r.URL.Path, "/info") {
				_ = WriteJSON(w, http.StatusOK, buildinfo.GetBuildInfo())
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	})
	mux.Use(cors.Handler(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}))
	return mux
}

//nolint:interfacer
func startSRV(lc fx.Lifecycle, mux chi.Router, port Port) *http.Server {
	log := zap.S()

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           mux,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Infof("Starting HTTP server on port %s", port)

			go func() {
				if err := server.ListenAndServe(); err != nil {
					log.Error(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping HTTP server.")
			err := server.Shutdown(ctx)
			if err != nil {
				log.Error(err)
			}
			log.Info("HTTP server has stopped")
			return nil
		},
	})
	return server
}
