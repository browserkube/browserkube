package reportportal

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/reportportal/goRP/v5/pkg/gorp"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/browserkube/browserkube/pkg/session"
	"github.com/browserkube/browserkube/pkg/wd"
	"github.com/browserkube/browserkube/pkg/wd/wdproto"
)

var Module = fx.Options(
	fx.Provide(
		newSettingsRepo,
		fx.Annotate(
			provideReportPortalPlugin,
			fx.ResultTags(`group:"wd-extensions"`),
		),
	),
)

func provideReportPortalPlugin(sr settingsRepo) wd.PluginOpts {
	return wd.PluginOpts{
		Weight: 250,
		Opts: []wd.PluginOpt{
			wd.WithBeforeSessionCreated(beforeSessionCreated(sr)),
			wd.WithAfterCommand(afterCommandHandler(sr)), //nolint:bodyclose
			wd.WithAfterCommand(findElementHandler(sr)),  //nolint:bodyclose
			wd.WithQuitSession(onQuitSession(sr)),
		},
	}
}

func onQuitSession(sr settingsRepo) func(next wd.OnSessionQuit) wd.OnSessionQuit {
	return func(next wd.OnSessionQuit) wd.OnSessionQuit {
		return func(ctx *wd.Context, s *session.Session) error {
			rp := s.Caps.BrowserKubeOpts.RP
			if rp == nil {
				return next(ctx, s)
			}
			shouldFinishTest := rp.FinishItem
			if shouldFinishTest {
				launchID := rp.LaunchID
				itemID := rp.ItemID
				project := rp.Project

				log := zap.S().With("project", project, "launchId", launchID, "itemId", itemID)
				settings, err := sr.FindByProjectName(ctx, project)
				if err != nil {
					log.Errorf("ReportPortal isn't properly configured: %s", err)
					return next(ctx, s)
				}
				rp := gorp.NewClient(settings.Host, project, settings.AuthToken)

				_, err = rp.FinishTest(itemID, &gorp.FinishTestRQ{
					LaunchUUID: launchID,
					FinishExecutionRQ: gorp.FinishExecutionRQ{
						Status: gorp.Statuses.Passed,
					},
				})
				if err != nil {
					log.Errorf("Unable to finish test in ReportPortal: %s", err)
				}
			}
			return next(ctx, s)
		}
	}
}

func beforeSessionCreated(sr settingsRepo) func(next wd.OnBeforeSessionStart) wd.OnBeforeSessionStart {
	return func(next wd.OnBeforeSessionStart) wd.OnBeforeSessionStart {
		return func(ctx *wd.Context, prq *httputil.ProxyRequest, sessionRQ *wdproto.NewSessionRQ, sID string) error {
			if sessionRQ.Capabilities.BrowserKubeOpts.RP == nil {
				// just skip to next plugin execution
				return next(ctx, prq, sessionRQ, sID)
			}
			launchID := sessionRQ.Capabilities.BrowserKubeOpts.RP.LaunchID
			itemID := sessionRQ.Capabilities.BrowserKubeOpts.RP.ItemID
			project := sessionRQ.Capabilities.BrowserKubeOpts.RP.Project
			if project == "" {
				// just skip to next plugin execution
				return next(ctx, prq, sessionRQ, sID)
			}

			log := zap.S().With("project", project)
			log.Info("ReportPortal launch detected")
			settings, err := sr.FindByProjectName(ctx, project)
			if err != nil {
				log.Errorf("ReportPortal isn't properly configured: %s", err)
				return next(ctx, prq, sessionRQ, sID)
			}

			rp := gorp.NewClient(settings.Host, project, settings.AuthToken)
			if launchID == "" {
				u := uuid.New()
				newLaunch, err := rp.StartLaunch(&gorp.StartLaunchRQ{
					StartRQ: gorp.StartRQ{
						Name:        "Browserkube",
						Description: "Browser session: " + sID,
						UUID:        &u,
						StartTime:   gorp.NewTimestamp(time.Now()),
						Attributes:  []*gorp.Attribute{{Parameter: gorp.Parameter{Key: "session", Value: sID}}},
					},
					Mode: gorp.LaunchModes.Default,
				})
				if err != nil {
					log.Errorf("Unable to start launch in ReportPortal: %s", err)
					// do not stop plugins execution chain
					return next(ctx, prq, sessionRQ, sID)
				}
				launchID = newLaunch.ID
				sessionRQ.Capabilities.BrowserKubeOpts.RP.LaunchID = newLaunch.ID
				log.Infow("ReportPortal launch has been created", "launchId", newLaunch.ID)
			}

			if itemID == "" {
				newTest, err := rp.StartTest(&gorp.StartTestRQ{
					Type:       gorp.TestItemTypes.Test,
					UniqueID:   uuid.New().String(),
					LaunchID:   launchID,
					HasStats:   true,
					Retry:      false,
					TestCaseID: uuid.NewString(),
					StartRQ: gorp.StartRQ{
						Name:       fmt.Sprintf("WebDriver session: %s", sID),
						Attributes: []*gorp.Attribute{{Parameter: gorp.Parameter{Key: "session", Value: sID}}},
						StartTime:  gorp.NewTimestamp(time.Now()),
					},
				})
				if err != nil {
					log.Errorf("Unable to create new item in ReportPortal: %s", err)
					// do not stop plugins execution chain
					return next(ctx, prq, sessionRQ, sID)
				}
				sessionRQ.Capabilities.BrowserKubeOpts.RP.ItemID = newTest.ID
				sessionRQ.Capabilities.BrowserKubeOpts.RP.FinishItem = true
				log.Infow("ReportPortal item has been created", "launchId", launchID, "itemID", newTest.ID)
			}

			return next(ctx, prq, sessionRQ, sID)
		}
	}
}

