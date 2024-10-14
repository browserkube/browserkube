package screenshot

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/browserkube/browserkube/browserkube/internal/api"
	"github.com/browserkube/browserkube/pkg/session"
	"github.com/browserkube/browserkube/pkg/wd"
	"github.com/browserkube/browserkube/pkg/wd/wdproto"
	"github.com/browserkube/browserkube/storage"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			provideScreenshotCapturePlugin,
			fx.ResultTags(`group:"wd-extensions"`),
		),
		fx.Annotate(
			provideScreenshotOnNotFoundPlugin,
			fx.ResultTags(`group:"wd-extensions"`),
		),
	),
)

func provideScreenshotCapturePlugin(store storage.BlobSessionStorage) wd.PluginOpts {
	return wd.PluginOpts{
		Weight: 250,
		Opts: []wd.PluginOpt{
			wd.WithAfterCommand(screenshotCapture(store)), //nolint:bodyclose
		},
	}
}
func provideScreenshotOnNotFoundPlugin(store storage.BlobSessionStorage) wd.PluginOpts {
	return wd.PluginOpts{
		Weight: 250,
		Opts: []wd.PluginOpt{
			wd.WithAfterCommand(screenshotIfNotFound(store)), //nolint:bodyclose
		},
	}
}

func screenshotCapture(store storage.BlobSessionStorage) func(next wd.OnAfterCommand) wd.OnAfterCommand {
	return func(next wd.OnAfterCommand) wd.OnAfterCommand {
		return func(ctx *wd.Context, rs *http.Response, sess *session.Session, command string) error {
			log := zap.S().With("sessionId", sess.ID)

			if command != "/screenshot" {
				return next(ctx, rs, sess, command)
			}

			rsPayload := &bytes.Buffer{}
			if _, err := io.Copy(rsPayload, rs.Body); err != nil {
				log.Errorf("failed to copy into rsPayload: %v", err)
				return next(ctx, rs, sess, command)
			}
			rs.Body = io.NopCloser(rsPayload)
			fileName := time.Now().UTC().Format("2006-01-02-15-04-05") + "-auto-screenshot.png"

			if err := store.SaveFile(ctx, sess.ID, api.ScreenshotsPath, &storage.BlobFile{
				FileName:    fileName,
				ContentType: "image/png",
				Content:     bytes.NewReader(rsPayload.Bytes()),
			}); err != nil {
				log.Errorf("failed to save sessionRecord: %v", err)
				return next(ctx, rs, sess, command)
			}

			log.Info("Screenshot has been saved: ", fileName)

			return next(ctx, rs, sess, command)
		}
	}
}

// TODO: cases need to be improved when automatic screenshots are required
func screenshotIfNotFound(store storage.BlobSessionStorage) func(next wd.OnAfterCommand) wd.OnAfterCommand {
	return func(next wd.OnAfterCommand) wd.OnAfterCommand {
		return func(ctx *wd.Context, rs *http.Response, sess *session.Session, command string) error {
			log := zap.S().With("sessionId", sess.ID)

			if rs.StatusCode != http.StatusNotFound {
				// ignore any status code other than 404. We need to handle only 'no such element' error.
				log.Errorf("Status code: %v", rs.StatusCode)

				return next(ctx, rs, sess, command)
			}

			screenshotBytes, err := wdproto.NewWebDriver(sess.Browser.Status.SeleniumURL, sess.ID).TakeScreenshot(ctx)
			if err != nil {
				log.Errorf("unable to take a screenshot: %s", err)
				return next(ctx, rs, sess, command)
			}

			var buf bytes.Buffer

			if err := json.NewEncoder(&buf).Encode(screenshotBytes); err != nil {
				log.Errorf("failed to encode screenshotBytes: %v", err)
				return next(ctx, rs, sess, command)
			}

			fileName := time.Now().UTC().Format("2006-01-02-15-04-05") + "-auto-screenshot.png"

			if err := store.SaveFile(ctx, sess.ID, api.ScreenshotsPath, &storage.BlobFile{
				FileName:    fileName,
				ContentType: "image/png",
				Content:     &buf,
			}); err != nil {
				log.Errorf("failed to save sessionRecord: %v", err)
				return next(ctx, rs, sess, command)
			}

			log.Info("Screenshot has been saved: ", fileName)

			return next(ctx, rs, sess, command)
		}
	}
}
