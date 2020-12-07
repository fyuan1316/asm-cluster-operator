package common

import "os"

const (
	istioNamespace         = "istio-system"
	defaultAlaudaNamespace = "cpaas-system"
	//defaultServiceMeshChartVersion = "v3.0"
)
const (
	alaudaNamespaceEnv = "ASM_ALAUDA_NAMESPACE"
	//serviceMeshChartVersionEnv = "SERVICEMESH_CHART_VERSION"
)

func GetLocalBaseDomain() string {
	if os.Getenv("LABEL_BASE_DOMAIN") != "" {
		return os.Getenv("LABEL_BASE_DOMAIN")
	}
	return "alauda.io"
}

func AlaudaNamespace() string {
	if ns := os.Getenv(alaudaNamespaceEnv); ns != "" {
		return ns
	}
	return defaultAlaudaNamespace
}

func IstioNamespace() string {
	return istioNamespace
}
