package reportcommand

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/browserkube/browserkube/browserkube/internal/api"
	"github.com/browserkube/browserkube/pkg/session"
	"github.com/browserkube/browserkube/pkg/wd"
	"github.com/browserkube/browserkube/storage"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			provideReportCommandPlugin,
			fx.ResultTags(`group:"wd-extensions"`),
		),
	),
)

func provideReportCommandPlugin(store storage.BlobSessionStorage) wd.PluginOpts {
	return wd.PluginOpts{
		Weight: 250,
		Opts: []wd.PluginOpt{
			wd.WithAfterCommand(fetchCommands(store)), //nolint:bodyclose
		},
	}
}

func fetchCommands(store storage.BlobSessionStorage) func(next wd.OnAfterCommand) wd.OnAfterCommand {
	return func(next wd.OnAfterCommand) wd.OnAfterCommand {
		return func(ctx *wd.Context, rs *http.Response, sess *session.Session, command string) error {
			log := zap.S().With("sessionId", sess.ID)

			opts := []trace.SpanStartOption{
				trace.WithAttributes(
					attribute.Key("sessionID").String(sess.ID),
				),
			}

			_, span := otel.Tracer("").Start(ctx, command, opts...)
			defer span.End()

			commandID := rs.Header.Get("commandID")

			commandLog := api.CommandLog{
				SessionID:  sess.ID,
				CommandID:  commandID,
				Method:     rs.Request.Method,
				Command:    command,
				Timestamp:  time.Now(),
				StatusCode: rs.StatusCode,
			}

			rsPayload := &bytes.Buffer{}
			if _, err := io.Copy(rsPayload, rs.Body); err != nil {
				log.Errorf("failed to copy into rsPayload: %v", err)
				return next(ctx, rs, sess, command)
			}
			rs.Body = io.NopCloser(rsPayload)
			commandLog.Response = rsPayload.Bytes()

			rqPayload := &bytes.Buffer{}
			if _, err := io.Copy(rqPayload, rs.Request.Body); err != nil {
				log.Errorf("failed to copy into rqPayload: %v", err)
				return next(ctx, rs, sess, command)
			}

			rs.Request.Body = io.NopCloser(rqPayload)
			commandLog.Request = rqPayload.Bytes()

			fileName := fmt.Sprintf("%03s.json", commandID)

			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(commandLog); err != nil {
				log.Errorf("failed to encode commandLog: %v", err)
				return next(ctx, rs, sess, command)
			}

			if err := store.SaveFile(ctx, sess.ID, api.CommandsPath, &storage.BlobFile{
				FileName:    fileName,
				ContentType: "application/json",
				Content:     &buf,
			}); err != nil {
				log.Errorf("failed to save sessionRecord: %v", err)
				return next(ctx, rs, sess, command)
			}

			log.Infow("command successfully saved to directory", "method", rs.Request.Method, "command", command)

			return next(ctx, rs, sess, command)
		}
	}
}
