package wdproto

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/browserkube/browserkube/pkg/session"
	browserkubeutil "github.com/browserkube/browserkube/pkg/util"
)

type Response struct {
	Value interface{} `json:"value"`
}

type Error struct {
	Error      string
	Message    string
	Stacktrace string
}

//go:generate easyjson
//easyjson:json
type NewSessionRQ struct {
	Capabilities    session.Capabilities `json:"desiredCapabilities,omitempty"`
	W3CCapabilities W3CCapabilities      `json:"capabilities,omitempty"`
}

type W3CCapabilities struct {
	Capabilities session.Capabilities    `json:"alwaysMatch,omitempty"`
	FirstMatch   []*session.Capabilities `json:"firstMatch,omitempty"`
}

type NewSessionRS struct {
	Value NewSessionRSValue `json:"value"`
}

type NewSessionRSValue struct {
	SessionID    string                 `json:"sessionId"`
	Capabilities map[string]interface{} `json:"capabilities"`
}

type CreateBrowserRequest struct {
	SessionName    string `json:"sessionName"`
	Platform       string `json:"platformName"`
	BrowserName    string `json:"browserName"`
	BrowserVersion string `json:"browserVersion"`
	Resolution     string `json:"screenResolution,omitempty"`
	RecordVideo    bool   `json:"recordVideo,omitempty"`
}
type WebDriver struct {
	client    *http.Client
	baseURL   string
	sessionID string
}

func NewWebDriver(baseURL, sessionID string) *WebDriver {
	return &WebDriver{baseURL: baseURL, sessionID: sessionID, client: http.DefaultClient}
}

func (wd *WebDriver) Maximize(ctx context.Context) error {
	err := wd.executeCommand(ctx, http.MethodPost, "/session/%s/window/maximize", nil, nil)
	return errors.WithStack(err)
}

func (wd *WebDriver) TakeScreenshot(ctx context.Context) ([]byte, error) {
	imgBase64String, err := newCommandExecutor[string](wd).Do(ctx, http.MethodGet, "/session/%s/screenshot", nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	imgBytes, err := base64.StdEncoding.DecodeString(imgBase64String)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return imgBytes, nil
}

func (wd *WebDriver) CurrentWindowHandle(ctx context.Context) (string, error) {
	h, err := newCommandExecutor[string](wd).Do(ctx, http.MethodPost, "/session/%s/window", nil)
	return h, errors.WithStack(err)
}

func (wd *WebDriver) Quit(ctx context.Context) error {
	var u string
	var err error
	if !strings.Contains(wd.baseURL, "session") {
		u, err = url.JoinPath(wd.baseURL, "session")
		if err != nil {
			return errors.WithStack(err)
		}
	}
	u, err = url.JoinPath(u, wd.sessionID)
	if err != nil {
		return errors.WithStack(err)
	}

	rq, err := http.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return errors.Wrap(err, "Unable to build quit session request")
	}
	resp, err := wd.client.Do(rq.WithContext(ctx)) //nolint:bodyclose
	if resp != nil {
		defer browserkubeutil.CloseQuietly(resp.Body)
	}

	if err == nil && resp.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("delete request failed: %w", err)
}

type commandExecutor[T any] interface {
	Do(ctx context.Context, method, command string, body interface{}) (T, error)
}
type commandExecutorFunc[T any] func(ctx context.Context, method, command string, body interface{}) (T, error)

func (f commandExecutorFunc[T]) Do(ctx context.Context, method, command string, body interface{}) (T, error) {
	return f(ctx, method, command, body)
}

func newCommandExecutor[T any](wd *WebDriver) commandExecutor[T] {
	return commandExecutorFunc[T](func(ctx context.Context, method, command string, body interface{}) (T, error) {
		resultPayload := struct {
			Value T `json:"value"`
		}{}
		if err := wd.executeCommand(ctx, method, command, body, &resultPayload); err != nil {
			return *new(T), errors.WithStack(err)
		}
		return resultPayload.Value, nil
	})
}

func (wd *WebDriver) executeCommand(ctx context.Context, method, command string, body, res interface{}) error {
	cmdURL, err := url.JoinPath(wd.baseURL, fmt.Sprintf(command, wd.sessionID))
	if err != nil {
		return errors.WithStack(err)
	}

	if body == nil {
		body = map[string]string{}
	}
	buf := &bytes.Buffer{}
	if err = json.NewEncoder(buf).Encode(body); err != nil {
		return errors.WithStack(err)
	}

	rq, err := http.NewRequestWithContext(ctx, method, cmdURL, buf)
	if err != nil {
		return errors.WithStack(err)
	}
	rs, err := wd.client.Do(rq) //nolint:bodyclose
	if err != nil {
		return errors.WithStack(err)
	}
	defer browserkubeutil.CloseQuietly(rs.Body)

	if rs.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(rs.Body)
		zap.S().Errorf("webdriver command execution error: %s", string(b))
		return errors.Errorf("unable to execute command. Status code: %d", rs.StatusCode)
	}

	if res != nil {
		if decErr := json.NewDecoder(rs.Body).Decode(res); decErr != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func BadGatewayError(w http.ResponseWriter, err error) {
	logger := zap.S()
	logger.Errorf("Bad Gateway Error: %+v", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadGateway)
	if encErr := json.NewEncoder(w).Encode(&Response{
		Value: Error{
			Error:   err.Error(),
			Message: "something went wrong",
		},
	}); encErr != nil {
		logger.Errorf("Bad gateway: %v", encErr)
	}
}
