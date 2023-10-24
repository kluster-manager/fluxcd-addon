package manager

import (
	"context"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"open-cluster-management.io/addon-framework/pkg/addonfactory"
	"open-cluster-management.io/api/addon/v1alpha1"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type HelmPreInstallationConfig map[string]interface{}

func GetInstallationValuesFromConfigmap(kubeclient client.Client) addonfactory.GetValuesFunc {
	return func(cluster *clusterv1.ManagedCluster, addon *v1alpha1.ManagedClusterAddOn) (addonfactory.Values, error) {
		overrideValues := addonfactory.Values{}
		for _, refConfig := range addon.Status.ConfigReferences {
			if refConfig.ConfigGroupResource.Group != "" ||
				refConfig.ConfigGroupResource.Resource != "configmaps" {
				continue
			}

			configMap := corev1.ConfigMap{}
			keyType := types.NamespacedName{Name: refConfig.Name, Namespace: refConfig.Namespace}
			if err := kubeclient.Get(context.TODO(), keyType, &configMap); err != nil {
				return nil, err
			}

			configData, _ := configMap.Data["values.yaml"]
			installationConfig := HelmPreInstallationConfig{}

			if err := yaml.Unmarshal([]byte(configData), &installationConfig); err != nil {
				return nil, err
			}

			values, err := addonfactory.JsonStructToValues(installationConfig)
			if err != nil {
				return nil, err
			}
			overrideValues = addonfactory.MergeValues(overrideValues, values)
		}
		return overrideValues, nil
	}
}
