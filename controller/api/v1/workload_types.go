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

// WorkloadSpec defines the desired state of Workload
type WorkloadSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Node is where workload running
	Node string `json:"node"`
	// Pod name this workload represent in same namespace
	Pod string `json:"pod"`
	// Pod UID
	PodUID string `json:"podUID"`
	// ContainerID is pause container id this workload bind to
	ContainerID string `json:"containerID"`
}

type Interface struct {
	// EndpointClaim name
	EndpointClaim string `json:"endpoint"`
	// Master if true, this is master interface
	Master bool `json:"master"`
}

// WorkloadStatus defines the observed state of Workload
type WorkloadStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Phase : Initial, Ready
	Phase string `json:"phase"`
	// Interfaces : interfaces bind to workloads, first interface is primary
	// name: {network}-{ip}
	Interfaces []Interface `json:"interfaces"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// Workload is the Schema for the workloads API
type Workload struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WorkloadSpec   `json:"spec,omitempty"`
	Status WorkloadStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// WorkloadList contains a list of Workload
type WorkloadList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Workload `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Workload{}, &WorkloadList{})
}
