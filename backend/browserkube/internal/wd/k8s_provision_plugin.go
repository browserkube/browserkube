package wd

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/browserkube/browserkube/browserkube/internal/provision"
	browserkubev1 "github.com/browserkube/browserkube/operator/api/v1"
	"github.com/browserkube/browserkube/pkg/session"
	"github.com/browserkube/browserkube/pkg/wd"
	"github.com/browserkube/browserkube/pkg/wd/wdproto"
)

func provideK8SProxyPlugin(serviceProvider provision.Provisioner) wd.PluginOpts {
	return wd.PluginOpts{
		Weight: 1,
		Opts: []wd.PluginOpt{
			wd.WithBeforeSessionCreated(provisionBrowserHandler(serviceProvider)),
			// wd.WithAfterSessionCreated(maximizeWindowOnStart()), //nolint:bodyclose
			wd.WithQuitSession(quitSessionHandler(serviceProvider)),
		},
	}
}

// afterCommandHandler deletes a pod when quit session is requested
func quitSessionHandler(serviceProvider provision.Provisioner) func(next wd.OnSessionQuit) wd.OnSessionQuit {
	return func(next wd.OnSessionQuit) wd.OnSessionQuit {
		return func(ctx *wd.Context, sess *session.Session) error {
			go func(srv *browserkubev1.Browser) {
				if dErr := serviceProvider.Delete(context.Background(), srv.Name); dErr != nil {
					zap.S().Error("Unable to delete provider", dErr)
				}
			}(sess.Browser)
			return next(ctx, sess)
		}
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

// provisionBrowserHandler provisions a browser pod before the session creation
func provisionBrowserHandler(serviceProvider provision.Provisioner) func(next wd.OnBeforeSessionStart) wd.OnBeforeSessionStart {
	return func(next wd.OnBeforeSessionStart) wd.OnBeforeSessionStart {
		return func(ctx *wd.Context, prq *httputil.ProxyRequest, sessionRQ *wdproto.NewSessionRQ, sessionID string) error {
			remoteSelenium, err := serviceProvider.Provision(ctx, sessionID, &sessionRQ.Capabilities)
			if err != nil {
				if remoteSelenium != nil {
					if dErr := serviceProvider.Delete(context.Background(), remoteSelenium.Name); dErr != nil {
						return errors.WithStack(dErr)
					}
				}
				return errors.WithStack(err)
			}

			pURL, err := url.Parse(remoteSelenium.Status.SeleniumURL)
			if err != nil {
				return errors.WithStack(err)
			}
			pURL.Path = prq.In.URL.Path
			prq.Out.URL = pURL

			return next(withBrowser(ctx, remoteSelenium), prq, sessionRQ, sessionID)
		}
	}
}

func maximize(ctx *wd.Context, sessionID string) error {
	sessionRemote, found := getBrowser(ctx)
	if !found {
		zap.S().Warn("Remote session isn't available while maximizing window")
		return errors.New("Session remote isn't found")
	}
	err := wdproto.NewWebDriver(sessionRemote.Status.SeleniumURL, sessionID).Maximize(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
