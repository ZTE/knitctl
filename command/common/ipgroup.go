package common

import (
	"io"
	"github.com/spf13/cobra"
	"fmt"
	"github.com/ZTE/knitctl/api"
	"strconv"
)

type IpgroupOption struct {

	//required with ipgroup
	FSize int
	Tenant string
	Network_id string
	CmdOut io.Writer
	CmdErr io.Writer
	Ips string

	//output type
	Output string
}

func (ipg *IpgroupOption) RunCreate(cmd *cobra.Command, args []string) error {
	if api.ContainsElem(args, "options"){
		cmd.HelpFunc()(cmd, args)
		return nil
	}
	if len(args) == 0{
		return fmt.Errorf("no set ipgroup name to create.\n" +
			"for example: knitctl create ipgroup ipgroup-name [options].")
	}

	if ipg.FSize == 0 && ipg.Ips == "" {
		return fmt.Errorf("no set required flag either --ips or --size." )
	}

	mapPost := make(map[string]interface{})

	mapPost["kind"] = "ipgroup"
	mapPost["apiVersion"] = "v1"

	mapMeta := make(map[string]interface{})
	mapMeta["tenant"] = ipg.Tenant
	mapMeta["name"] = args[0]
	mapMeta["network"] = ipg.Network_id

	mapSpec := make(map[string]interface{})
	if ipg.Ips != ""{
		mapSpec["ips"] = ipg.Ips
		mapSpec["size"] = ""
	}else if (ipg.FSize > 0){
		mapSpec["size"] = strconv.Itoa(ipg.FSize)
		mapSpec["ips"] = ""
	}

	mapPost["spec"] = mapSpec
	mapPost["metadata"] = mapMeta

	return api.CreateIPG(mapPost, ipg.CmdOut)
}


func (ipg *IpgroupOption) RunDelete(cmd *cobra.Command, args []string) error {
	if api.ContainsElem(args, "options"){
		cmd.HelpFunc()(cmd, args)
		return nil
	}
	if len(args) == 0{
		return fmt.Errorf("no set ipgroup name to delete.\n" +
			"for example: knitctl delete ipgroup ipgroup-name [options].")
	}

	return api.DeleteIPG(ipg.Tenant, args[0], ipg.CmdOut)
}

// Run performs the get operation.
func (options *IpgroupOption) RunGetIpg(cmd *cobra.Command, args []string) error {
	if len(args) >= 0 {
		nameOrId := ""
		if len(args) == 1{
			nameOrId = args[0]
		}
		if options.Output != "" && options.Output != "wide"{
			return fmt.Errorf("invalid output options value. Supported value: wide.\n" +
				"for example: knitctl get ipgroup --tenant=my-tenant --output=wide.")
		}
		return api.GetIPG(options.Tenant, nameOrId, options.Output, options.CmdOut)

	}
	return nil
}

func (options *IpgroupOption) RunSetIpg(cmd *cobra.Command, args []string) error {
	//invoke api to get result and list it
	if len(args) > 0 {
		nameOrId := args[0]
		if options.FSize == 0 && options.Ips == "" {
			return fmt.Errorf("no set required flag either --ips or --size." )
		}
		ips := ""
		var strSize string = ""
		if options.Ips != ""{
			ips = options.Ips
		}else if (options.FSize > 0){
			strSize = strconv.Itoa(options.FSize)
		}
		return api.PutIPG(nameOrId, options.Tenant, ips, options.Network_id, strSize, options.CmdOut)
	}
	return fmt.Errorf("no set ipgroup name to update.\n" +
		"for example: knitctl set ipgroup ipgroup-name [options].")
}

