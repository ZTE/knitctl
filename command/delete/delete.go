package delete
import (
	"github.com/spf13/cobra"
	"io"
	"github.com/ZTE/knitctl/api"
	"github.com/ZTE/knitctl/command/common"
	"fmt"
)

var (
	delete_short = ("Delete command by filenames, command and names.")
	delete_long = (`Delete command by filenames, command and names.
YAML format are accepted.
`)

	delete_example = (`  # Delete a network/tenant using the type and name specified in knitter.yaml.
  knitctl delete -f ./knitter.yaml

  # Delete network with name in tenant
  knitctl delete network my-network --tenant=my-user

  # Delete ipgroup with name in tenant
  knitctl delete ipgroup my-ipg --tenant=my-user

  # Delete tenant with name.
  kubectl delete tenant my-user`)
)

type DeleteOptions struct {
	Out    io.Writer
	ErrOut io.Writer

	FilenameOptions common.FileRecusiveOptions

	tenant string
}

func NewCmdDelete(out io.Writer, errOut io.Writer) *cobra.Command {
	options := DeleteOptions{
		Out : out,
		ErrOut : errOut,
	}
	cmd := &cobra.Command{
		Use: "delete [flags])",
		DisableFlagsInUseLine: true,
		Short:   delete_short,
		Long:    delete_long,
		Example: delete_example,
		Run: func(cmd *cobra.Command, args []string) {
			if err := options.RunDelete(cmd, args); err != nil {
				api.ExcuteError(err, options.ErrOut)
			}
		},
		SuggestFor: []string{"rm", "remove"},
		ValidArgs:  []string{"network", "tenant", "ipgroup"},
		ArgAliases: []string{"nw", "ipg"},
	}

	cmd.Flags().BoolVarP(&options.FilenameOptions.Recursive, "recursive", "R", options.FilenameOptions.Recursive, "Process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.")
	cmd.Flags().StringArrayVarP(&options.FilenameOptions.Filenames, "filename", "f", options.FilenameOptions.Filenames, "Filename, or directory to use to delete resources.")
	cmd.MarkFlagRequired("filename")

	cmd.AddCommand(NewCmdDeleteNW(out, errOut))
	cmd.AddCommand(NewCmdDeleteTenant(out, errOut))
	cmd.AddCommand(NewCmdDeleteIPG(out, errOut))



	//cmdutil.AddIncludeUninitializedFlag(cmd)
	return cmd
}

func (o *DeleteOptions) RunDelete(cmd *cobra.Command, args []string) error {
	if api.ContainsElem(args, "options"){
		cmd.HelpFunc()(cmd, args)
		return nil
	}
	if len(args) > 0{
		//delete -f aaa.yaml
		return fmt.Errorf("invalid argument '%s' to delete resource.\n " +
			`Use "knitctl delete --help" for more information about a given command.`, args)
	}
	fileNames := o.FilenameOptions.Filenames
	recursive := o.FilenameOptions.Recursive
	for _, oneFile := range fileNames {
		if (!recursive) {
			//it means no -r or --recursive flag
			return api.DeleteResource(oneFile, o.ErrOut, o.Out)
		} else {
			//it means that we must find all yaml file in fileName directory. then excute every file
			allFile := api.GetAllFileFormat(oneFile, ".yaml")
			for _, oneFile := range allFile {
				err := api.DeleteResource(oneFile, o.ErrOut, o.Out)
				if err != nil {
					return err
				}
			}

		}
	}

	return nil
}
