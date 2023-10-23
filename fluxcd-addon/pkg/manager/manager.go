package manager

import (
	"context"
	"embed"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/component-base/version"
	"k8s.io/klog/v2"
	"open-cluster-management.io/addon-framework/pkg/addonfactory"
	"open-cluster-management.io/addon-framework/pkg/addonmanager"
	cmdfactory "open-cluster-management.io/addon-framework/pkg/cmd/factory"
	"open-cluster-management.io/addon-framework/pkg/utils"
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
	mgr, err := addonmanager.New(kubeconfig)
	if err != nil {
		return err
	}
	agent, err := addonfactory.NewAgentAddonFactory(AddonName, FS, "manifests/helloworld").
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
