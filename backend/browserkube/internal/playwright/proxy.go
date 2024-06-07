package playwright

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	_ "gocloud.dev/blob/fileblob" // Register some standard stuff
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	browserkubev1 "github.com/browserkube/browserkube/operator/api/v1"
	"github.com/browserkube/browserkube/pkg/sessionresult"
	"github.com/browserkube/browserkube/pkg/websocketproxy"
	"github.com/browserkube/browserkube/storage"
)

const (
	ScreenshotDir = "/screenshots"
)

type Proxy struct {
	*websocketproxy.WebsocketProxy
	LogBuffer       bytes.Buffer
	SessionRecorder storage.BlobSessionStorage
	screenshoter    *Screenshoter
	ctx             context.Context //nolint: containedctx
}

func NewProxy(ctx context.Context, u *url.URL, sessionID string, sessionRecorder storage.BlobSessionStorage) (*Proxy, error) {
	dirPath := path.Join(ScreenshotDir, sessionID)

	s := &Screenshoter{
		requested: map[int]any{},
		lastID:    10000,
		mu:        sync.RWMutex{},
		dirPath:   dirPath,
	}

	pp := &Proxy{
		LogBuffer:       bytes.Buffer{},
		screenshoter:    s,
		ctx:             ctx,
		SessionRecorder: sessionRecorder,
	}

	opts := []websocketproxy.ProxyOpt{
		websocketproxy.WithOutgoingMiddleware(pp.record),
		websocketproxy.WithIncomingMiddleware(pp.record),
		websocketproxy.WithOutgoingMiddleware(pp.screenshotRecord),
		websocketproxy.WithOutgoingMiddleware(pp.takeScreenshot),
	}

	proxy, err := websocketproxy.NewProxy(u, opts...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pp.WebsocketProxy = proxy
	return pp, nil
}

func (pp *Proxy) record(msg *websocketproxy.Message) error {
	_, span := otel.Tracer("").Start(pp.ctx, msg.Method)
	defer span.End()

	if strings.Contains(msg.GUID, "page@") {
		pp.screenshoter.page = msg.GUID
	}

	newMsg, err := json.Marshal(&msg)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = pp.LogBuffer.Write(newMsg)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

type Screenshoter struct {
	requested map[int]any
	page      string

	lastID  int
	dirPath string

	mu sync.RWMutex
}

func (pp *Proxy) takeScreenshot(msg *websocketproxy.Message) error {
	if msg.Result == nil && pp.screenshoter.page == "" {
		return nil
	}
	if data, ok := msg.Result.(map[string]interface{}); ok {
		if matchesValue, exists := data["matches"]; exists {
			if matchesBool, ok := matchesValue.(bool); ok {
				if !matchesBool {
					return nil
				}
				pp.screenshoter.mu.Lock()
				defer pp.screenshoter.mu.Unlock()

				id := pp.screenshoter.lastID
				pp.screenshoter.lastID--

				msg := &websocketproxy.Message{
					ID:     id,
					GUID:   pp.screenshoter.page,
					Method: "screenshot",
					Params: map[string]interface{}{
						"type": "png",
					},
					Metadata: map[string]interface{}{
						"wallTime": time.Now().Nanosecond(),
						"apiName":  "page.screenshot",
						"internal": false,
					},
				}

				newMsg, err := json.Marshal(&msg)
				if err != nil {
					return errors.WithStack(err)
				}

				err = pp.WSConn.WriteMessage(websocket.TextMessage, newMsg)
				if err != nil {
					return errors.WithStack(err)
				}
				pp.screenshoter.requested[id] = struct{}{}
			}
		}
	}

	return nil
}

func (pp *Proxy) screenshotRecord(msg *websocketproxy.Message) error {
	pp.screenshoter.mu.RLock()
	defer pp.screenshoter.mu.RUnlock()
	if _, found := pp.screenshoter.requested[msg.ID]; !found || msg.Result == nil {
		return nil
	}

	if data, ok := msg.Result.(map[string]interface{}); ok {
		if binaryValue, exists := data["binary"]; exists {
			if binaryStr, ok := binaryValue.(string); ok {
				imgBytes, err := base64.StdEncoding.DecodeString(binaryStr)
				if err != nil {
					return websocketproxy.ErrDoNotSend
				}

				if err = os.MkdirAll(pp.screenshoter.dirPath, 0o777); err != nil {
					return websocketproxy.ErrDoNotSend
				}

				// Time formatting as YYYYMMDD_HHMMSS
				currentTime := time.Now().Format("20060102-150405")
				fileName := fmt.Sprintf("screenshot-%s", currentTime)

				f, err := os.CreateTemp(pp.screenshoter.dirPath, fileName)
				if err != nil {
					return websocketproxy.ErrDoNotSend
				}
				defer f.Close()

				if err = json.NewEncoder(f).Encode(imgBytes); err != nil {
					return websocketproxy.ErrDoNotSend
				}
				return websocketproxy.ErrDoNotSend
			}
		}
	}

	return nil
}

func (pp *Proxy) SaveSessionRecord(ctx context.Context, sessionID string) error {
	err := pp.SessionRecorder.SaveFile(ctx, sessionID, "", &storage.BlobFile{
		FileName:    sessionresult.MessageLogFileName,
		ContentType: "text/plain",
		Content:     &pp.LogBuffer,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (pp *Proxy) SaveScreenshotRecord(ctx context.Context, sessionID string) error {
	if _, err := os.Stat(pp.screenshoter.dirPath); os.IsNotExist(err) {
		return nil
	}

	files, err := os.ReadDir(pp.screenshoter.dirPath)
	if err != nil {
		return errors.WithStack(err)
	}

	var buf bytes.Buffer
	for _, file := range files {
		filePath := path.Join(pp.screenshoter.dirPath, file.Name())

		zap.S().Info("screenshot file name: ", file.Name())

		data, fErr := os.ReadFile(filePath)
		if fErr != nil {
			return errors.WithStack(err)
		}
		buf.Write(data)
		defer buf.Reset()
		err = pp.SessionRecorder.SaveFile(ctx, sessionID, "", &storage.BlobFile{
			FileName:    file.Name(),
			ContentType: "image/png",
			Content:     &buf,
		})
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (pp *Proxy) SaveBrowserLogRecord(ctx context.Context, sessionID string, browserLog io.ReadCloser) error {
	err := pp.SessionRecorder.SaveFile(ctx, sessionID, "", &storage.BlobFile{
		FileName:    sessionresult.BrowserLogFileName,
		ContentType: "text/plain",
		Content:     browserLog,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (pp *Proxy) SaveSessionResult(
	ctx context.Context,
	sessionID string,
	browser *browserkubev1.Browser,
	repo sessionresult.Repository,
) error {
	sr := &sessionresult.Result{
		SessionResult: browserkubev1.SessionResult{
			ObjectMeta: metav1.ObjectMeta{
				Name:        browser.Name,
				Namespace:   browser.Namespace,
				Labels:      browser.Labels,
				Annotations: browser.Annotations,
			},
			Spec: browserkubev1.SessionResultSpec{
				StartedAt:    browser.CreationTimestamp,
				Browser:      browser.Spec,
				BrowserImage: browser.Status.Image,
				Files:        browserkubev1.SessionResultFiles{},
			},
		},
	}

	if browser.DeletionTimestamp != nil {
		sr.SessionResult.Spec.FinishedAt = *browser.DeletionTimestamp
	}
	if sessionFileExists(ctx, pp.SessionRecorder, sessionresult.BrowserLogFileName, sessionID) {
		sr.Spec.Files.BrowserLog = path.Join(sessionID, sessionresult.BrowserLogFileName)
	}

	if _, err := repo.Create(ctx, sr); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func sessionFileExists(ctx context.Context, store storage.Storage, fileName, sessionID string) bool {
	exists, err := store.Exists(ctx, sessionID, fileName)
	return err == nil && exists
}
