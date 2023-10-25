package manager

import (
	"context"
	"embed"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"

	fluxapi1alpha1 "github.com/kluster-management/fluxcd-addon/api/api/v1alpha1"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/component-base/version"
	"k8s.io/klog/v2"
	"open-cluster-management.io/addon-framework/pkg/addonfactory"
	"open-cluster-management.io/addon-framework/pkg/addonmanager"
	agentapi "open-cluster-management.io/addon-framework/pkg/agent"
	cmdfactory "open-cluster-management.io/addon-framework/pkg/cmd/factory"
	_ "open-cluster-management.io/api/addon/v1alpha1"
	workapiv1 "open-cluster-management.io/api/work/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:embed manifests
//go:embed manifests/flux2
//go:embed manifests/flux2/templates/_helper.tpl
var FS embed.FS

const (
	AddonName         = "fluxcd-addon"
	AgentManifestsDir = "manifests/flux2"

	AgentHealthProberName      = "helm-controller"
	AgentHealthProberNamespace = "flux-system"
)

func NewManagerCommand() *cobra.Command {
	cmd := cmdfactory.
		NewControllerCommandConfig(AddonName, version.Get(), runManagerController).
		NewCommand()
	cmd.Use = "manager"
	cmd.Short = "Starts the addon manager controller"

	return cmd
}

func runManagerController(ctx context.Context, kubeConfig *rest.Config) error {
	kubeClient, err := getKubeClient(kubeConfig)
	if err != nil {
		klog.Errorf("Creating kube client failed: `%v`", err)
		return err
	}

	mgr, err := addonmanager.New(kubeConfig)
	if err != nil {
		return err
	}
	agent, err := addonfactory.NewAgentAddonFactory(AddonName, FS, AgentManifestsDir).
		WithConfigGVRs(
			schema.GroupVersionResource{Group: FluxCDConfigGroup, Version: FluxCDConfigVersion, Resource: FluxCDConfigResource},
		).
		WithGetValuesFuncs(GetConfigValues(kubeClient)).
		WithAgentHealthProber(agentHealthProber()).
		BuildHelmAgentAddon()
	if err != nil {
		klog.Error("Failed to build agent: `%v`", err)
		return err
	}

	if err = mgr.AddAgent(agent); err != nil {
		return err
	}

	go mgr.Start(ctx)
	<-ctx.Done()

	return nil
}

func getKubeClient(kubeConfig *rest.Config) (client.Client, error) {
	scheme := runtime.NewScheme()
	err := fluxapi1alpha1.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}

	return client.New(kubeConfig, client.Options{Scheme: scheme})
}

func agentHealthProber() *agentapi.HealthProber {
	return &agentapi.HealthProber{
		Type: agentapi.HealthProberTypeWork,
		WorkProber: &agentapi.WorkHealthProber{
			ProbeFields: []agentapi.ProbeField{
				{
					ResourceIdentifier: workapiv1.ResourceIdentifier{
						Group:     "apps",
						Resource:  "deployments",
						Name:      AgentHealthProberName,
						Namespace: AgentHealthProberNamespace,
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
