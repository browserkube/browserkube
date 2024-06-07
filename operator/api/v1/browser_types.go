/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BrowserSpec defines the desired state of Browser
type BrowserSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Platform       string `json:"platformName,omitempty"`
	BrowserVersion string `json:"browserVersion"`
	BrowserName    string `json:"browserName"`
	Type           string `json:"type"`
	Timezone       string `json:"timeZone,omitempty"`

	// nolint: tagliatelle
	EnableVNC bool `json:"enableVNC,omitempty"`
	// video recording options
	EnableVideo      bool               `json:"enableVideo,omitempty"`
	ScreenResolution string             `json:"screenResolution,omitempty"`
	Extensions       []BrowserExtension `json:"extensions,omitempty"`

	// +optional
	Caps []byte `json:"caps,omitempty"`
}

type BrowserExtension struct {
	ExtensionID string `json:"extensionId,omitempty"`
	UpdateURL   string `json:"updateUrl,omitempty"`
	Version     string `json:"version,omitempty"`
}

const (
	TypeWebDriver  = "WEBDRIVER"
	TypePlaywright = "PLAYWRIGHT"
)

// BrowserStatus defines the observed state of Browser
type BrowserStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Phase       Phase      `json:"phase"`
	Reason      Reason     `json:"reason,omitempty"`
	Message     string     `json:"message,omitempty"`
	PodName     string     `json:"podName"`
	Host        string     `json:"host,omitempty"`
	SeleniumURL string     `json:"seleniumURL,omitempty"`
	PortConfig  PortConfig `json:"portConfig,omitempty"`
	Image       string     `json:"image,omitempty"`
	VncPass     string     `json:"vncPass,omitempty"`
}
type Reason string

const (
	ReasonVersionNotSupported  = "Version isn't supported"
	ReasonPlatformNotSupported = "Platform isn't supported"
	ReasonConfigNotFound       = "Browser config isn't found"
	ReasonUnknownSessionType   = "Session type unknown"
	ReasonUnknown              = "Unknown"
)

type PortConfig struct {
	Sidecar    string `json:"sidecar,omitempty"`
	Browser    string `json:"browser,omitempty"`
	FileServer string `json:"fileServer,omitempty"`
	Clipboard  string `json:"clipboard,omitempty"`
	VNC        string `json:"vnc,omitempty"`
	DevTools   string `json:"devTools,omitempty"`
}

type Phase string

var (
	PhasePending    Phase = "Pending"
	PhaseRunning    Phase = "Running"
	PhaseTerminated Phase = "Terminated"
	PhaseFailed     Phase = "Failed"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Browser is the Schema for the browsers API
type Browser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BrowserSpec   `json:"spec,omitempty"`
	Status BrowserStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BrowserList contains a list of Browser
type BrowserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Browser `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Browser{}, &BrowserList{})
}
