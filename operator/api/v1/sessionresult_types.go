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

// SessionResultSpec defines the desired state of SessionResult
type SessionResultSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	StartedAt    metav1.Time        `json:"startedAt,omitempty"`
	FinishedAt   metav1.Time        `json:"finishedAt,omitempty"`
	Browser      BrowserSpec        `json:"browser,omitempty"`
	BrowserImage string             `json:"browserImage,omitempty"`
	Files        SessionResultFiles `json:"files"`
}

type SessionResultFiles struct {
	BrowserLog string `json:"browserLog"`
	Video      string `json:"video"`
	Bookmarks  string `json:"bookmarks"`
}

// Status defines the observed state of Session
type Status string

// SessionResultStatus defines the observed state of SessionResult
type SessionResultStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SessionResult is the Schema for the sessionresults API
type SessionResult struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SessionResultSpec   `json:"spec,omitempty"`
	Status SessionResultStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SessionResultList contains a list of SessionResult
type SessionResultList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SessionResult `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SessionResult{}, &SessionResultList{})
}
