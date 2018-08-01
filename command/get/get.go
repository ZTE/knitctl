package get

import (
	"github.com/spf13/cobra"
	"io"
	"fmt"
	"github.com/ZTE/knitctl/api"
)

var (
	getShort = ("Display one or many resource.")
	getLong = (`Display one or many resource.`)

	getExample = (`  # List a tenant with name.
  knitctl get tenant my-user

  # List all tenants.
  knitctl get tenant

  # List a network with name in tenant.
  knitctl get network my-network --tenant=my-user

  # List all networks in tenant.
  knitctl get network --tenant=my-user

  # List a ipgroup with name in tenant.
  knitctl get ipgroup my-ipg --tenant=my-user

  # List all ipgroups in tenant.
  knitctl get ipgroup --tenant=my-user

  # List a pod with name in tenant.
  knitctl get pod my-pod --tenant=my-user

  # List all pods in tenant.
  knitctl get pod --tenant=my-user
`)
)

func NewCmdGet(out io.Writer, errOut io.Writer) *cobra.Command {

	cmd := &cobra.Command{
		Use: "get [(TYPE [NAME])] [flags]",
		Short:   getShort,
		Long:    getLong,
		Example: getExample,
		Run: func(cmd *cobra.Command, args []string) {
			api.ExcuteError(RunGet(cmd, args), errOut)
		},
		SuggestFor: []string{"list", "ps"},
		ValidArgs:  []string{"network", "tenant", "ipgroup", "pod"},
		ArgAliases: []string{"nw", "ipg", "pod"},
	}

	cmd.AddCommand(NewCmdGetTN(out, errOut))
	cmd.AddCommand(NewCmdGetNW(out, errOut))
	cmd.AddCommand(NewCmdGetIpg(out, errOut))
	cmd.AddCommand(NewCmdGetPod(out, errOut))

	return cmd
}

func RunGet(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		err := fmt.Errorf("invalid argument '%s' to list resource.\n " +
			`Use "knitctl get --help" for more information about a given command.`, args)
		return err
	}
	cmd.HelpFunc()(cmd, args)
	return nil
}
