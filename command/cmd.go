package command

import (
	"io"
	"github.com/ZTE/knitctl/api"
	"github.com/spf13/cobra"
	"github.com/ZTE/knitctl/command/create"
	"github.com/ZTE/knitctl/command/get"
	"github.com/ZTE/knitctl/command/set"
	"github.com/ZTE/knitctl/command/delete"
	"github.com/ZTE/knitctl/command/config"
	"github.com/ZTE/knitctl/command/version"
)

var (
	knitterctlDes = (api.APPName + " controls the knitter manager.")
)

func NewKnitctlCommand(out, err io.Writer) *cobra.Command {

	cmds := &cobra.Command{
		Use:   api.APPName,
		Short: knitterctlDes,
		Long:  knitterctlDes,
		Run:   runHelp,
	}

	cmds.AddCommand(create.NewCommandCreate(out, err))
	cmds.AddCommand(get.NewCmdGet(out, err))
	cmds.AddCommand(set.NewCmdSet(out, err))
	cmds.AddCommand(delete.NewCmdDelete(out, err))
	cmds.AddCommand(config.NewCommandConfig(out, err))
	cmds.AddCommand(version.NewCommandVersion(out, err))

	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
