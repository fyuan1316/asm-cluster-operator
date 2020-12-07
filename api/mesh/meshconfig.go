package mesh

import (
	"gomod.alauda.cn/asm/cluster-operator/pkg/lib"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type ServiceMesh struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ServiceMeshSpec `json:"spec"`
}

type ServiceMeshSpec struct {
	Cluster                    string            `json:"cluster"`
	RegistryAddress            string            `json:"registryAddress"`
	IngressHost                string            `json:"ingressHost"`
	GlobalIngressHost          string            `json:"globalIngressHost"`
	ExternalHost               string            `json:"externalHost"`
	ClusterNodeIps             []string          `json:"clusterNodeIps"`
	TraceSampling              float64           `json:"traceSampling"`
	HighAvailability           bool              `json:"highAvailability"`
	PrometheusURL              string            `json:"prometheusURL"`
	ServiceMonitorLabels       map[string]string `json:"serviceMonitorLabels,omitempty"`
	IstioSidecarInjectorPolicy *bool             `json:"istioSidecarInjectorPolicy,omitempty"`
	IstioSidecar               struct {
		CpuValue    string `json:"cpuValue"`
		MemoryValue string `json:"memoryValue"`
	} `json:"istioSidecar"`

	IPRanges struct {
		Ranges []string `json:"ranges"`
	} `json:"ipranges"`
	Elasticsearch struct {
		IsDefault bool   `json:"isDefault"`
		Password  string `json:"password,omitempty"`
		Username  string `json:"username,omitempty"`
		URL       string `json:"url"`
	} `json:"elasticsearch"`
	IngressScheme string     `json:"ingressScheme"`
	Egress        EgressSpec `json:"egress"`
}

type EgressSpec struct {
	DeployMode EgressDeployMode `json:"deployMode"`
	HostNames  []string         `json:"hostNames,omitempty"`
}

func (e EgressSpec) IsFixedDeploy() bool {
	if e.DeployMode == EgressFixedDeploy {
		return true
	}
	return false
}
func (e EgressSpec) GenUUIDByHosts() string {
	if e.DeployMode == EgressFixedDeploy {
		seed := strings.Join(e.HostNames, "|")
		hashId := lib.ComputeHash(seed)
		return hashId
	}
	return "FloatDeploy"
}

type EgressDeployMode string

const (
	EgressFixedDeploy EgressDeployMode = "Fixed"
	EgressFloatDeploy EgressDeployMode = "Float"
)
