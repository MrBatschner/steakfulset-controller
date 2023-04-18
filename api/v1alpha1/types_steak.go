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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SteakfulSetKind is the kind of a SteakfulSet
const SteakKind = "Steak"

// CookLevel is a type that describes how thoroughly a Steak should be cooked
type CookLevel string

const (
	RARE        CookLevel = "rare"
	MEDIUM_RARE CookLevel = "medium_rare"
	MEDIUM      CookLevel = "medium"
	MEDIUM_WELL CookLevel = "medium_well"
	WELL        CookLevel = "well_done"
)

// Fat is a type that describes the amount of fat on a Steak
type Fat string

const (
	FAT_LEAN   Fat = "lean"
	FAT_MEDIUM Fat = "medium"
	FAT_JUICY  Fat = "juicy"
)

// Variant is a type that describes the variant/kind of meat of a Steak
type Variant string

const (
	// VARIANT_BEEF is a beef steak
	VARIANT_BEEF Variant = "beef"
	// VARIANT_PORK is a pork steak
	VARIANT_PORK Variant = "pork"
	// VARIANT_CHICKEN is a chicken steak
	VARIANT_CHICKEN Variant = "chicken"
	// VARIANT_SALMON is a salmon steak
	VARIANT_SALMON Variant = "salmon"
	// VARIANT_TOFU is a vegetarian tofu steak
	VARIANT_TOFU Variant = "tofu"
)

// SteakSpec describes the desired state of s Steak
type SteakSpec struct {
	// CookLevel defines how strong the Steak should be cooked
	CookLevel CookLevel `json:"cookLevel"`
	// Fat defines how much fat a Steak should have
	Fat Fat `json:"fat"`
	// Weight is the desired weight of a Steak
	Weight int `json:"weight"`
	// Variant is the kind of Steak
	Variant Variant `json:"variant"`
}

// SteakStatus defines the observed state of Steak
type SteakStatus struct {
	// CookStatus shows the cooking level of a Steak
	CookStatus string `json:"cookStatus"`
	// ServingWeight describes the actual weight of the cooked Steak
	ServingWeight int `json:"servingWeight"`
	// Served denotes wether the Steak is still being cooked or served
	Served bool `json:"served"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=stk
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:JSONPath=.spec.weight,type=integer,name=Weight,description="The weight of this Steak."
//+kubebuilder:printcolumn:JSONPath=.spec.cookLevel,type=string,name=Cooked,description="The level of how thorough this Steak is cooked."
//+kubebuilder:printcolumn:JSONPath=.spec.fat,type=string,name=Fat,description="The juicyness of this Steak."
//+kubebuilder:printcolumn:JSONPath=.spec.variant,type=string,name=Variant,description="The kind of this Steak."

// Steak is the Schema for the steaks API
type Steak struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SteakSpec   `json:"spec,omitempty"`
	Status SteakStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SteakList contains a list of Steak
type SteakList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Steak `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Steak{}, &SteakList{})
}
