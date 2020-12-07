package clusterRegister

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	Group   = "clusterregistry.k8s.io"
	Version = "v1alpha1"
)

var (
	// ClusterRegistryGVK is group version used to register these objects
	ClusterRegistryGVK = schema.GroupVersionKind{Group: Group, Version: Version, Kind: "Cluster"}
	ClusterRegistryGVR = schema.GroupVersionResource{Group: Group, Version: Version, Resource: "clusters"}

	// GroupVersion is group version used to register these objects
	//ClusterRegistryGroupVersion = schema.GroupVersion{Group: "clusterregistry.k8s.io", Version: "v1alpha1"}
)

// Cluster contains information about a cluster in a cluster registry.
type Cluster struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the specification of the cluster. This may or may not be
	// reconciled by an active controller.
	Spec ClusterSpec `json:"spec,omitempty"`

	// Status is the status of the cluster.
	// +optional
	Status ClusterStatus `json:"status,omitempty"`
}

// ClusterSpec defines the desired state of Cluster
type ClusterSpec struct {
	// KubernetesAPIEndpoints represents the endpoints of the API server for this
	// cluster.
	KubernetesAPIEndpoints KubernetesAPIEndpoints `json:"kubernetesApiEndpoints,omitempty" protobuf:"bytes,1,opt,name=kubernetesApiEndpoints"`

	// AuthInfo contains public information that can be used to authenticate
	// to and authorize with this cluster. It is not meant to store private
	// information (e.g., tokens or client certificates) and cluster registry
	// implementations are not expected to provide hardened storage for
	// secrets.
	AuthInfo AuthInfo `json:"authInfo,omitempty" protobuf:"bytes,2,opt,name=authInfo"`
}

// ClusterStatus defines the observed state of Cluster
type ClusterStatus struct {
	// Conditions contains the different condition statuses for this cluster.
	Conditions []ClusterCondition `json:"conditions,omitempty" protobuf:"bytes,1,rep,name=conditions"`

	// Version contains the kubernetes version for this cluster.
	Version string `json:"version,omitempty" protobuf:"bytes,2,rep,name=version"`

	// Type contains the type of this cluster.
	Type ClusterType `json:"type,omitempty" protobuf:"bytes,3,rep,name=type"`
}

// KubernetesAPIEndpoints represents the endpoints for one and only one
// Kubernetes API server.
type KubernetesAPIEndpoints struct {
	// ServerEndpoints specifies the address(es) of the Kubernetes API server’s
	// network identity or identities.
	ServerEndpoints []ServerAddressByClientCIDR `json:"serverEndpoints,omitempty" protobuf:"bytes,1,rep,name=serverEndpoints"`

	// CABundle contains the certificate authority information.
	// +optional
	CABundle []byte `json:"caBundle,omitempty" protobuf:"bytes,2,opt,name=caBundle"`
}

// ServerAddressByClientCIDR helps clients determine the server address that
// they should use, depending on the ClientCIDR that they match.
type ServerAddressByClientCIDR struct {
	// The CIDR with which clients can match their IP to figure out if they should
	// use the corresponding server address.
	// +optional
	ClientCIDR string `json:"clientCIDR,omitempty" protobuf:"bytes,1,opt,name=clientCIDR"`
	// Address of this server, suitable for a client that matches the above CIDR.
	// This can be a hostname, hostname:port, IP or IP:port.
	ServerAddress string `json:"serverAddress,omitempty" protobuf:"bytes,2,opt,name=serverAddress"`
}

// AuthInfo holds information that describes how a client can get
// credentials to access the cluster. For example, OAuth2 client registration
// endpoints and supported flows, or Kerberos server locations.
type AuthInfo struct {
	// User references an object that contains implementation-specific details
	// about how a user should authenticate against this cluster.
	// +optional
	User *ObjectReference `json:"user,omitempty" protobuf:"bytes,1,opt,name=user"`

	// Controller references an object that contains implementation-specific
	// details about how a controller should authenticate. A simple use case for
	// this would be to reference a secret in another namespace that stores a
	// bearer token that can be used to authenticate against this cluster's API
	// server.
	Controller *ObjectReference `json:"controller,omitempty" protobuf:"bytes,2,opt,name=controller"`
}

// ObjectReference contains enough information to let you inspect or modify the referred object.
type ObjectReference struct {
	// Kind contains the kind of the referent, e.g., Secret or ConfigMap
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +kubebuilder:validation:Enum=Secret
	Kind string `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`

	// Name contains the name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	Name string `json:"name,omitempty" protobuf:"bytes,2,opt,name=name"`

	// Namespace contains the namespace of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`
}

// ClusterConditionType marks the kind of cluster condition being reported.
// +kubebuilder:validation:Enum=NotAccessible;ComponentNotHealthy;NodeNotReady
type ClusterConditionType string

const (
	// NotAccessible means that if the cluster's apiserver can't be accessible
	// using the version endpoint without any auth.
	NotAccessible ClusterConditionType = "NotAccessible"
	// ComponentNotHealthy means that if any component of the cluster's control plane
	// components is unhealthy.
	ComponentNotHealthy ClusterConditionType = "ComponentNotHealthy"
	// NodeNotReady means that if any node of the cluster's nodes is not ready.
	NodeNotReady ClusterConditionType = "NodeNotReady"
)

// ClusterCondition contains condition information for a cluster.
type ClusterCondition struct {
	// Type is the type of the cluster condition.
	Type ClusterConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=ClusterConditionType"`

	// Status is the status of the condition. One of True, False, Unknown.
	Status v1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=ConditionStatus"`

	// LastHeartbeatTime is the last time this condition was updated.
	LastHeartbeatTime metav1.Time `json:"lastHeartbeatTime,omitempty" protobuf:"bytes,3,opt,name=lastHeartbeatTime"`

	// LastTransitionTime is the last time the condition changed from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`

	// Reason is a (brief) reason for the condition's last status change.
	Reason string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`

	// Message is a human-readable message indicating details about the last status change.
	Message string `json:"message,omitempty" protobuf:"bytes,6,opt,name=message"`
}

// ClusterType marks the kind of cluster type being reported.
// +kubebuilder:validation:Enum=Kubernetes;OpenShift
type ClusterType string
