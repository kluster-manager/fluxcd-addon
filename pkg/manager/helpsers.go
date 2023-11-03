package manager

import (
	"context"
	"fmt"
	fluxcnfv1alpha "github.com/kluster-management/fluxcd-addon/api/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"open-cluster-management.io/addon-framework/pkg/addonfactory"
	agentapi "open-cluster-management.io/addon-framework/pkg/agent"
	"open-cluster-management.io/api/addon/v1alpha1"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
	workapiv1 "open-cluster-management.io/api/work/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// FluxCDConfigVersion defines the API version used for FluxCDConfigs.
	FluxCDConfigVersion = "v1alpha1"

	// FluxCDConfigResource is the resource name for FluxCDConfig objects.
	FluxCDConfigResource = "fluxcdconfigs"

	// FluxCDConfigGroup is the group name for FluxCDConfig objects.
	FluxCDConfigGroup = "fluxcd.open-cluster-management.io"
)

// GetConfigValues returns a function that retrieves and transforms configuration values from
// FluxCDConfig objects. The function fetches configuration references from a ManagedClusterAddOn,
// reads corresponding FluxCDConfig objects, extracts their specifications, and converts them into
// addonfactory.Values. These values are then used to customize configuration of addon-agent.
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

// agentHealthProber returns an instance of the agent's health prober. It is used for
// probing and checking the health status of the agent.
func agentHealthProber() *agentapi.HealthProber {
	return &agentapi.HealthProber{
		Type: agentapi.HealthProberTypeWork,
		WorkProber: &agentapi.WorkHealthProber{
			ProbeFields: []agentapi.ProbeField{
				{
					ResourceIdentifier: workapiv1.ResourceIdentifier{
						Group:     "apps",
						Resource:  "deployments",
						Name:      "helm-controller",
						Namespace: AgentInstallNamespace,
					},
					ProbeRules: []workapiv1.FeedbackRule{
						{
							Type: workapiv1.WellKnownStatusType,
						},
					},
				},
				{
					ResourceIdentifier: workapiv1.ResourceIdentifier{
						Group:     "apps",
						Resource:  "deployments",
						Name:      "source-controller",
						Namespace: AgentInstallNamespace,
					},
					ProbeRules: []workapiv1.FeedbackRule{
						{
							Type: workapiv1.WellKnownStatusType,
						},
					},
				},
			},
			HealthCheck: func(identifier workapiv1.ResourceIdentifier, result workapiv1.StatusFeedbackResult) error {
				if len(result.Values) == 0 {
					return fmt.Errorf("no values are probed for deployment %s/%s", identifier.Namespace, identifier.Name)
				}
				for _, value := range result.Values {
					if value.Name != "ReadyReplicas" {
						continue
					}

					if *value.Value.Integer >= 1 {
						return nil
					}

					return fmt.Errorf("readyReplica is %d for deployement %s/%s", *value.Value.Integer, identifier.Namespace, identifier.Name)
				}
				return fmt.Errorf("readyReplica is not probed")
			},
		},
	}
}
