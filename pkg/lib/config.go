package lib

import (
	"errors"
	pkgerrors "github.com/pkg/errors"
	"gomod.alauda.cn/asm/cluster-operator/api/clusterRegister"
	"k8s.io/client-go/rest"
)

func RestConfigOfCluster(c *clusterRegister.Cluster, token string) (*rest.Config, error) {
	if len(c.Spec.KubernetesAPIEndpoints.ServerEndpoints) == 0 || c.Spec.KubernetesAPIEndpoints.ServerEndpoints[0].ServerAddress == "" {
		return nil, pkgerrors.WithStack(errors.New("could not found an avaliable server endpoint from cluster config."))
	}

	config := &rest.Config{
		Host: c.Spec.KubernetesAPIEndpoints.ServerEndpoints[0].ServerAddress,
	}
	if len(c.Spec.KubernetesAPIEndpoints.CABundle) == 0 {
		config.TLSClientConfig.Insecure = true
	} else {
		config.TLSClientConfig.CAData = c.Spec.KubernetesAPIEndpoints.CABundle
	}

	config.BearerToken = token
	return config, nil
}
