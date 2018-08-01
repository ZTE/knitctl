package get

import (
	"github.com/spf13/cobra"
	"io"
	"github.com/ZTE/knitctl/api"
	"github.com/ZTE/knitctl/command/common"
)

var (
	getPoDes = ("Display one or many pod in tenant.")

	getPoExample = (`  # List a pod with name in tenant.
  knitctl get pod my-pod --tenant=my-user

  # List all pods in tenant.
  knitctl get pod --tenant=my-user
`)
)

func NewCmdGetPod(out io.Writer, errOut io.Writer) *cobra.Command {
	options := &common.PodOption{
		CmdErr:errOut,
		CmdOut:out,
	}

	cmd := &cobra.Command{
		Use: "pod [NAME] [flags]",
		Short:   getPoDes,
		Long:    getPoDes,
		Example: getPoExample,
		Run: func(cmd *cobra.Command, args []string) {
			api.ExcuteError(options.RunGetPod(cmd, args), errOut)
		},
		Aliases: []string{"po"},
	}

	cmd.Flags().StringVarP(&options.Tenant, "tenant", "t", options.Tenant, "If present, list pod(s) in tenant name. ")
	cmd.Flags().StringVarP(&options.Output, "output", "o", options.Output, "If present, list all columns of pod(s) in tenant name. Supported: wide.")
	cmd.MarkFlagRequired("tenant")

	return cmd
}


