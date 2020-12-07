package v1alpha1

import (
	"gomod.alauda.cn/asm/cluster-operator/api/mesh"
	"sigs.k8s.io/yaml"
)

func (in ServiceMesh) ToServiceMeshType() (*mesh.ServiceMesh, error) {
	meshSpec := &mesh.ServiceMesh{}
	b := []byte(in.Spec.Config)
	err := yaml.Unmarshal(b, meshSpec)
	if err != nil {
		return nil, err
	}
	return meshSpec, nil
}
