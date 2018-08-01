package common

import (
	"io"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ZTE/knitctl/api"
)

type NetworkOption struct {
	//post
	Public bool
	Gateway string
	Cidr string

	Tenant string

	//set host route
	Destination string
	Nexthop string

	CmdOut io.Writer
	CmdErr io.Writer

	//output
	Output string
}

func (o *NetworkOption) RunCreate(cmd *cobra.Command, args []string) error {
	if api.ContainsElem(args, "options"){
		cmd.HelpFunc()(cmd, args)
		return nil
	}
	if len(args) == 0{
		return fmt.Errorf("no set network name to create.\n" +
			"for example: knitctl create network network-name [options].")
	}
	mapPost := make(map[string]interface{})
	mapMeta := make(map[string]interface{})
	mapSpec := make(map[string]interface{})
	mapMeta["tenant"] = o.Tenant
	mapMeta["name"] = args[0]

	mapSpec["public"] = o.Public
	mapSpec["gateway"] = o.Gateway
	mapSpec["cidr"] = o.Cidr

	mapPost["kind"] = "network"
	mapPost["apiVersion"] = "v1"

	mapPost["metadata"] = mapMeta
	mapPost["spec"] = mapSpec

	return api.CreateNW(mapPost, o.CmdOut)
}

func (o *NetworkOption) RunDelete(cmd *cobra.Command, args []string) error {
	if api.ContainsElem(args, "options"){
		cmd.HelpFunc()(cmd, args)
		return nil
	}
	if len(args) == 0{
		return fmt.Errorf("no set network name to delete.\n" +
			"for example: knitctl delete network network-name [options].")
	}

	return api.DeleteNetworkFromFiled(o.Tenant, args[0], o.CmdOut)
}

func (options *NetworkOption) RunGetNW(cmd *cobra.Command, args []string) error {
	if len(args) >= 0 {
		nameOrId := ""
		if len(args) == 1{
			nameOrId = args[0]
		}
		if options.Output != "" && options.Output != "wide"{
			return fmt.Errorf("invalid output options value. Supported value: wide.\n" +
				"for example: knitctl get network --tenant=my-tenant --output=wide.")
		}
		return api.GetNetwork(options.Tenant, nameOrId, options.Output, options.CmdOut)

	}
	return nil
}

// Run performs the get operation.
func (options *NetworkOption) RunSetNw(cmd *cobra.Command, args []string) error {

	if len(args) > 1 {
		nameOrId := args[0]
		//set network
		strHostRoute := api.StruHostRoute{
			Destination: options.Destination,
			Nexthop:     options.Nexthop,
		}
		return api.PutNetwork(options.Tenant, nameOrId, []api.StruHostRoute{strHostRoute}, options.CmdOut)

	}
	return fmt.Errorf("no set network name to update.\n" +
		"for example: knitctl set network network-name [options].")
}
