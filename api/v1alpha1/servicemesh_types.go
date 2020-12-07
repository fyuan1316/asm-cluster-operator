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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ServiceMeshSpec defines the desired state of ServiceMesh
type ServiceMeshSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ServiceMesh. Edit ServiceMesh_types.go to remove/update
	Foo          string           `json:"foo,omitempty"`
	Version      string           `json:"version"`
	GlobalConfig GlobalConfigSpec `json:"globalConfig"`
	//Config       unstructured.Unstructured `json:"config"`
	Config string `json:"config"`
}

type GlobalConfigSpec struct {
	ClusterName   string `json:"clusterName"`
	ServerAddress string `json:"serverAddress"`
	CABundle      string `json:"caBundle"`
	Token         string `json:"token"`
}

// ServiceMeshStatus defines the observed state of ServiceMesh
type ServiceMeshStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Phase     string      `json:"phase"`
	Condition []Condition `json:"conditions,omitempty"`
	Version   string      `json:"version,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ServiceMesh is the Schema for the servicemeshes API
type ServiceMesh struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceMeshSpec   `json:"spec,omitempty"`
	Status ServiceMeshStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ServiceMeshList contains a list of ServiceMesh
type ServiceMeshList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceMesh `json:"items"`
}

var Phase = struct {
	Unknown      string
	Pending      string
	Provisioning string
	Running      string
	Failed       string
	Deleting     string

	Cancel string
}{
	"",
	"Pending",
	"Provisioning",
	"Running",
	"Failed",
	"Deleting",
	"Cancel",
}

func init() {
	SchemeBuilder.Register(&ServiceMesh{}, &ServiceMeshList{})
}
