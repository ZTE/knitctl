package set

import (
	"github.com/ZTE/knitctl/api"
	"github.com/ZTE/knitctl/command/common"
	"github.com/spf13/cobra"
	"io"
)

var (
	setTNDes = ("Update tenant quota.")

	setTNExample = (`  # Update a tenant with specified quota.
  knitctl set tenant my-user --quota=2`)
)

// NewCmdSet creates a command object for the generic "set" command, which
// retrieves one or more Command from a server.
func NewCmdSetTN(out io.Writer, errOut io.Writer) *cobra.Command {
	options := &common.TenantOption{
		CmdOut: out,
		CmdErr: errOut,
	}

	cmd := &cobra.Command{
		Use:     "tenant [NAME] [flags]",
		Short:   setTNDes,
		Long:    setTNDes,
		Example: setTNExample,
		Run: func(cmd *cobra.Command, args []string) {
			api.ExcuteError(options.RunSetTN(cmd, args), errOut)
		},
		SuggestFor: []string{"ten", "nant"},
	}

	cmd.Flags().StringVarP(&options.Quota, "quota", "q", options.Quota, "If present, update tenant quota(number of network). ")
	//quota is required when command is "knitctl set tenant --quota"
	cmd.MarkFlagRequired("quota")

	return cmd
}
