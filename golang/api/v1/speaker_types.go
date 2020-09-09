/*


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

// SpeakerSpec defines the desired state of Speaker
type SpeakerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// NOTE Custom spec fields

	// Speaker First name
	FirstName string `json:"firstName,omitempty"`
	// Speaker Last name
	LastName string `json:"lastName,omitempty"`
	// Speaker Avatar url
	Avatar string `json:"avatar,omitempty"`
}

// SpeakerStatus defines the observed state of Speaker
type SpeakerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// NOTE Custom status fields

	// Number of sessions made by the speaker
	Sessions string `json:"sessions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Speaker is the Schema for the speakers API
type Speaker struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SpeakerSpec   `json:"spec,omitempty"`
	Status SpeakerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SpeakerList contains a list of Speaker
type SpeakerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Speaker `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Speaker{}, &SpeakerList{})
}
