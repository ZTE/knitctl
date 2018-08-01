package set

import (
	"github.com/spf13/cobra"
	"io"
	"github.com/ZTE/knitctl/api"
	"github.com/ZTE/knitctl/command/common"
)

var (
	setIpgDes = ("Update ipgroup with specified options")

	setIpgExample = (`  # Update a ipgroup ips in tenant by ips.
  knitctl set ipgroup my-ipg --ips=10.10.10.5 --network=my-nw --tenant=my-user

  # Update a ipgroup ips in tenant by size, automatically assign ips.
  knitctl set ipgroup my-ipg --size=2 --network=my-nw --tenant=my-user

  # Update a ipgroup ips in tenant by ips and size, only ips can work.
  knitctl set ipgroup my-ipg --ips=10.10.10.5 --size=2 --network=my-nw --tenant=my-user`)
)

// NewCmdSet creates a command object for the generic "set" command, which
// retrieves one or more Command from a server.
func NewCmdSetIpg(out io.Writer, errOut io.Writer) *cobra.Command {
	options := &common.IpgroupOption{
		CmdOut:out,
		CmdErr:errOut,
	}

	cmd := &cobra.Command{
		Use: "ipgroup [NAME] [flags]",
		Short:   setIpgDes,
		Long:    setIpgDes,
		Example: setIpgExample,
		Run: func(cmd *cobra.Command, args []string) {
			api.ExcuteError(options.RunSetIpg(cmd, args), errOut)
		},
		SuggestFor: []string{"ip", "group"},
		Aliases:  []string{"ipg"},
	}

	cmd.Flags().StringVarP(&options.Tenant, "tenant", "t", options.Tenant, "If present, update ips of ipgroup in tenant name. ")
	cmd.Flags().StringVarP(&options.Network_id, "network", "", options.Network_id, "If present, update ipgroup in network name or id.")
	cmd.Flags().StringVarP(&options.Ips, "ips", "", options.Ips, "If present, update ips of ipgroup in tenant by ips, ips shoud be separated by comma. Note that user can update ips by selecting either ips or size. If present both, ips will work.")
	cmd.Flags().IntVarP(&options.FSize, "size", "", options.FSize, "If present, update ips of ipgroup by size, knitter manager can automatically assign ips. Note that user can update ips by selecting either ips or size. If present both, ips will work.")

	cmd.MarkFlagRequired("tenant")
	cmd.MarkFlagRequired("network")

	return cmd
}