// afterCommandHandler deletes a pod when quit session is requested
func afterCommandHandler(sr settingsRepo) func(next wd.OnAfterCommand) wd.OnAfterCommand {
	return func(next wd.OnAfterCommand) wd.OnAfterCommand {
		return func(ctx *wd.Context, rs *http.Response, sess *session.Session, command string) error {
			if sess.Caps.BrowserKubeOpts.RP == nil {
				// just skip to next plugin execution
				return next(ctx, rs, sess, command)
			}
			launchID := sess.Caps.BrowserKubeOpts.RP.LaunchID
			itemID := sess.Caps.BrowserKubeOpts.RP.ItemID
			project := sess.Caps.BrowserKubeOpts.RP.Project
			if project == "" {
				// just skip to next plugin execution
				return next(ctx, rs, sess, command)
			}

			log := zap.S().With("project", project, "launchId", launchID, "itemId", itemID)
			if launchID == "" || itemID == "" {
				log.Warn("Launch ID or item ID isn't found")
				return next(ctx, rs, sess, command)
			}

			settings, err := sr.FindByProjectName(ctx, project)
			if err != nil {
				log.Errorf("ReportPortal isn't properly configured: %s", err)
				return next(ctx, rs, sess, command)
			}
			rp := gorp.NewClient(settings.Host, project, settings.AuthToken)

			if _, err = rp.SaveLog(&gorp.SaveLogRQ{
				ItemID:     itemID,
				LaunchUUID: launchID,
				Level:      gorp.LogLevelInfo,
				LogTime:    gorp.NewTimestamp(time.Now()),
				Message:    fmt.Sprintf("WebDriver command: %s %s", rs.Request.Method, command),
			}); err != nil {
				log.Errorf("Unable to send log to ReportPortal: %s", err)
				return next(ctx, rs, sess, command)
			}

			log.Infof("Log has been reported for ReportPortal")
			return next(ctx, rs, sess, command)
		}
	}
}

// findElementHandler handles findElement webdriver commands. In case when element not found it takes a
// screenshot of current windows and attaches to log in ReportPortal
func findElementHandler(sr settingsRepo) func(next wd.OnAfterCommand) wd.OnAfterCommand {
	return func(next wd.OnAfterCommand) wd.OnAfterCommand {
		return func(ctx *wd.Context, rs *http.Response, sess *session.Session, command string) error {
			if !(strings.HasSuffix(command, "/element") || strings.HasSuffix(command, "/elements")) {
				return next(ctx, rs, sess, command)
			}
			if rs.StatusCode != http.StatusNotFound {
				// ignore any status code other than 404. We need to handle only 'no such element' error.
				return next(ctx, rs, sess, command)
			}

			log := zap.S()
			screenshotBytes, err := wdproto.NewWebDriver(sess.Browser.Status.SeleniumURL, sess.ID).TakeScreenshot(ctx)
			if err != nil {
				log.Errorf("unable to take a screenshot: %s", err)
				return next(ctx, rs, sess, command)
			}

			project := sess.Caps.BrowserKubeOpts.RP.Project
			launchID := sess.Caps.BrowserKubeOpts.RP.LaunchID
			itemID := sess.Caps.BrowserKubeOpts.RP.ItemID
			log = zap.S().With("project", project, "launchId", launchID, "itemId", itemID)
			settings, err := sr.FindByProjectName(ctx, project)
			if err != nil {
				log.Errorf("ReportPortal isn't properly configured: %s", err)
				return next(ctx, rs, sess, command)
			}

			fileName := uuid.New().String() + ".png"
			_, err = gorp.NewClient(settings.Host, project, settings.AuthToken).SaveLogMultipart(
				[]*gorp.SaveLogRQ{
					{
						LaunchUUID: launchID,
						ItemID:     itemID,
						Level:      gorp.LogLevelError,
						LogTime:    gorp.NewTimestamp(time.Now()),
						Message: fmt.Sprintf(
							"Element not found for item '%s', command: %s %s",
							itemID, rs.Request.Method, command,
						),
						Attachment: gorp.Attachment{Name: fileName},
					},
				},
				[]gorp.Multipart{
					&gorp.ReaderMultipart{FileName: fileName, Reader: bytes.NewReader(screenshotBytes)},
				},
			)
			if err != nil {
				log.Errorf("unable to send log with attachment to ReportPortal: %s", err)
				return next(ctx, rs, sess, command)
			}

			return next(ctx, rs, sess, command)
		}
	}
}
