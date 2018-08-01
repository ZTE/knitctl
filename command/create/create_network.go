package create

import (
	"github.com/spf13/cobra"
	"io"
	"github.com/ZTE/knitctl/api"
	"github.com/ZTE/knitctl/command/common"
)

var (
	NWShort = ("Create a network in tenant with the specified name, gateway, cidr and public.")
	NWLong = (`Create a network in tenant with the specified name, gateway, cidr and public.`)

	NWExample = (`  # Create a new network named my-network in tenant named my-user
  knitctl create network my-network --public=false --gateway=10.10.10.1 --cidr=10.10.10.0/24 --tenant=my-user`)
)


// NewCmdCreateNW is a macro command to create a new network
func NewCmdCreateNW(cmdOut io.Writer, errOut io.Writer) *cobra.Command {
	options := &common.NetworkOption{
		CmdOut: cmdOut,
		CmdErr:errOut,
	}

	cmd := &cobra.Command{
		Use: "network [NAME] [flags]",
		DisableFlagsInUseLine: true,
		Short:                 NWShort,
		Long:                  NWLong,
		Example:               NWExample,
		Run: func(cmd *cobra.Command, args []string) {
			api.ExcuteError(options.RunCreate(cmd, args), options.CmdErr)
		},
		SuggestFor: []string{"net", "ntwork", "work"},
		Aliases: []string{"nw"},
	}
	cmd.Flags().BoolVarP(&options.Public, "public", "p", options.Public, `If present, set a public/private network to create, default false.`)
	cmd.Flags().StringVarP(&options.Gateway, "gateway", "g", options.Gateway, `If present, set gateway of creating network.`)
	cmd.Flags().StringVarP(&options.Cidr, "cidr", "c", options.Cidr, "If present, set cidr of creating network.")
	cmd.Flags().StringVarP(&options.Tenant, "tenant", "t", options.Tenant, "If present, set tenant where knitctl create network")

	cmd.MarkFlagRequired("gateway")
	cmd.MarkFlagRequired("cidr")
	cmd.MarkFlagRequired("tenant")

	//cmdutil.AddApplyAnnotationFlags(cmd)

	return cmd
}

