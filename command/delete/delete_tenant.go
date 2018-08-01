package delete

import (
	"github.com/spf13/cobra"
	"io"
	"github.com/ZTE/knitctl/api"
	"github.com/ZTE/knitctl/command/common"
)

var (
	tenantDeleteShort = ("Delete a tenant with the specified name.")
	tenantDeleteLong = (`Delete a tenant with the specified name.`)

	tenantDeleteExample = (`  # Delete a tenant named my-tenant
  knitctl delete tenant my-tenant`)
)

// NewCmdCreateTenant is a macro command to create a new tenant
func NewCmdDeleteTenant(cmdOut io.Writer, errOut io.Writer) *cobra.Command {
	options := &common.TenantOption{
		CmdOut: cmdOut,
		CmdErr: errOut,
	}

	cmd := &cobra.Command{
		Use: "tenant [NAME]",
		DisableFlagsInUseLine: true,
		Aliases:               []string{"t"},
		Short:                 tenantDeleteShort,
		Long:                  tenantDeleteLong,
		Example:               tenantDeleteExample,
		Run: func(cmd *cobra.Command, args []string) {
			api.ExcuteError(options.RunDelete(cmd, args), options.CmdErr)
		},
		SuggestFor: []string{"ten", "nant"},
	}

	//cmdutil.AddApplyAnnotationFlags(cmd)

	return cmd
}
