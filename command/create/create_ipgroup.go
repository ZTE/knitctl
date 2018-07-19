package create

import (
	"github.com/ZTE/knitctl/api"
	"github.com/ZTE/knitctl/command/common"
	"github.com/spf13/cobra"
	"io"
)

var (
	ipgShort   = `Create a ipgroup in tenant with the specified name, size, and network_id.`
	ipgLong    = "Create a ipgroup in tenant with the specified name, size, and network_id."
	ipgExample = (`  # Create a new ipgroup named my-ipg in tenant named my-user
  knitctl create ipgroup my-ipg --size=1 --network=my-nw --tenant=my-user`)
)

func NewCmdCreateIpgroup(cmdOut io.Writer, cmdErr io.Writer) *cobra.Command {
	options := &common.IpgroupOption{
		CmdOut: cmdOut,
		CmdErr: cmdErr,
	}

	cmd := &cobra.Command{
		Use:     "ipgroup [NAME] [flags]",
		Short:   ipgShort,
		Long:    ipgLong,
		Example: ipgExample,
		Run: func(cmd *cobra.Command, args []string) {
			api.ExcuteError(options.RunCreate(cmd, args), options.CmdOut)
		},
		SuggestFor: []string{"ip", "group"},
		Aliases:    []string{"ipg"},
	}
	cmd.Flags().StringVarP(&options.Network_id, "network", "n", options.Network_id, `If present, set network of creating ipgroup.`)
	cmd.Flags().StringVarP(&options.Ips, "ips", "", options.Ips, "If present, set ips of creating ipgroup.")
	cmd.Flags().StringVarP(&options.Tenant, "tenant", "t", options.Tenant, "If present, set tenant where knitctl create ipgroup.")

	cmd.MarkFlagRequired("network")
	cmd.MarkFlagRequired("tenant")
	cmd.MarkFlagRequired("ips")

	return cmd
}
