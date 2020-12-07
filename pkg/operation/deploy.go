package operation

import (
	"context"
	"fmt"
	"gomod.alauda.cn/asm/cluster-operator/api/mesh"
	"gomod.alauda.cn/asm/cluster-operator/pkg/common"
	"gomod.alauda.cn/asm/cluster-operator/pkg/lib"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"strconv"
	"strings"
	"time"
)

func (s StagesHandler) Deploy() Resulter {
	stages, err := s.deployStages()
	if err != nil {
		return Result{
			State:   "",
			Reason:  "",
			Err:     err,
			Changed: false,
		}
	}
	ctx := s.MeshContext.Ctx
	err = s.processStages(ctx, stages, func(items []*executeItem) error {
		for _, item := range items {
			//if err := m.createOrUpdateHelmRequest(item.payload.(*lib.HelmRequest)); err != nil {
			if err := CreateOrUpdateHelmRequest(item.payload.(*lib.HelmRequest)); err != nil {
				item.err = err
				return err
			}
		}
		return WaitStageDeployFinished(ctx, items)
	})
	if err != nil {
		return Result{
			State:   "",
			Reason:  "",
			Err:     err,
			Changed: false,
		}
	} else {
		return nil
	}
}

func (s StagesHandler) processStages(ctx context.Context, stages [][]*executeItem, processStage func([]*executeItem) error) error {
	resetFailedReason := func() error {
		for _, items := range stages {
			for _, item := range items {
				if err := UpdateHelmRequestFailedReason(item, item.err); err != nil {
					return err
				}
			}
		}
		return nil
	}
	defer resetFailedReason()
	if err := resetFailedReason(); err != nil {
		return err
	}
	for _, items := range stages {
		for _, item := range items {
			if item.preCheck != nil {
				s.Log.Info("run precheck")
				if err := loopUntil(ctx, 5*time.Second, 10, item.preCheck); err != nil {
					item.err = err
					return err
				}
			}
		}
		for _, item := range items {
			if item.preRun != nil {
				s.Log.Info("run prerun")
				if err := item.preRun(); err != nil {
					item.err = err
					return err
				}
			}
		}
		if err := processStage(items); err != nil {
			return err
		}
		for _, item := range items {
			if item.postRun != nil {
				s.Log.Info("run postrun")
				if err := item.postRun(); err != nil {
					item.err = err
					return err
				}
			}
		}
		for _, item := range items {
			if item.postCheck != nil {
				s.Log.Info("run postcheck")
				if err := loopUntil(ctx, 5*time.Second, 10, item.postCheck); err != nil {
					item.err = err
					return err
				}
			}
		}
	}
	return nil
}
func (s StagesHandler) deployStages() ([][]*executeItem, error) {
	sm := s.MeshConfig
	logger := s.Log.WithName("deploy")
	asmInitReq := asmInitReqFactory(sm)
	//istioReq := istioReqFactory(sm)
	jaegerOperatorReq := jaegerOperatorReqFactory(sm)
	//alaudaServiceMeshReq := alaudaServiceMeshReqFactory(sm)
	return [][]*executeItem{
		{
			{payload: asmInitReq,
				preRun: func() error {
					if err := EnsureIstioNamespace(); err != nil {
						return err
					}
					return nil
				}},
			{payload: jaegerOperatorReq,
				postRun: func() error {
					logger.Info("ensure jaeger global ingress endpoint resoucres")
					if err := EnsureJaegerGlobalIngressService(sm); err != nil {
						return err
					}
					return nil
				},
			},
		},
		/*fy
		{
			{payload: istioReq,
				preRun: func() error {
					if sm.Spec.Egress.DeployMode == mesh.EgressFixedDeploy {
						logger.Info("egress gateway deployMode: fixed, label nodes")
						return m.labelNodeForEgressGW(sm)
					}
					logger.Info("egress gateway deployMode: float")
					return nil
				},
				postRun: func() error {
					err := m.reloadGrafana(sm)
					if err != nil {
						return err
					}
					logger.Info("ensure grafana global ingress endpoint resoucres")
					return m.ensureGrafanaGlobalIngressService(sm)
				}},
			{payload: alaudaServiceMeshReq},
		},*/
	}, nil
}

