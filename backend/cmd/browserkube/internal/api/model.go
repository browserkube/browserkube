package api

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

var defaultResolutions = []string{
	"1024x768",
	"1280x960",
	"1280x1024",
	"1600x1200",
	"1440x900",
	"1920x1080",
}

type (
	StatusStats struct {
		All        int `json:"all"`
		Running    int `json:"running"`
		Connecting int `json:"connecting"`
		Queued     int `json:"queued"`
	}

	Browser struct {
		Platform    string   `json:"platformName"`
		Name        string   `json:"name"`
		Version     string   `json:"version"`
		Image       string   `json:"image"`
		Type        string   `json:"type"`
		Resolutions []string `json:"resolutions"`
	}
	Status struct {
		QuotesLimit int           `json:"quotesLimit"`
		MaxTimeout  time.Duration `json:"maxTimeout"`
		Stats       StatusStats   `json:"stats"`
	}

	WSMessage[T any] struct {
		Name    string `json:"name,omitempty"`
		Payload T      `json:"payload,omitempty"`
	}

	//nolint:maligned
	Session struct {
		ID               string                 `json:"id,omitempty"`
		Name             string                 `json:"name,omitempty"`
		Image            string                 `json:"image,omitempty"`
		Type             string                 `json:"type,omitempty"`
		State            string                 `json:"state,omitempty"`
		Details          map[string]interface{} `json:"details,omitempty"`
		Platform         string                 `json:"platformName,omitempty"`
		Manual           bool                   `json:"manual,omitempty"`
		Browser          string                 `json:"browser,omitempty"`
		BrowserVersion   string                 `json:"browserVersion,omitempty"`
		VncOn            bool                   `json:"vncOn,omitempty"`
		VncPsw           string                 `json:"vncPsw,omitempty"`
		LogsOn           bool                   `json:"logsOn,omitempty"`
		LogsRefAddr      string                 `json:"logsRefAddr,omitempty"`
		VideoRecOn       bool                   `json:"videoRecOn,omitempty"`
		VideoCodec       string                 `json:"videoCodec,omitempty"`
		VideoName        string                 `json:"videoName,omitempty"`
		VideoFps         int                    `json:"videoFps,omitempty"`
		VideoRefAddr     string                 `json:"videoRefAddr,omitempty"`
		VideoSize        *Resolution            `json:"videoSize,omitempty"`
		ScreenResolution string                 `json:"screenResolution,omitempty"`
		CreatedAt        Timestamp              `json:"createdAt,omitempty"`
	}
	Timestamp time.Time

	Resolution struct {
		Height int `json:"height"`
		Width  int `json:"width"`
	}

	SessionResult struct {
		Session
	}

	CommandLog struct {
		SessionID  string    `json:"sessionId"`
		CommandID  string    `json:"commandId"`
		Method     string    `json:"method"`
		Command    string    `json:"command"`
		Request    []byte    `json:"request"`
		StatusCode int       `json:"statusCode"`
		Response   []byte    `json:"response"`
		Timestamp  time.Time `json:"timestamp"`
	}

	CommandLogResponse struct {
		Commands     []CommandLog `json:"commands"`
		NewPageToken string       `json:"newPageToken"`
	}
)

func (t Timestamp) MarshalJSON() ([]byte, error) {
	tt := time.Time(t).UnixMilli()
	bytes, err := json.Marshal(tt)
	return bytes, errors.WithStack(err)
}

func NewWSMessage[T any](name string, payload T) *WSMessage[T] {
	return &WSMessage[T]{Name: name, Payload: payload}
}
