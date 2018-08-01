package get

import (
	"github.com/spf13/cobra"
	"io"
	"github.com/ZTE/knitctl/api"
	"github.com/ZTE/knitctl/command/common"
)

var (
	getIpgDes = ("Display one or many ipgroup in tenant.")

	getIpgExample = (`  # List a ipgroup with name in tenant.
  knitctl get ipgroup my-ipg --tenant=my-user

  # List all ipgroups in tenant.
  knitctl get ipgroup --tenant=my-user
`)
)

func NewCmdGetIpg(out io.Writer, errOut io.Writer) *cobra.Command {
	options := &common.IpgroupOption{
		CmdOut: out,
		CmdErr:errOut,
	}

	cmd := &cobra.Command{
		Use: "ipgroup [NAME] [flags]",
		Short:   getIpgDes,
		Long:    getIpgDes,
		Example: getIpgExample,
		Run: func(cmd *cobra.Command, args []string) {
			api.ExcuteError(options.RunGetIpg(cmd, args), errOut)
		},
		SuggestFor: []string{"ip", "group"},
		Aliases: []string{"ipg"},
	}

	cmd.Flags().StringVarP(&options.Tenant, "tenant", "t", options.Tenant, "If present, list ipgroup(s) in tenant name. ")
	cmd.Flags().StringVarP(&options.Output, "output", "o", options.Output, "If present, list all columns of ipgroup(s) in tenant name. Supported: wide.")

	cmd.MarkFlagRequired("tenant")

	return cmd
}
