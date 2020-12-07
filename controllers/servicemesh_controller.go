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

package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	operatorv1alpha1 "gomod.alauda.cn/asm/cluster-operator/api/v1alpha1"
	errordef "gomod.alauda.cn/asm/cluster-operator/pkg/error"
	"gomod.alauda.cn/asm/cluster-operator/pkg/operation"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ServiceMeshReconciler reconciles a ServiceMesh object
type ServiceMeshReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=operator.asm.alauda.io,resources=servicemeshes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=operator.asm.alauda.io,resources=servicemeshes/status,verbs=get;update;patch

func (r *ServiceMeshReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	var err error
	rt := ctrl.Result{}
	ctx := context.Background()
	log := r.Log.WithValues("servicemesh", req.NamespacedName)
	log.Info(fmt.Sprintf("Starting reconcile loop for %v", req.NamespacedName))
	defer log.Info(fmt.Sprintf("Finish reconcile loop for %v", req.NamespacedName))

	mesh := &operatorv1alpha1.ServiceMesh{}
	if err := r.Get(context.Background(), req.NamespacedName, mesh); err != nil {
		if errors.IsNotFound(err) {
			return rt, nil
		}
		return rt, err
	}
	// todo fy validate mesh struct, throw error if any
	meshConfig, err := mesh.ToServiceMeshType()
	if err != nil {
		fmt.Println(err)
		return rt, err
	}

	meshCopy := mesh.DeepCopy()
	provision := operation.NewProvisionHandler(
		ctx,
		meshCopy,
		r.Log, r.Scheme, r.Client, r.Recorder,
		meshConfig)
	result := provision.Do()

	if result != nil && result.IsStatusChanged() {
		// get all subresource statuses

		// update service mesh status
		if err := r.Status().Update(context.Background(), meshCopy, &client.UpdateOptions{}); err != nil {
			r.Log.Error(err, "update mesh status error")
		}
	}

	return handleErr(rt, err)
}
func handleErr(rt ctrl.Result, err error) (ctrl.Result, error) {
	rt.Requeue = true
	if err != nil {
		if err == errordef.ErrNeedRetry {
			rt.RequeueAfter = errordef.ReconcileAfterDuration
		}
		err = nil
	} else {
		rt.RequeueAfter = errordef.HealthCheckReconcileAfterDuration
	}
	return rt, err
}

func (r *ServiceMeshReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1alpha1.ServiceMesh{}).
		Complete(r)
}
