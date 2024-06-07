package reportlog

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/browserkube/browserkube/browserkube/internal/provision"
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

func provideReportLogPlugin(serviceProvider provision.Provisioner, store storage.BlobSessionStorage) wd.PluginOpts {
	return wd.PluginOpts{
		Weight: 250,
		Opts: []wd.PluginOpt{
			wd.WithQuitSession(fetchLogsHook(serviceProvider, store)),
		},
	}
}

func fetchLogsHook(serviceProvider provision.Provisioner, store storage.BlobSessionStorage) func(next wd.OnSessionQuit) wd.OnSessionQuit {
	return func(next wd.OnSessionQuit) wd.OnSessionQuit {
		return func(ctx *wd.Context, s *session.Session) error {
			log := zap.S().With("sessionId", s.ID)

			podLogs, err := serviceProvider.Logs(ctx, s.Browser.Status.PodName, false)
			if err != nil {
				log.Errorf("Unable to get pod logs: %v", err)
				return next(ctx, s)
			}
			defer podLogs.Close()

			if err := store.SaveFile(ctx, s.ID, "", &storage.BlobFile{
				FileName:    sessionresult.BrowserLogFileName,
				ContentType: "text/plain",
				Content:     podLogs,
			}); err != nil {
				log.Errorf("Failed to save sessionRecord: %v", err)
				return next(ctx, s)
			}

			log.Info("session result has been saved")
			return next(ctx, s)
		}
	}
}
