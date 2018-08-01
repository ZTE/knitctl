package get

import (
	"github.com/spf13/cobra"
	"io"
	"github.com/ZTE/knitctl/api"
	"github.com/ZTE/knitctl/command/common"
)

var (
	getTenantDes = ("Display one or many tenant.")

	getTenantExample = (`  # List a tenant with name.
  knitctl get tenant my-user

  # List all tenants.
  knitctl get tenant
`)
)

func NewCmdGetTN(out io.Writer, errOut io.Writer) *cobra.Command {

	options := &common.TenantOption{
		CmdOut: out,
		CmdErr: errOut,
	}
	cmd := &cobra.Command{
		Use: "tenant [NAME]",
		Short:   getTenantDes,
		Long:    getTenantDes,
		Example: getTenantExample,
		Run: func(cmd *cobra.Command, args []string) {
			api.ExcuteError(options.RunGetTN(cmd, args), errOut)
		},
		SuggestFor: []string{"ten", "nant"},
	}

	cmd.Flags().StringVarP(&options.Output, "output", "o", options.Output, "If present, list all columns of network(s) in tenant name. Supported: wide.")

	return cmd
}
