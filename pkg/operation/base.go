package operation

import (
	"context"
	"github.com/go-logr/logr"
	"gomod.alauda.cn/asm/cluster-operator/api/mesh"
	operatorv1alpha1 "gomod.alauda.cn/asm/cluster-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewProvisionHandler(
	Ctx context.Context,
	meshCopy *operatorv1alpha1.ServiceMesh,
	log logr.Logger,
	scheme *runtime.Scheme,
	client client.Client,
	recorder record.EventRecorder,
	meshConfig *mesh.ServiceMesh,
) Handler {
	handler := &StagesHandler{
		MeshHandler: MeshHandler{
			MeshContext: MeshContext{
				Ctx:        Ctx,
				Log:        log,
				Recorder:   recorder,
				Scheme:     scheme,
				Client:     client,
				MeshRef:    meshCopy,
				MeshConfig: meshConfig,
			},
		},
		MeshStatusHandler: MeshStatusHandler{
			MeshRef: meshCopy,
		},
	}
	return handler
}

type StagesHandler struct {
	MeshHandler
	MeshStatusHandler
}

func (s StagesHandler) Do() Resulter {
	// edge check
	if s.MeshRef.Status.IsUnknownPhase() || s.MeshRef.Status.IsPendingPhase() {
		s.MeshRef.Status.ChangeToPhaseProvisioning()
		return &Result{Changed: true}
	}
	// 终止
	if s.MeshRef.Status.IsCancel() {
		s.MeshRef.Status.ChangeToPhaseCancel()
		return &Result{Changed: true}
	}

	// 是否删除
	if s.MeshRef.IsDeletingPhase() {
		return s.Remove()
	}

	// 部署
	if s.MeshRef.IsDeployingPhase() {
		return s.Deploy()
	}

	// 升级
	if s.MeshRef.IsUpgradingPhase() {
		return s.Upgrade()
	}

	return nil
}

func (s StagesHandler) Remove() Resulter {
	panic("implement me")
}

func (s StagesHandler) Upgrade() Resulter {
	panic("implement me")
}

func (s StagesHandler) HealthCheck() {
	panic("implement me")
}

func (s StagesHandler) ProcessRunning() {
	panic("implement me")
}

type MeshHandler struct {
	MeshContext
}
type MeshStatusHandler struct {
	MeshRef *operatorv1alpha1.ServiceMesh
}

type MeshContext struct {
	Ctx context.Context
	client.Client
	Log        logr.Logger
	Recorder   record.EventRecorder
	Scheme     *runtime.Scheme
	MeshRef    *operatorv1alpha1.ServiceMesh
	MeshConfig *mesh.ServiceMesh
}

type executeItem struct {
	payload   interface{}
	err       error
	preRun    func() error
	postRun   func() error
	preCheck  func() (bool, error)
	postCheck func() (bool, error)
}
