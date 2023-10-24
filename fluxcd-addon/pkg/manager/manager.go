package manager

import (
	"context"
	"embed"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/component-base/version"
	"k8s.io/klog/v2"
	"open-cluster-management.io/addon-framework/pkg/addonfactory"
	"open-cluster-management.io/addon-framework/pkg/addonmanager"
	cmdfactory "open-cluster-management.io/addon-framework/pkg/cmd/factory"
	"open-cluster-management.io/addon-framework/pkg/utils"
	_ "open-cluster-management.io/api/addon/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:embed manifests
//go:embed manifests/helloworld
//go:embed manifests/helloworld/templates/_helpers.tpl
var FS embed.FS

const (
	AddonName = "fluxcd-addon"
)

func NewManagerCommand() *cobra.Command {
	cmd := cmdfactory.
		NewControllerCommandConfig(AddonName, version.Get(), runManagerController).
		NewCommand()
	cmd.Use = "manager"
	cmd.Short = "Starts the addon manager controller"

	return cmd
}

func runManagerController(ctx context.Context, kubeconfig *rest.Config) error {
	kubeclient, err := client.New(kubeconfig, client.Options{})
	if err != nil {
		klog.Errorf("Creating kube client failed: `%v`", err)
		return err
	}

	mgr, err := addonmanager.New(kubeconfig)
	if err != nil {
		return err
	}
	agent, err := addonfactory.NewAgentAddonFactory(AddonName, FS, "manifests/helloworld").
		WithConfigGVRs(
			schema.GroupVersionResource{Version: "v1", Resource: "configmaps"},
		).
		WithGetValuesFuncs(GetInstallationValuesFromConfigmap(kubeclient)).
		WithAgentHealthProber(utils.NewDeploymentProber(types.NamespacedName{Name: "helm-controller", Namespace: "flux-system"})).
		BuildHelmAgentAddon()
	if err != nil {
		klog.Error("Failed to build agent: %v", err)
		return err
	}

	if err = mgr.AddAgent(agent); err != nil {
		return err
	}

	go mgr.Start(ctx)
	<-ctx.Done()

	return nil
}
