package main

import (
	goflag "flag"
	"fmt"
	"github.com/kluster-management/fluxcd-addon/pkg/manager"
	"os"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/rand"
	utilflag "k8s.io/component-base/cli/flag"
	logs "k8s.io/component-base/logs/api/v1"
	"open-cluster-management.io/addon-framework/pkg/version"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	flag.CommandLine.SetNormalizeFunc(utilflag.WordSepNormalizeFunc)
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	logs.AddFlags(logs.NewLoggingConfiguration(), flag.CommandLine)

	command := newCommand()
	if err := command.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fluxcd",
		Short: "fluxcd addon manager",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
			}
			os.Exit(1)
		},
	}
	if v := version.Get().String(); len(v) == 0 {
		cmd.Version = "<unknown>"
	} else {
		cmd.Version = v
	}

	cmd.AddCommand(manager.NewManagerCommand())

	return cmd
}
