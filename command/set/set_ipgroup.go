package set

import (
	"github.com/ZTE/knitctl/api"
	"github.com/ZTE/knitctl/command/common"
	"github.com/spf13/cobra"
	"io"
)

var (
	setIpgDes = ("Update ipgroup with specified options")

	setIpgExample = (`  # Update a ipgroup with specified host route in tenant.
  knitctl set ipgroup my-ipg --ips=10.10.10.5 --network=my-nw --tenant=my-user"`)
)

// NewCmdSet creates a command object for the generic "set" command, which
// retrieves one or more Command from a server.
func NewCmdSetIpg(out io.Writer, errOut io.Writer) *cobra.Command {
	options := &common.IpgroupOption{
		CmdOut: out,
		CmdErr: errOut,
	}

	cmd := &cobra.Command{
		Use:     "ipgroup [NAME] [flags]",
		Short:   setIpgDes,
		Long:    setIpgDes,
		Example: setIpgExample,
		Run: func(cmd *cobra.Command, args []string) {
			api.ExcuteError(options.RunSetIpg(cmd, args), errOut)
		},
		SuggestFor: []string{"ip", "group"},
		Aliases:    []string{"ipg"},
	}

	cmd.Flags().StringVarP(&options.Tenant, "tenant", "t", options.Tenant, "If present, update name ,ips, and size of ipgroup in tenant name. ")
	cmd.Flags().StringVarP(&options.Ips, "ips", "", options.Ips, "If present, update ips of ipgroup in tenant name. Note that if present size and ips, ips will work.")
	cmd.Flags().StringVarP(&options.Network_id, "network", "", options.Network_id, "If present, update ipgroup int network name or id.")

	cmd.MarkFlagRequired("tenant")
	cmd.MarkFlagRequired("ips")
	cmd.MarkFlagRequired("network")

	return cmd
}
