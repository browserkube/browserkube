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
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BrowserSetSpec defines the desired state of BrowserSet
type BrowserSetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	DefaultTimezone string `json:"defaultTimezone"`

	// +optional
	PodSpec *BrowserPodSpec `json:"podSpec"`
	// +optional
	WebDriver map[string]BrowsersConfig `json:"webdriver,omitempty"`
	// +optional
	Playwright map[string]BrowsersConfig `json:"playwright,omitempty"`
}

type BrowserPodSpec struct {
	TerminationGracePeriodSeconds *int64               `json:"terminationGracePeriodSeconds,omitempty"`
	ActiveDeadlineSeconds         *int64               `json:"activeDeadlineSeconds,omitempty"`
	DNSPolicy                     corev1.DNSPolicy     `json:"dnsPolicy,omitempty"`
	NodeSelector                  map[string]string    `json:"nodeSelector,omitempty"`
	ServiceAccountName            string               `json:"serviceAccountName,omitempty"`
	NodeName                      string               `json:"nodeName,omitempty"`
	Affinity                      *corev1.Affinity     `json:"affinity,omitempty"`
	SchedulerName                 string               `json:"schedulerName,omitempty"`
	Tolerations                   []corev1.Toleration  `json:"tolerations,omitempty"`
	HostAliases                   []corev1.HostAlias   `json:"hostAliases,omitempty"`
	PriorityClassName             string               `json:"priorityClassName,omitempty"`
	Priority                      *int32               `json:"priority,omitempty"`
	DNSConfig                     *corev1.PodDNSConfig `json:"dnsConfig,omitempty"`
}

type BrowsersConfig struct {
	DefaultVersion string `json:"defaultVersion"`
	// +optional
	DefaultPath string                   `json:"defaultPath"`
	Versions    map[string]BrowserConfig `json:"versions"`
}

type BrowserConfig struct {
	Provider string `json:"provider"`
	Image    string `json:"image"`
	Port     string `json:"port"`
	// +optional
	Path string `json:"path"`
	// +optional
	Timezone string `json:"timezone"`
	// +optional
	Spec *BrowserPodSpec `json:"spec"`

	// +optional
	EnableVideo bool `json:"enableVideo,omitempty"`
	// +optional
	AwsAccessKeyID string `json:"awsAccessKeyID,omitempty"`
	// +optional
	AwsSecretAccessKey string `json:"awsSecretAccessKey,omitempty"`

	// TODO not supported yet
	// +optional
	StartupTimeout time.Duration `json:"startupTimeout"`
	// +optional
	SessionDeleteTimeout time.Duration `json:"sessionDeleteTimeout"`
}

// BrowserSetStatus defines the observed state of BrowserSet
type BrowserSetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// BrowserSet is the Schema for the browsersets API
type BrowserSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BrowserSetSpec   `json:"spec,omitempty"`
	Status BrowserSetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BrowserSetList contains a list of BrowserSet
type BrowserSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BrowserSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BrowserSet{}, &BrowserSetList{})
}
