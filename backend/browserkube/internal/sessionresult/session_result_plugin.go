package sessionresult

import (
	"context"
	"path"

	"go.uber.org/fx"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	browserkubev1 "github.com/browserkube/browserkube/operator/api/v1"
	"github.com/browserkube/browserkube/pkg/session"
	"github.com/browserkube/browserkube/pkg/sessionresult"
	"github.com/browserkube/browserkube/pkg/wd"
	"github.com/browserkube/browserkube/storage"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			provideSessionResultPlugin,
			fx.ResultTags(`group:"wd-extensions"`),
		),
	),
)

func provideSessionResultPlugin(sessionResultsRepo sessionresult.Repository, store storage.BlobSessionStorage) wd.PluginOpts {
	return wd.PluginOpts{
		Weight: 1,
		Opts: []wd.PluginOpt{
			wd.WithQuitSession(quitSessionHandler(sessionResultsRepo, store)),
		},
	}
}

func sessionFileExists(storage storage.BlobSessionStorage, fileName, sessionID string) bool {
	exists, err := storage.Exists(context.Background(), sessionID, fileName)
	return err == nil && exists
}

// quitSessionHandler creates session result on session quit
func quitSessionHandler(sessionResultsRepo sessionresult.Repository, store storage.BlobSessionStorage) func(next wd.OnSessionQuit) wd.OnSessionQuit {
	return func(next wd.OnSessionQuit) wd.OnSessionQuit {
		return func(ctx *wd.Context, s *session.Session) error {
			log := zap.S().With("sessionId", s.ID)

			sr := &sessionresult.Result{
				SessionResult: browserkubev1.SessionResult{
					ObjectMeta: metav1.ObjectMeta{
						Name:        s.Browser.Name,
						Namespace:   s.Browser.Namespace,
						Labels:      s.Browser.Labels,
						Annotations: s.Browser.Annotations,
					},
					Spec: browserkubev1.SessionResultSpec{
						StartedAt:    s.Browser.CreationTimestamp,
						Browser:      s.Browser.Spec,
						BrowserImage: s.Browser.Status.Image,
						Files:        browserkubev1.SessionResultFiles{},
					},
				},
			}
			if s.Browser.DeletionTimestamp != nil {
				sr.SessionResult.Spec.FinishedAt = *s.Browser.DeletionTimestamp
			}

			if s.Browser.Spec.EnableVideo && sessionFileExists(store, sessionresult.VideoFileName, s.ID) {
				sr.Spec.Files.Video = path.Join(s.ID, sessionresult.VideoFileName)
			}

			if sessionFileExists(store, sessionresult.BrowserLogFileName, s.ID) {
				sr.Spec.Files.BrowserLog = path.Join(s.ID, sessionresult.BrowserLogFileName)
			}

			_, err := sessionResultsRepo.Create(ctx, sr)
			if err != nil {
				log.Error("unable to create session result", err)
				return next(ctx, s)
			}

			log.Info("saved session result")
			return next(ctx, s)
		}
	}
}
