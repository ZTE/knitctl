package version

import (
	"github.com/spf13/cobra"
	"io"
	"github.com/ZTE/knitctl/api"
	"fmt"
)

var (
	versionDes = (`list knitter server version.`)

	versionExample = (`  # list knitter version.
  knitctl version`)
)

func NewCommandVersion(out io.Writer, errOut io.Writer) *cobra.Command {

	cmd := &cobra.Command{
		Use: "version",
		DisableFlagsInUseLine: true,
		Short:   versionDes,
		Long:    versionDes,
		Example: versionExample,
		Run: func(cmd *cobra.Command, args []string) {
			api.ExcuteError(RunVersion(cmd, args, out), errOut)
		},
	}
	return cmd
}

func RunVersion(cmd *cobra.Command, args []string, out io.Writer) error {
	if api.ContainsElem(args, "options"){
		cmd.HelpFunc()(cmd, args)
		return nil
	}
	if len(args) > 0{
		return fmt.Errorf("invalid argument to get version.\n " +
			`Use "knitctl version --help" for more information about a given command.`)

	}
	strRe, err := api.GetKnitterVersion()
	api.ExcInfo(strRe, out)
	return err
}