const (
	asmInitChart           lib.ChartName = "asm-init"
	istioInitChart         lib.ChartName = "istio-init"
	istioChart             lib.ChartName = "istio-install"
	jaegerOperatorChart    lib.ChartName = "jaeger-operator"
	alaudaServiceMeshChart lib.ChartName = "cluster-asm"
)

func asmInitReqFactory(sm *mesh.ServiceMesh) *lib.HelmRequest {

	alaudaNamespace := common.AlaudaNamespace()
	req := &lib.HelmRequest{
		TypeMeta: v1.TypeMeta{
			APIVersion: helmRequestGVK.GroupVersion().String(),
			Kind:       helmRequestGVK.Kind,
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      lib.HelmRequestName(sm.Spec.Cluster, asmInitChart),
			Namespace: alaudaNamespace,
		},
		Spec: lib.HelmRequestSpec{
			ClusterName:          sm.Spec.Cluster,
			InstallToAllClusters: false,
			Dependencies:         nil,
			ReleaseName:          lib.HelmRequestName(sm.Spec.Cluster, asmInitChart),
			Chart:                fullChartName(asmInitChart),
			Namespace:            alaudaNamespace,
			Version:              chartVersion(asmInitChart),
			Values: treeValue(map[string]interface{}{
				"install.mode": "asm",
			}),
		},
	}
	return req
}

/*
func istioReqFactory(sm *mesh.ServiceMesh) *lib.HelmRequest {

	u, err := url.Parse(sm.Spec.PrometheusURL)
	if err != nil {
		return nil
	}

	var enablePolicy = "disabled"
	if sm.Spec.IstioSidecarInjectorPolicy != nil && *sm.Spec.IstioSidecarInjectorPolicy == true {
		enablePolicy = "enabled"
	}

	req := &lib.HelmRequest{
		TypeMeta: v1.TypeMeta{
			APIVersion: helmRequestGVK.GroupVersion().String(),
			Kind:       helmRequestGVK.Kind,
		},

		ObjectMeta: v1.ObjectMeta{
			Name:      lib.HelmRequestName(sm.Spec.Cluster, istioChart),
			Namespace: common.AlaudaNamespace(),
		},
		Spec: lib.HelmRequestSpec{
			ClusterName:          sm.Spec.Cluster,
			InstallToAllClusters: false,
			Dependencies:         []string{lib.HelmRequestName(sm.Spec.Cluster, asmInitChart)},
			ReleaseName:          lib.HelmRequestName(sm.Spec.Cluster, istioChart),
			Chart:                fullChartName(istioChart),
			Version:              chartVersion(istioChart),
			Namespace:            common.IstioNamespace(),
			Values: treeValue(map[string]interface{}{
				"global.hub":                      fmt.Sprintf("%s/%s", sm.Spec.RegistryAddress, alaudak8sRepo),
				"global.labelBaseDomain":          common.GetLocalBaseDomain(),
				"istio.includeIPRanges":           strings.Join(sm.Spec.IPRanges.Ranges, ","),
				"pilot.traceSampling":             fmt.Sprintf("%.2f", sm.Spec.TraceSampling),
				"prometheus.url":                  u.Host,
				"proxy.autoInject":                enablePolicy,
				"proxy.cpuLimit":                  fmt.Sprintf("%s", sm.Spec.IstioSidecar.CpuValue),
				"proxy.memoryLimit":               fmt.Sprintf("%s", sm.Spec.IstioSidecar.MemoryValue),
				"proxy.cpuRequest":                fmt.Sprintf("%s", sm.Spec.IstioSidecar.CpuValue),
				"proxy.memoryRequest":             fmt.Sprintf("%s", sm.Spec.IstioSidecar.MemoryValue),
				"egressGateways.fixedDeploy":      sm.Spec.Egress.IsFixedDeploy(),
				"egressGateways.autoscaleEnabled": false,
				"egressGateways.replicaCount":     1,
				"egressGateways.deployrevision":   sm.Spec.Egress.GenUUIDByHosts(),
			}),
		}}

	setTreeValue(req.Spec.Values, "grafana.rootUrl", sm.Spec.GlobalIngressHost+globalGrafanaNamePrefix+sm.Spec.Cluster)

	if sm.Spec.HighAvailability {
		setTreeValue(req.Spec.Values, "pilot.replicaCount", 2)
		setTreeValue(req.Spec.Values, "pilot.hpaSpec.minReplicas", 2)
		setTreeValue(req.Spec.Values, "policy.replicaCount", 2)
		setTreeValue(req.Spec.Values, "policy.autoscaleMin", 2)
		setTreeValue(req.Spec.Values, "telemetry.replicaCount", 2)
		setTreeValue(req.Spec.Values, "telemetry.autoscaleMin", 2)
		setTreeValue(req.Spec.Values, "ingressGateways.replicaCount", 2)
		setTreeValue(req.Spec.Values, "ingressGateways.autoscaleMin", 2)
		setTreeValue(req.Spec.Values, "egressGateways.autoscaleEnabled", true)
		setTreeValue(req.Spec.Values, "egressGateways.autoscaleMin", 2)
	}
	return req
}
*/

