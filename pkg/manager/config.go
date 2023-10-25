package manager

import (
	"context"
	fluxcnfv1alpha "github.com/kluster-management/fluxcd-addon/api/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"open-cluster-management.io/addon-framework/pkg/addonfactory"
	"open-cluster-management.io/api/addon/v1alpha1"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	FluxCDConfigVersion  = "v1alpha1"
	FluxCDConfigResource = "fluxcdconfigs"
	FluxCDConfigGroup    = "fluxcd.open-cluster-management.io"
)

func GetConfigValues(kc client.Client) addonfactory.GetValuesFunc {
	return func(cluster *clusterv1.ManagedCluster, addon *v1alpha1.ManagedClusterAddOn) (addonfactory.Values, error) {
		overrideValues := addonfactory.Values{}
		for _, refConfig := range addon.Status.ConfigReferences {
			if refConfig.ConfigGroupResource.Group != FluxCDConfigGroup ||
				refConfig.ConfigGroupResource.Resource != FluxCDConfigResource {
				continue
			}

			fluxCDConfig := fluxcnfv1alpha.FluxCDConfig{}
			keyType := types.NamespacedName{Name: refConfig.Name, Namespace: refConfig.Namespace}

			if err := kc.Get(context.TODO(), keyType, &fluxCDConfig); err != nil {
				return nil, err
			}

			fluxCDConfigSpec := fluxCDConfig.Spec
			values, err := addonfactory.JsonStructToValues(fluxCDConfigSpec)
			if err != nil {
				return nil, err
			}
			overrideValues = addonfactory.MergeValues(overrideValues, values)
		}

		return overrideValues, nil
	}
}
