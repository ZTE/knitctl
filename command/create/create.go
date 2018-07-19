package create

import (
	"io"

	"fmt"
	"github.com/ZTE/knitctl/api"
	"github.com/ZTE/knitctl/command/common"
	"github.com/spf13/cobra"
)

type CreateOptions struct {
	//PrintFlags  *PrintFlags
	//RecordFlags *genericclioptions.RecordFlags

	FilenameOptions  common.FileRecusiveOptions
	EditBeforeCreate bool
	Out              io.Writer
	ErrOut           io.Writer

	//Recorder genericclioptions.Recorder
	//PrintObj func(obj kruntime.Object) error
}

var (
	createShort = ("Create resource(s) from a file or stdin.")
	createLong  = (`Create resource(s) from a file.
YAML formats are accepted.`)

	createExample = (`  # Create resource(s) using the data in network.yaml.
  knitctl create -f ./resource.yaml`)
)

func NewCreateOptions(out, errOut io.Writer) *CreateOptions {
	return &CreateOptions{
		//PrintFlags:  NewPrintFlags("created"),
		//RecordFlags: genericclioptions.NewRecordFlags(),

		Out:    out,
		ErrOut: errOut,
	}
}

func NewCommandCreate(out io.Writer, errOut io.Writer) *cobra.Command {
	o := NewCreateOptions(out, errOut)

	cmd := &cobra.Command{
		Use: "create -f FILENAME",
		DisableFlagsInUseLine: true,
		Short:   createShort,
		Long:    createLong,
		Example: createExample,
		Run: func(cmd *cobra.Command, args []string) {
			api.ExcuteError(o.RunCreate(cmd, args), o.ErrOut)
		},
		SuggestFor: []string{"post"},
		ValidArgs:  []string{"network", "tenant", "ipgroup"},
		ArgAliases: []string{"nw", "ipg"},
	}

	//cmdutil.AddFilenameOptionFlags(cmd, &o.FilenameOptions, usage)
	cmd.Flags().BoolVarP(&o.FilenameOptions.Recursive, "recursive", "R", o.FilenameOptions.Recursive, "Process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.")
	cmd.Flags().StringArrayVarP(&o.FilenameOptions.Filenames, "filename", "f", o.FilenameOptions.Filenames, "Filename, or directory to create resources.")
	cmd.MarkFlagRequired("filename")

	cmd.AddCommand(NewCmdCreateNW(out, errOut))
	cmd.AddCommand(NewCmdCreateTenant(out, errOut))
	cmd.AddCommand(NewCmdCreateIpgroup(out, errOut))
	return cmd
}

func (o *CreateOptions) RunCreate(cmd *cobra.Command, args []string) error {
	if api.ContainsElem(args, "options") {
		cmd.HelpFunc()(cmd, args)
		return nil
	}
	if len(args) == 0 {
		//it means that create network/tenant based on yaml file(-f a.yaml)
		fileNames := o.FilenameOptions.Filenames
		recursive := o.FilenameOptions.Recursive
		for _, oneFile := range fileNames {
			if !recursive {
				//it means no -r or --recursive flag
				return api.CreateResource(oneFile, o.ErrOut, o.Out)
			} else {
				//it means that we must find all yaml file in fileName directory. then excute every file
				allFile := api.GetAllFileFormat(oneFile, ".yaml")
				for _, oneFile := range allFile {
					err := api.CreateResource(oneFile, o.ErrOut, o.Out)
					if err != nil {
						return err
					}
				}

			}
		}
	}

	if len(args) > 0 {
		err := fmt.Errorf("invalid argument '%s' to create resource.\n "+
			`Use "knitctl create --help" for more information about a given command.`, args)
		o.ErrOut.Write([]byte(err.Error()))
		return err
	}
	return nil
}
