package get

import (
	"github.com/spf13/cobra"
	"io"
	"github.com/ZTE/knitctl/api"
	"github.com/ZTE/knitctl/command/common"
)

var (
	getNWDes = ("Display one or many network in tenant.")

	getNWExample = (`  # List a network with name in tenant.
  knitctl get network my-network --tenant=my-user

  # List all networks in tenant.
  knitctl get network --tenant=my-user
`)
)

func NewCmdGetNW(out io.Writer, errOut io.Writer) *cobra.Command {
	options := &common.NetworkOption{
		CmdOut: out,
		CmdErr: errOut,
	}

	cmd := &cobra.Command{
		Use: "network [NAME] [flags]",
		Short:   getNWDes,
		Long:    getNWDes,
		Example: getNWExample,
		Run: func(cmd *cobra.Command, args []string) {
			api.ExcuteError(options.RunGetNW(cmd, args), errOut)
		},
		SuggestFor: []string{"net", "ntwork", "work"},
		Aliases: []string{"nw"},
	}

	cmd.Flags().StringVarP(&options.Tenant, "tenant", "t", options.Tenant, "If present, list network(s) in tenant name. ")
	cmd.Flags().StringVarP(&options.Output, "output", "o", options.Output, "If present, list all columns of network(s) in tenant name. Supported: wide.")

	cmd.MarkFlagRequired("tenant")

	return cmd
}
