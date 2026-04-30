package session

import (
	"encoding/json"
	"encoding/json/jsontext"

	jsonv2 "encoding/json/v2"

	"go.uber.org/zap"

	browserkubev1 "github.com/browserkube/browserkube/operator/api/v1"
)

type Session struct {
	ID      string
	State   string
	Browser *browserkubev1.Browser
	Caps    *Capabilities
}

type Capabilities struct {
	Extra           jsontext.Value  `json:",inline"`
	Platform        string          `json:"platformName,omitzero"`
	BrowserVersion  string          `json:"browserVersion,omitzero"`
	BrowserName     string          `json:"browserName,omitzero"`
	Timezone        string          `json:"timeZone,omitzero"`
	BrowserKubeOpts BrowserKubeOpts `json:"browserkube:options,omitzero"`
}

func (v Capabilities) MarshalJSON() ([]byte, error) {
	type alias Capabilities
	return jsonv2.Marshal(alias(v))
}

func (v *Capabilities) UnmarshalJSON(data []byte) error {
	type alias Capabilities
	return jsonv2.Unmarshal(data, (*alias)(v))
}

//nolint:maligned
type BrowserKubeOpts struct {
	Extra              jsontext.Value                   `json:",inline"`
	RP                 *ReportPortalOpts                `json:"reportportal,omitzero"       schema:"-"`
	User               string                           `json:"user,omitzero"               schema:"-"`
	Token              string                           `json:"token,omitzero"              schema:"-"`
	Name               string                           `json:"name,omitzero"               schema:"-"`
	VideoFileName      string                           `json:"videoFileName,omitzero"      schema:"-"`
	Type               string                           `json:"type,omitzero"               schema:"-"`
	Manual             bool                             `json:"manual,omitzero"             schema:"-"`
	EnableVideo        bool                             `json:"enableVideo,omitzero"        schema:"enableVideo"`
	ScreenResolution   string                           `json:"screenResolution,omitzero"   schema:"screenResolution"`
	SessionTimeout     int                              `json:"sessionTimeout,omitzero"     schema:"sessionTimeout"`
	SessionIdleTimeout int                              `json:"sessionIdleTimeout,omitzero" schema:"sessionIdleTimeout"`
	EnableVNC          bool                             `json:"enableVNC,omitzero"          schema:"enableVNC"` //nolint:tagliatelle
	Extensions         []browserkubev1.BrowserExtension `json:"extensions,omitempty"        schema:"-"`
}

func (v BrowserKubeOpts) MarshalJSON() ([]byte, error) {
	type alias BrowserKubeOpts
	return jsonv2.Marshal(alias(v))
}

func (v *BrowserKubeOpts) UnmarshalJSON(data []byte) error {
	type alias BrowserKubeOpts
	return jsonv2.Unmarshal(data, (*alias)(v))
}

type ReportPortalOpts struct {
	Project    string `json:"project,omitzero"`
	LaunchID   string `json:"launchId,omitzero"`
	ItemID     string `json:"itemId,omitzero"`
	FinishItem bool   `json:"finishItem,omitzero"`
}

func (v *Capabilities) String() string {
	str, err := json.Marshal(v)
	if err != nil {
		zap.S().Error("unable to marshall caps", err)
		return ""
	}
	return string(str)
}
