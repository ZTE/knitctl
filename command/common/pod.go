package common

import (
	"fmt"
	"github.com/ZTE/knitctl/api"
	"github.com/spf13/cobra"
	"io"
)

type PodOption struct {
	Tenant string
	CmdOut io.Writer
	CmdErr io.Writer

	//output
	Output string
}

// Run performs the get operation.
func (options *PodOption) RunGetPod(cmd *cobra.Command, args []string) error {
	if len(args) >= 0 {
		nameOrId := ""
		if len(args) == 1 {
			nameOrId = args[0]
		}
		if options.Output != "" && options.Output != "wide" {
			return fmt.Errorf("invalid output options value. Supported value: wide.\n" +
				"for example: knitctl get ipgroup --tenant=my-tenant --output=wide.")
		}

		//only print two columns. output kind not work.
		return api.GetPod(options.Tenant, nameOrId, options.CmdOut)

	}
	return nil
}
