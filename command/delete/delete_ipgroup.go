package delete

import (
	"github.com/spf13/cobra"
	"github.com/ZTE/knitctl/api"
	"io"
	"github.com/ZTE/knitctl/command/common"
)

var (
	deleteIPGDes = ("Delete a ipgroup in tenant with the specified name.")

	deleteIPGExample = (`  # Delete a ipgroup named my-ipg in tenant named my-user
  knitctl delete ipgroup my-ipg --tenant=my-user`)
)

func NewCmdDeleteIPG(cmdOut io.Writer, cmdErr io.Writer) *cobra.Command {
	options := &common.IpgroupOption{
		CmdOut: cmdOut,
		CmdErr: cmdErr,
	}

	cmd := &cobra.Command{
		Use: "ipgroup [NAME] [flags]",
		DisableFlagsInUseLine: true,
		Short:                 deleteIPGDes,
		Long:                  deleteIPGDes,
		Example:               deleteIPGExample,
		Run: func(cmd *cobra.Command, args []string) {
			api.ExcuteError(options.RunDelete(cmd, args), options.CmdOut)
		},
		SuggestFor: []string{"ip", "group"},
		Aliases: []string{"ipg"},
	}
	cmd.Flags().StringVarP(&options.Tenant, "tenant", "t", options.Tenant, "If present, set tenant where knitctl delete ipgroup.")

	cmd.MarkFlagRequired("tenant")

	return cmd
}