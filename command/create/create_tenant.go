package create

import (
	"github.com/ZTE/knitctl/api"
	"github.com/ZTE/knitctl/command/common"
	"github.com/spf13/cobra"
	"io"
)

var (
	tenantShort = ("Create a tenant with the specified name.")
	tenantLong  = (`Create a tenant with the specified name.`)

	tenantExample = (`  # Create a new tenant named my-tenant
  knitctl create tenant my-tenant`)
)

// NewCmdCreateTenant is a macro command to create a new tenant
func NewCmdCreateTenant(cmdOut io.Writer, errOut io.Writer) *cobra.Command {
	options := &common.TenantOption{
		CmdOut: cmdOut,
		CmdErr: errOut,
	}

	cmd := &cobra.Command{
		Use: "tenant [NAME]",
		DisableFlagsInUseLine: true,
		Short:   tenantShort,
		Long:    tenantLong,
		Example: tenantExample,
		Run: func(cmd *cobra.Command, args []string) {
			api.ExcuteError(options.RunCreate(cmd, args), options.CmdErr)
		},
		SuggestFor: []string{"ten", "nant"},
	}

	//cmdutil.AddApplyAnnotationFlags(cmd)

	return cmd
}
