// This file is licensed under the Apache Software License, v.2.0 except as
// noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SteakfulSetKind is the kind of a SteakfulSet
const SteakfulSetKind = "SteakfulSet"

// SteakfulSetSpec describes the desired state of a SteakfulSet
type SteakfulSetSpec struct {
	// Guests defines the number of guests we expect for this SteakfulSet
	Guests int `json:"guests"`

	// Steak defines the desired state of the Steaks to be prepared by this SteakfulSet
	Steak Steak `json:"steak"`
}

// SteakfulSetStatus describes the observed state of a SteakfulSet
type SteakfulSetStatus struct {
	SteaksServed []corev1.ObjectReference `json:"steaksServed"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=stks
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:JSONPath=.spec.guests,type=integer,name=Guests,description="The number of guests to be served, aka replicaCount."
//+kubebuilder:printcolumn:JSONPath=.spec.steak.spec.weight,type=integer,name=Weight,description="The weight of the Steaks to be served."
//+kubebuilder:printcolumn:JSONPath=.spec.steak.spec.cookLevel,type=string,name=Cooked,description="The level of how thorough the Steak should be cooked."
//+kubebuilder:printcolumn:JSONPath=.spec.steak.spec.fat,type=string,name=Fat,description="The juicyness of the Steaks to be served."
//+kubebuilder:printcolumn:JSONPath=.spec.steak.spec.variant,type=string,name=Variant,description="The kind of Steak to be served."
//+kubebuilder:printcolumn:JSONPath=.metadata.creationTimestamp,type=date,name=Age

// SteakfulSet is the Schema for the steakfulsets API
type SteakfulSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SteakfulSetSpec   `json:"spec,omitempty"`
	Status SteakfulSetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SteakfulSetList contains a list of SteakfulSets
type SteakfulSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SteakfulSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SteakfulSet{}, &SteakfulSetList{})
}
