package main

import (
	"github.com/browserkube/browserkube/browserkube/internal/screenshot"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/fx"
	"k8s.io/utils/env"

	_ "github.com/browserkube/browserkube/browserkube/docs"
	"github.com/browserkube/browserkube/browserkube/internal/api"
	"github.com/browserkube/browserkube/browserkube/internal/api/swagger"
	"github.com/browserkube/browserkube/browserkube/internal/playwright"
	"github.com/browserkube/browserkube/browserkube/internal/provision"
	provisionk8s "github.com/browserkube/browserkube/browserkube/internal/provision/k8s"
	"github.com/browserkube/browserkube/browserkube/internal/reportcommand"
	"github.com/browserkube/browserkube/browserkube/internal/reportlog"
	"github.com/browserkube/browserkube/browserkube/internal/reportportal"
	"github.com/browserkube/browserkube/browserkube/internal/reportvideo"
	"github.com/browserkube/browserkube/browserkube/internal/sessionresult"
	"github.com/browserkube/browserkube/browserkube/internal/wd"
	browserkubeapp "github.com/browserkube/browserkube/pkg/app"
	browserkubehttp "github.com/browserkube/browserkube/pkg/http"
	"github.com/browserkube/browserkube/pkg/opentelemetry"
	"github.com/browserkube/browserkube/pkg/storage"
)

//	@title			BrowserKube
//	@version		1.0
//	@description	Kubernetes-native browser farm
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Andrei Varabyeu
//	@contact.email	andrei.varabyeu@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath	/browserkube

//nolint:gosec // not a credentials
const nsSecret = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

func main() {
	browserkubeapp.Run(
		// generic modules
		browserkubehttp.Module,
		swagger.Module,
		opentelemetry.Module,

		storage.Module,

		// provision modules
		provisionk8s.Module,
		fx.Provide(
			provideProvisionConfig,
		),
		fx.Provide(
			provideHTTPClient,
		),

		// proxy modules
		wd.Module,
		playwright.Module,
		reportportal.Module,
		reportlog.Module,
		// TODO: cases need to be improved when automatic screenshots are required
		screenshot.Module,
		reportvideo.Module,
		reportcommand.Module,

		sessionresult.Module,

		// main ui module
		api.Module,
	)
}

func provideProvisionConfig() (*provision.Config, error) {
	cfg := &provision.Config{
		BrowserNS: env.GetString("BROWSER_NS", ""),
	}
	if cfg.BrowserNS == "" {
		var err error
		cfg.BrowserNS, err = getCurrentNamespace()
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to find out current namespace: %v", err)
		}
	}
	return cfg, nil
}

func getCurrentNamespace() (string, error) {
	ns, err := os.ReadFile(nsSecret)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return string(ns), nil
}

func provideHTTPClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
	}
}
