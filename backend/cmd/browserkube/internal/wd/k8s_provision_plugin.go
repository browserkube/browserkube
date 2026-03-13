package wd

import (
	"context"
	"net/http/httputil"
	"net/url"

	"github.com/browserkube/browserkube/cmd/browserkube/internal/wd/wdctx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/browserkube/browserkube/cmd/browserkube/internal/provision"
	browserkubev1 "github.com/browserkube/browserkube/operator/api/v1"
	"github.com/browserkube/browserkube/pkg/session"
	"github.com/browserkube/browserkube/pkg/wd"
	"github.com/browserkube/browserkube/pkg/wd/wdproto"
)

func NewK8SProxyPlugins(serviceProvider provision.Provisioner) []wd.PluginOpt {
	return []wd.PluginOpt{
		wd.WithBeforeSessionCreated(provisionBrowserHandler(serviceProvider)),
		wd.WithQuitSession(destroyBrowserHandler(serviceProvider)),
	}
}

// destroyBrowserHandler deletes a pod when quit session is requested
func destroyBrowserHandler(serviceProvider provision.Provisioner) func(next wd.OnSessionQuit) wd.OnSessionQuit {
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

			return next(wdctx.WithBrowser(ctx, remoteSelenium), prq, sessionRQ, sessionID)
		}
	}
}
