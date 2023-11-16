package manager

import (
	"context"
	"embed"
	fluxapi1alpha1 "github.com/kluster-manager/fluxcd-addon/api/api/v1alpha1"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/component-base/version"
	"k8s.io/klog/v2"
	"open-cluster-management.io/addon-framework/pkg/addonfactory"
	"open-cluster-management.io/addon-framework/pkg/addonmanager"
	cmdfactory "open-cluster-management.io/addon-framework/pkg/cmd/factory"
	"open-cluster-management.io/api/addon/v1alpha1"
	_ "open-cluster-management.io/api/addon/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:embed agent-manifests
//go:embed agent-manifests/flux2
//go:embed agent-manifests/flux2/templates/_helper.tpl
var FS embed.FS

const (
	// AddonName represents the name of the FluxCD addon.
	AddonName = "fluxcd-addon-manager"

	// AgentManifestsDir is the directory containing Flux2 agent-manifests.
	AgentManifestsDir = "agent-manifests/flux2"

	// AgentInstallNamespace is the namespace where the FluxCD addon will be installed.
	AgentInstallNamespace = "flux-system"
)

var (
	// scheme is a runtime.Scheme for managing API resources.
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(fluxapi1alpha1.AddToScheme(scheme))
}

// NewManagerCommand creates a command for starting the addon manager controller.
func NewManagerCommand() *cobra.Command {
	cmd := cmdfactory.
		NewControllerCommandConfig(AddonName, version.Get(), runManagerController).
		NewCommand()
	cmd.Use = "manager"
	cmd.Short = "Starts the addon manager controller"

	return cmd
}

// runManagerController initializes and runs the addon manager controller.
// It sets up the required Kubernetes client, agent, and manager to manage the addon.
func runManagerController(ctx context.Context, kubeConfig *rest.Config) error {
	kubeClient, err := client.New(kubeConfig, client.Options{Scheme: scheme})
	if err != nil {
		klog.Errorf("Creating kube client failed: `%v`", err)
		return err
	}

	// Create an instance of the addon manager.
	mgr, err := addonmanager.New(kubeConfig)
	if err != nil {
		return err
	}

	// Initialize the agent addon factory and configure it.
	agent, err := addonfactory.NewAgentAddonFactory(AddonName, FS, AgentManifestsDir).
		WithConfigGVRs(
			schema.GroupVersionResource{Group: FluxCDConfigGroup, Version: FluxCDConfigVersion, Resource: FluxCDConfigResource},
		).
		WithGetValuesFuncs(GetConfigValues(kubeClient)).
		WithAgentHealthProber(agentHealthProber()).
		WithAgentInstallNamespace(func(addon *v1alpha1.ManagedClusterAddOn) string { return AgentInstallNamespace }).
		BuildHelmAgentAddon()
	if err != nil {
		klog.Error("Failed to build agent: `%v`", err)
		return err
	}

	// Add the agent to the manager.
	if err = mgr.AddAgent(agent); err != nil {
		return err
	}

	go mgr.Start(ctx)
	<-ctx.Done()

	return nil
}
