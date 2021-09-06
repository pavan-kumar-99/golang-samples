/*
Copyright 2021.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AcmAutoSyncSpec defines the desired state of AcmAutoSync
type AcmAutoSyncSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of AcmAutoSync. Edit acmautosync_types.go to remove/update
	SecretName string `json:"secretName,omitempty"`
	AcmArn     string `json:"acmArn,omitempty"`
}

// AcmAutoSyncStatus defines the observed state of AcmAutoSync
type AcmAutoSyncStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AcmAutoSync is the Schema for the acmautosyncs API
type AcmAutoSync struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AcmAutoSyncSpec   `json:"spec,omitempty"`
	Status AcmAutoSyncStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AcmAutoSyncList contains a list of AcmAutoSync
type AcmAutoSyncList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AcmAutoSync `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AcmAutoSync{}, &AcmAutoSyncList{})
}
