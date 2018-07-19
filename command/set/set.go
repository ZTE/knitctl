package set

import (
	"fmt"
	"github.com/ZTE/knitctl/api"
	"github.com/spf13/cobra"
	"io"
)

var (
	setDescription = ("Update resource with specified options")

	setExample = (`  # Update a tenant with specified quota.
  knitctl set tenant my-user --quota=2

  # Update a ipgroup in tenant with specified ips, size, and network_id.
  knitctl set ipgroup my-ipg --network=my-nw --tenant=my-user --size=2 --ips="10.10.10.2,10.10.10.3"`)
)

func NewCmdSet(out io.Writer, errOut io.Writer) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "set [(TYPE [NAME]) [flags]",
		Short:   setDescription,
		Long:    setDescription,
		Example: setExample,
		Run: func(cmd *cobra.Command, args []string) {
			api.ExcuteError(RunSet(cmd, args), errOut)
		},
		SuggestFor: []string{"update", "put"},
		ValidArgs:  []string{"network", "tenant", "ipgroup"},
		ArgAliases: []string{"nw", "ipg"},
	}
	cmd.AddCommand(NewCmdSetTN(out, errOut))
	//cmd.AddCommand(NewCmdSetNw(out, errOut))
	cmd.AddCommand(NewCmdSetIpg(out, errOut))

	return cmd
}

func RunSet(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("invalid argument '%s' to update resource.\n "+
			`Use "knitctl set --help" for more information about a given command.`, args)
	}
	cmd.HelpFunc()(cmd, args)
	return nil
}
