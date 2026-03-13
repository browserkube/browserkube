package pluginmaximize

import (
	"errors"
	"net/http"

	"github.com/browserkube/browserkube/cmd/browserkube/internal/provision"
	"github.com/browserkube/browserkube/cmd/browserkube/internal/wd/wdctx"
	"github.com/browserkube/browserkube/pkg/wd"
	"github.com/browserkube/browserkube/pkg/wd/wdproto"
	"go.uber.org/zap"
)

func NewMaximizeOnStartPlugin(serviceProvider provision.Provisioner) []wd.PluginOpt {
	return []wd.PluginOpt{
		wd.WithAfterSessionCreated(maximizeWindowOnStart()),
	}
}

// maximizeWindowOnStart maximizes window on start
//
//nolint:deadcode
func maximizeWindowOnStart() func(next wd.OnAfterSessionStart) wd.OnAfterSessionStart {
	return func(next wd.OnAfterSessionStart) wd.OnAfterSessionStart {
		return func(ctx *wd.Context, rs *http.Response, sID string) error {
			if err := maximize(ctx, sID); err != nil {
				zap.S().Errorf("unable to maximize window: %+v", err)
			}
			return next(ctx, rs, sID)
		}
	}
}
func maximize(ctx *wd.Context, sessionID string) error {
	sessionRemote, found := wdctx.GetBrowser(ctx)
	if !found {
		zap.S().Warn("Remote session isn't available while maximizing window")
		return errors.New("session remote isn't found")
	}
	go func() {
		err := wdproto.NewWebDriver(sessionRemote.Status.SeleniumURL, sessionID).Maximize(ctx)
		if err != nil {
			zap.S().Error("Unable to maximize browser", err)
		}
	}()
	return nil
}
