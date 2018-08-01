package create

import (
	"io"
	"github.com/spf13/cobra"
	"github.com/ZTE/knitctl/api"
	"github.com/ZTE/knitctl/command/common"
)

var(
	ipgShort = `Create a ipgroup in tenant with the specified name, size or ips, and network_id.`
	ipgLong = "Create a ipgroup in tenant with the specified name, size or ips, and network_id."
	ipgExample = (`  # Create a new ipgroup named my-ipg in tenant named my-user by size.
  knitctl create ipgroup my-ipg --size=1 --network=my-nw --tenant=my-user

  # Create a new ipgroup named my-ipg in tenant named my-user by ips.
  knitctl create ipgroup my-ipg --ips=10.10.10.5,10.10.10.6 --network=my-nw --tenant=my-user

  # Create a new ipgroup named my-ipg in tenant named my-user by ips and size, only ips can work.
  knitctl create ipgroup my-ipg --ips=10.10.10.5,10.10.10.6 --size=3 --network=my-nw --tenant=my-user`)
)

func NewCmdCreateIpgroup(cmdOut io.Writer, cmdErr io.Writer) *cobra.Command {
	options := &common.IpgroupOption{
		CmdOut: cmdOut,
		CmdErr: cmdErr,
	}

	cmd := &cobra.Command{
		Use: "ipgroup [NAME] [flags]",
		Short:                 ipgShort,
		Long:                  ipgLong,
		Example:               ipgExample,
		Run: func(cmd *cobra.Command, args []string) {
			api.ExcuteError(options.RunCreate(cmd, args), options.CmdOut)
		},
		SuggestFor: []string{"ip", "group"},
		Aliases: []string{"ipg"},
	}
	cmd.Flags().StringVarP(&options.Tenant, "tenant", "t", options.Tenant, "If present, set tenant where knitctl create ipgroup.")
	cmd.Flags().StringVarP(&options.Network_id, "network", "n", options.Network_id, `If present, set network of creating ipgroup.`)
	cmd.Flags().StringVarP(&options.Ips, "ips", "", options.Ips, "If present, set ips of creating ipgroup by ips, which should be separated by comma. Note that user can set ips by selecting either ips or size. If present both, ips will work.")
	cmd.Flags().IntVarP(&options.FSize, "size", "", options.FSize, "If present, set ips of creating ipgroup by size, knitter manager can automatically assign ips. Note that user can set ips by selecting either ips or size. If present both, ips will work.")

	cmd.MarkFlagRequired("network")
	cmd.MarkFlagRequired("tenant")

	return cmd
}
