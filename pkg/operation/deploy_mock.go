package operation

import (
	"context"
	"fmt"
	"gomod.alauda.cn/asm/cluster-operator/api/mesh"
	"gomod.alauda.cn/asm/cluster-operator/pkg/lib"
	"time"
)

func CreateOrUpdateHelmRequest(tplReq *lib.HelmRequest) error {
	fmt.Printf("deploy create request, %s", tplReq.Name)
	return nil
}

func WaitStageDeployFinished(ctx context.Context, items []*executeItem) error {
	return loopItemsUntil(ctx, 5*time.Second, 5, items, func(item *executeItem) (bool, error) {
		request := item.payload.(*lib.HelmRequest)
		fmt.Printf("waiting for helmrequest %s synced ", request.Name)
		//synced, err := m.isHelmRequestSynced(request)
		return true, nil
	})
}
func UpdateHelmRequestFailedReason(item *executeItem, err1 error) error {
	fmt.Println("deploy UpdateHelmRequestFailedReason")
	return nil
}
func EnsureIstioNamespace() error {
	fmt.Println("deploy EnsureIstioNamespace")
	return nil
}

func EnsureJaegerGlobalIngressService(sm *mesh.ServiceMesh) error {
	fmt.Println("deploy EnsureJaegerGlobalIngressService")
	return nil
}
