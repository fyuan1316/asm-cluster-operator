package lib

import (
	"fmt"
	"gomod.alauda.cn/asm/cluster-operator/pkg/common"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var asmHelmRequestFailedReason = "asm." + common.GetLocalBaseDomain() + "/failedreason"

type helmRequestPhase string

const (
	helmRequestFailed   helmRequestPhase = "Failed"
	helmRequestSynced   helmRequestPhase = "Synced"
	helmRequestPending  helmRequestPhase = "Pending"
	helmRequestUnknown  helmRequestPhase = "Unknown"
	helmRequestDeleting helmRequestPhase = "Deleting"
)

type ChartName string

func HelmRequestName(cluster string, name ChartName) string {
	return fmt.Sprintf("asm-%s-%s", cluster, name)
}

const (
	asmInitChart           ChartName = "asm-init"
	istioInitChart         ChartName = "istio-init"
	istioChart             ChartName = "istio-install"
	jaegerOperatorChart    ChartName = "jaeger-operator"
	alaudaServiceMeshChart ChartName = "cluster-asm"
)

type HelmRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Status            HelmRequestStatus `json:"status"`
	Spec              HelmRequestSpec   `json:"spec"`
}

type HelmRequestSpec struct {
	ClusterName          string                 `json:"clusterName,omitempty"`
	InstallToAllClusters bool                   `json:"installToAllClusters,omitempty"`
	Dependencies         []string               `json:"dependencies,omitempty"`
	ReleaseName          string                 `json:"releaseName,omitempty"`
	Chart                string                 `json:"chart,omitempty"`
	Version              string                 `json:"version,omitempty"`
	Namespace            string                 `json:"namespace,omitempty"`
	ValuesFrom           []ValuesFromSource     `json:"valuesFrom,omitempty"`
	Values               map[string]interface{} `json:"values,omitempty"`
}
type ValuesFromSource struct {
	ConfigMapKeyRef *v1.ConfigMapKeySelector `json:"configMapKeyRef,omitempty"`
	SecretKeyRef    *v1.SecretKeySelector    `json:"secretKeyRef,omitempty"`
}

type HelmRequestStatus struct {
	Phase          helmRequestPhase `json:"phase,omitempty"`
	LastSpecHash   string           `json:"lastSpecHash,omitempty"`
	SyncedClusters []string         `json:"syncedClusters,omitempty"`
	Notes          string           `json:"notes,omitempty"`
	Version        string           `json:"version,omitempty"`
}