func jaegerOperatorReqFactory(sm *mesh.ServiceMesh) *lib.HelmRequest {
	req := &lib.HelmRequest{
		TypeMeta: v1.TypeMeta{
			APIVersion: helmRequestGVK.GroupVersion().String(),
			Kind:       helmRequestGVK.Kind,
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      lib.HelmRequestName(sm.Spec.Cluster, jaegerOperatorChart),
			Namespace: common.AlaudaNamespace(),
		},
		Spec: lib.HelmRequestSpec{
			ClusterName:          sm.Spec.Cluster,
			InstallToAllClusters: false,
			ReleaseName:          lib.HelmRequestName(sm.Spec.Cluster, jaegerOperatorChart),
			Dependencies:         nil,
			Version:              chartVersion(jaegerOperatorChart),
			Chart:                fullChartName(jaegerOperatorChart),
			Namespace:            common.IstioNamespace(),
			Values: treeValue(map[string]interface{}{
				"fullnameOverride":       "jaeger-operator",
				"image.repository":       fmt.Sprintf("%s/%s/jaeger-operator", sm.Spec.RegistryAddress, alaudak8sRepo),
				"global.labelBaseDomain": common.GetLocalBaseDomain(),
			}),
		}}
	return req
}

/*
func alaudaServiceMeshReqFactory(sm *mesh.ServiceMesh) *lib.HelmRequest {

	req := &lib.HelmRequest{
		TypeMeta: v1.TypeMeta{
			APIVersion: helmRequestGVK.GroupVersion().String(),
			Kind:       helmRequestGVK.Kind,
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      lib.HelmRequestName(sm.Spec.Cluster, alaudaServiceMeshChart),
			Namespace: common.AlaudaNamespace(),
		},
		Spec: lib.HelmRequestSpec{
			ClusterName:          sm.Spec.Cluster,
			InstallToAllClusters: false,
			Dependencies: []string{
				lib.HelmRequestName(sm.Spec.Cluster, istioChart),
				lib.HelmRequestName(sm.Spec.Cluster, jaegerOperatorChart),
				lib.HelmRequestName(sm.Spec.Cluster, asmInitChart)},
			Chart:       fullChartName(alaudaServiceMeshChart),
			ReleaseName: lib.HelmRequestName(sm.Spec.Cluster, alaudaServiceMeshChart),
			Version:     chartVersion(alaudaServiceMeshChart),
			Namespace:   common.AlaudaNamespace(),
			Values: treeValue(map[string]interface{}{
				"global.registry.address":          sm.Spec.RegistryAddress,
				"global.labelBaseDomain":           common.GetLocalBaseDomain(),
				"global.scheme":                    sm.Spec.IngressScheme,
				"global.useNodePort":               true,
				"prometheus.url":                   sm.Spec.PrometheusURL,
				"prometheus.serviceMonitorLabels":  sm.Spec.ServiceMonitorLabels,
				"grafana.url":                      sm.Spec.GlobalIngressHost + globalGrafanaNamePrefix + sm.Spec.Cluster,
				"jaeger.query.basepath":            "/" + globalJaegerNamePrefix + sm.Spec.Cluster,
				"jaeger.url":                       sm.Spec.GlobalIngressHost + globalJaegerNamePrefix + sm.Spec.Cluster,
				"jaeger.elasticsearch.serverurl":   sm.Spec.Elasticsearch.URL,
				"jaeger.elasticsearch.username":    wrapInt(sm.Spec.Elasticsearch.Username),
				"jaeger.elasticsearch.password":    wrapInt(sm.Spec.Elasticsearch.Password),
				"jaeger.elasticsearch.indexprefix": fmt.Sprintf("asm-cluster-%s", sm.Spec.Cluster),
				"flagger.metricsServer":            sm.Spec.PrometheusURL,
				"flagger.selectorLabels":           fmt.Sprintf("service.%s/name,app,name,app.kubernetes.io/name", common.GetLocalBaseDomain()),
			}),
		}}

	if sm.Spec.HighAvailability {
		setTreeValue(req.Spec.Values, "jaeger.query.replicas", 2)
		setTreeValue(req.Spec.Values, "jaeger.collector.replicas", 2)
	}

	return req
}

*/
func treeValue(values map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range values {
		setTreeValue(result, k, v)
	}
	return result
}

