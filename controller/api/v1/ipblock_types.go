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

// IPBlockSpec defines the desired state of IPBlock
type IPBlockSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Network is network name reference
	Network string `json:"network"`
	// SubnetSlice is subnets in blocks
	SubnetSlice SubnetSlice `json:"subnetSlice,omitempty"`
	// NetMask is netmask ip block occupy
	NetMask int `json:"netMask"`
	// NodesAffinity is nodes bind this ip block, if empty, affinity all nodes
	NodesAffinity []string `json:"nodesAffinity,omitempty"`
}

// IPBlockStatus defines the observed state of IPBlock
type IPBlockStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Allocations, ip allocations
	Allocations map[string]int `json:"allocations,omitempty"`
	// Unallocated, ip unallocated
	Unallocated []string `json:"unallocated"`

	// Phase
	// Pending: ip block is waiting for allocating
	// Running: ip block allocated
	Phase NetworkPhase `json:"phase"`
	// Message : message for phase
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:subresource:status
// IPBlock is the Schema for the ipblocks API
type IPBlock struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IPBlockSpec   `json:"spec,omitempty"`
	Status IPBlockStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IPBlockList contains a list of IPBlock
type IPBlockList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IPBlock `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IPBlock{}, &IPBlockList{})
}
