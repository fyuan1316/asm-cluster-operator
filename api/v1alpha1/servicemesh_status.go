package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

func (in ServiceMeshStatus) IsUnknownPhase() bool {
	if in.Phase == Phase.Unknown {
		return true
	}
	return false
}
func (in ServiceMeshStatus) IsPendingPhase() bool {
	if in.Phase == Phase.Pending {
		return true
	}
	return false
}

func (in ServiceMeshStatus) IsCancel() bool {
	if in.getCondition(ConditionCancel) == ConditionTrue {
		return true
	}
	return false
}
func (in ServiceMesh) IsDeletingPhase() bool {
	return !in.DeletionTimestamp.IsZero()
}
func (in ServiceMeshStatus) isProvisioningPhase() bool {
	if in.Phase == Phase.Provisioning {
		return true
	}
	return false
}
func (in ServiceMesh) IsDeployingPhase() bool {
	if in.Status.isProvisioningPhase() {
		if in.Status.Version == "" {
			return true
		}
	}
	return false
}
func (in ServiceMesh) IsUpgradingPhase() bool {
	if in.Status.isProvisioningPhase() {
		if in.Status.Version != "" && in.Status.Version != in.Spec.Version {
			return true
		}
	}
	return false
}

///

func (in ServiceMeshStatus) getCondition(t ConditionType) corev1.ConditionStatus {
	for _, cond := range in.Condition {
		if cond.Type == t {
			return cond.Status
		}
	}
	return ConditionUnknown
}
func (in *ServiceMeshStatus) ChangeToPhaseProvisioning() {
	in.Phase = Phase.Provisioning
}
func (in *ServiceMeshStatus) ChangeToPhaseCancel() {
	in.Phase = Phase.Cancel
	in.ConvertConditionCancelToFalse()
}

///

func (in *ServiceMeshStatus) ConvertConditionCancelToFalse() {
	in.changeProvisionCondition("", "", ConditionCancel, ConditionFalse)
}
func (in *ServiceMeshStatus) changeProvisionCondition(reason, message string, conditionType ConditionType, conditionStatus corev1.ConditionStatus) {
	if in.Condition != nil {
		for idx, condition := range in.Condition {
			if condition.Type == conditionType {
				cond := &in.Condition[idx]
				cond.Status = conditionStatus
				cond.Reason = reason
				cond.Message = message
				break
			}
		}
	}
}
