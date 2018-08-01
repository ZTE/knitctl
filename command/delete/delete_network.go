package delete

import (
	"github.com/spf13/cobra"
	"io"
	"github.com/ZTE/knitctl/api"
	"github.com/ZTE/knitctl/command/common"
)

var (
	deleteNWShort = ("Delete a network in tenant with the specified name.")
	deleteNWLong = (`Delete a network in tenant with the specified name.`)

	deleteNWExample = (`  # Delete a network named my-network in tenant named my-user.
  knitctl delete network my-network --tenant=my-user
`)
)

// NewCmdCreateNW is a macro command to create a new network
func NewCmdDeleteNW(cmdOut io.Writer, errOut io.Writer) *cobra.Command {
	options := &common.NetworkOption{
		CmdOut: cmdOut,
		CmdErr:errOut,
	}

	cmd := &cobra.Command{
		Use: "network [NAME] [flags]",
		DisableFlagsInUseLine: true,
		Short:                 deleteNWShort,
		Long:                  deleteNWLong,
		Example:               deleteNWExample,
		Run: func(cmd *cobra.Command, args []string) {
			api.ExcuteError(options.RunDelete(cmd, args), options.CmdErr)
		},
		SuggestFor: []string{"net", "ntwork", "work"},
		Aliases: []string{"nw"},
	}
	cmd.Flags().StringVarP(&options.Tenant, "tenant", "t", options.Tenant, "If present, set tenant where knitctl delete network.")

	cmd.MarkFlagRequired("tenant")

	return cmd
}
