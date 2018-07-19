package common

import (
	"fmt"
	"github.com/ZTE/knitctl/api"
	"github.com/spf13/cobra"
	"io"
)

type TenantOption struct {
	CmdOut io.Writer
	CmdErr io.Writer

	Quota string

	//output
	Output string
}

func (o *TenantOption) RunCreate(cmd *cobra.Command, args []string) error {
	if api.ContainsElem(args, "options") {
		cmd.HelpFunc()(cmd, args)
		return nil
	}
	if len(args) == 0 {
		return fmt.Errorf("no set tenant name to create.\n" +
			"for example: knitctl create tenant tenant-name.")
	}
	mapPost := make(map[string]interface{})

	mapPost["kind"] = "tenant"
	mapPost["apiVersion"] = "v1"

	mapMeta := make(map[string]interface{})
	mapMeta["name"] = args[0]
	mapPost["metadata"] = mapMeta

	return api.CreateTN(mapPost, o.CmdOut)
}

func (o *TenantOption) RunDelete(cmd *cobra.Command, args []string) error {
	if api.ContainsElem(args, "options") {
		cmd.HelpFunc()(cmd, args)
		return nil
	}
	if len(args) == 0 {
		return fmt.Errorf("no set tenant name to delete.\n" +
			"for example: knitctl delete tenant tenant-name.")
	}

	return api.DeleteTNFromFiled(args[0], o.CmdOut)
}

func (o *TenantOption) RunGetTN(cmd *cobra.Command, args []string) error {
	//invoke api to get result and list it
	if len(args) >= 0 {
		nameOrId := ""
		if len(args) == 1 {
			nameOrId = args[0]
		}
		if o.Output != "" && o.Output != "wide" {
			return fmt.Errorf("invalid output options value. Supported value: wide.\n" +
				"for example: knitctl get ipgroup --tenant=my-tenant --output=wide.")
		}
		return api.GetTN(nameOrId, o.Output, o.CmdOut)

	}
	return nil
}

// Run performs the get operation.
func (options *TenantOption) RunSetTN(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		nameOrId := args[0]

		return api.PutTN(nameOrId, options.Quota, options.CmdOut)

	}
	return fmt.Errorf("no set tenant name to update.\n" +
		"for example: knitctl set tenant tenant-name [options].")
}
