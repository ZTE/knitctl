package set

import (
	"github.com/ZTE/knitctl/api"
	"github.com/ZTE/knitctl/command/common"
	"github.com/spf13/cobra"
	"io"
)

var (
	setNwDes = ("Update network with specified options")

	setNwExample = (`  # Update a network in tenant with specified destination and nexthop.
  knitctl set network my-network --tenant=my-user --destination=10.10.10.1 --nexthop=10.10.9.1`)
)

func NewCmdSetNw(out io.Writer, errOut io.Writer) *cobra.Command {
	options := &common.NetworkOption{
		CmdOut: out,
		CmdErr: errOut,
	}

	cmd := &cobra.Command{
		Use:     "network [NAME] [flags]",
		Short:   setNwDes,
		Long:    setNwDes,
		Example: setNwExample,
		Run: func(cmd *cobra.Command, args []string) {
			api.ExcuteError(options.RunSetNw(cmd, args), errOut)
		},
		SuggestFor: []string{"net", "ntwork", "work"},
		Aliases:    []string{"nw"},
	}

	cmd.Flags().StringVarP(&options.Tenant, "tenant", "t", options.Tenant, "If present, update network in tenant name. ")
	cmd.Flags().StringVarP(&options.Destination, "destination", "d", options.Destination, "If present, update destination of network host route in tenant. ")
	cmd.Flags().StringVarP(&options.Nexthop, "nexthop", "", options.Nexthop, "If present, update nexthop of network host route in tenant. ")

	cmd.MarkFlagRequired("tenant")
	cmd.MarkFlagRequired("destination")
	cmd.MarkFlagRequired("nexthop")

	return cmd
}