func setTreeValue(tree map[string]interface{}, key string, value interface{}) {
	tmp := tree
	keys := strings.Split(key, ".")
	count := len(keys)
	for i, k := range keys {
		if i < count-1 {
			if _, ok := tmp[k]; !ok {
				tmp[k] = map[string]interface{}{}
			}
			tmp = tmp[k].(map[string]interface{})
		} else {
			tmp[k] = value
		}
	}
}
func fullChartName(name lib.ChartName) string {
	repo := defaultChartRepo
	if value := os.Getenv(fmt.Sprintf(ChartRepoEnv, strings.ToUpper(string(name)))); value != "" {
		repo = value
	}
	return fmt.Sprintf("%s/%s", repo, name)
}

func chartVersion(name lib.ChartName) string {
	return os.Getenv(fmt.Sprintf(ChartVersionEnv, strings.ToUpper(string(name))))
}

const (
	alaudak8sRepo    = "asm"
	ChartRepoEnv     = "CHARTREPO_%s"
	ChartVersionEnv  = "CHARTVERSION_%s"
	defaultChartRepo = "stable"
)
const (
	istioCrdGroup           = "istio.io"
	flaggerCrdGroup         = "flagger.app"
	minIstioCrdGroups       = 1
	globalJaegerNamePrefix  = "jaeger-asm-"
	globalJaegerNodePort    = 30668
	globalJaegerSvcPort     = 16686
	globalGrafanaNamePrefix = "grafana-asm-"
	globalGrafanaNodePort   = 30667
	globalGrafanaSvcPort    = 3000
)

func wrapInt(i string) string {
	if _, err := strconv.Atoi(i); err == nil {
		return fmt.Sprintf(`"%s"`, i)
	}
	return i
}

const (
	alaudaDomain = "alauda.io"
)

var (
	helmRequestGVK = &schema.GroupVersionKind{
		Group:   fmt.Sprintf("app.%s", alaudaDomain),
		Version: "v1alpha1",
		Kind:    "HelmRequest",
	}
	servicemeshGVK = &schema.GroupVersionKind{
		Group:   fmt.Sprintf("asm.%s", alaudaDomain),
		Version: "v1alpha1",
		Kind:    "ServiceMesh",
	}
)
