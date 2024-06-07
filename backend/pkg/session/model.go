package session

import (
	"encoding/json"

	"github.com/mailru/easyjson"
	"go.uber.org/zap"

	browserkubev1 "github.com/browserkube/browserkube/operator/api/v1"
)

type Session struct {
	ID      string
	State   string
	Browser *browserkubev1.Browser
	Caps    *Capabilities
}

//go:generate easyjson
//easyjson:json
type Capabilities struct {
	easyjson.UnknownFieldsProxy
	Platform        string          `json:"platformName,omitempty"`
	BrowserVersion  string          `json:"browserVersion,omitempty"`
	BrowserName     string          `json:"browserName,omitempty"`
	Timezone        string          `json:"timeZone,omitempty"`
	BrowserKubeOpts BrowserKubeOpts `json:"browserkube:options,omitempty"`
}

//go:generate easyjson
//easyjson:json
//nolint:maligned
type BrowserKubeOpts struct {
	easyjson.UnknownFieldsProxy
	RP               *ReportPortalOpts `json:"reportportal,omitempty"     schema:"-"`
	User             string            `json:"user,omitempty"             schema:"-"`
	Token            string            `json:"token,omitempty"            schema:"-"`
	Name             string            `json:"name,omitempty"             schema:"-"`
	VideoFileName    string            `json:"videoFileName,omitempty"    schema:"-"`
	Type             string            `json:"type,omitempty"             schema:"-"`
	Manual           bool              `json:"manual,omitempty"           schema:"-"`
	EnableVideo      bool              `json:"enableVideo,omitempty"      schema:"enableVideo"`
	ScreenResolution string            `json:"screenResolution,omitempty" schema:"screenResolution"`

	//nolint: tagliatelle
	EnableVNC  bool                             `json:"enableVNC,omitempty"  schema:"enableVNC"`
	Extensions []browserkubev1.BrowserExtension `json:"extensions,omitempty" schema:"-"`
}

type ReportPortalOpts struct {
	Project    string `json:"project,omitempty"`
	LaunchID   string `json:"launchId,omitempty"`
	ItemID     string `json:"itemId,omitempty"`
	FinishItem bool   `json:"finishItem,omitempty"`
}

func (v *Capabilities) String() string {
	str, err := json.Marshal(v)
	if err != nil {
		zap.S().Error("unable to marshall caps", err)
		return ""
	}
	return string(str)
}
