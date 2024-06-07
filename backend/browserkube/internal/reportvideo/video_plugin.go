package reportvideo

import (
	"context"
	"net/http"
	"net/url"
	"path"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/browserkube/browserkube/pkg/session"
	"github.com/browserkube/browserkube/pkg/sessionresult"
	"github.com/browserkube/browserkube/pkg/wd"
	"github.com/browserkube/browserkube/storage"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			provideReportLogPlugin,
			fx.ResultTags(`group:"wd-extensions"`),
		),
	),
)

func provideReportLogPlugin(client *http.Client, storage storage.BlobSessionStorage) wd.PluginOpts {
	return wd.PluginOpts{
		Weight: 251,
		Opts: []wd.PluginOpt{
			wd.WithQuitSession(fetchVideoHook(client, storage)),
		},
	}
}

func stopVideo(ctx context.Context, baseURL string, httpClient *http.Client) error {
	cli := NewClient(&Config{BaseURL: baseURL}, WithClient(httpClient))
	return cli.Stop(ctx)
}

func fetchVideoHook(client *http.Client, store storage.BlobSessionStorage) func(next wd.OnSessionQuit) wd.OnSessionQuit {
	return func(next wd.OnSessionQuit) wd.OnSessionQuit {
		return func(ctx *wd.Context, s *session.Session) error {
			if !s.Browser.Spec.EnableVideo {
				return next(ctx, s)
			}

			log := zap.S().With("sessionId", s.ID)

			baseURL, err := url.Parse(s.Browser.Status.SeleniumURL)
			if err != nil {
				log.Errorf("unable parse url: %v", err)
				return next(ctx, s)
			}

			baseURL.Path = ""
			err = stopVideo(ctx, baseURL.String(), client)
			if err != nil {
				log.Errorf("unable to stop recording: %v", err)
				return next(ctx, s)
			}

			baseURL.Path = path.Join(sessionresult.VideosPath, sessionresult.VideoFileName)

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL.String(), nil)
			if err != nil {
				log.Errorf("unable to create http request: %v", err)
				return next(ctx, s)
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Errorf("unable to get video from url: %v", err)
				return next(ctx, s)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				log.Errorf("unable to get video from url, status code %v", resp.StatusCode)
				return next(ctx, s)
			}

			if err := store.SaveFile(ctx, s.ID, "", &storage.BlobFile{
				FileName:    sessionresult.VideoFileName,
				ContentType: "video/mp4",
				Content:     resp.Body,
			}); err != nil {
				log.Errorf("failed to save video: %v", err)
				return next(ctx, s)
			}

			log.Infof("video has been saved to blob storage")
			return next(ctx, s)
		}
	}
}
