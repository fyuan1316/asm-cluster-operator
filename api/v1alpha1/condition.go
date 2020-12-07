package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConditionType string

const (
	ConditionInitialized ConditionType = "Initialized"
	ConditionReady       ConditionType = "Ready"
	ConditionReclaimed   ConditionType = "Reclaimed"
	ConditionCancel      ConditionType = "Cancel"
)
const (
	ConditionUnknown corev1.ConditionStatus = corev1.ConditionUnknown
	ConditionTrue    corev1.ConditionStatus = corev1.ConditionTrue
	ConditionFalse   corev1.ConditionStatus = corev1.ConditionFalse
)

type Condition struct {
	// Type of condition.
	// +required
	Type ConditionType `json:"type" description:"type of status condition"`

	// Status of the condition, one of True, False, Unknown.
	// +required
	Status corev1.ConditionStatus `json:"status" description:"status of the condition, one of True, False, Unknown"`

	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// We use VolatileTime in place of metav1.Time to exclude this from creating equality.Semantic
	// differences (all other things held constant).
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" description:"last time the condition transit from one status to another"`

	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" description:"one-word CamelCase reason for the condition's last transition"`

	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty" description:"human-readable message indicating details about last transition"`
}
