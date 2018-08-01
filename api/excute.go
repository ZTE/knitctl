package api

import (
	"fmt"
	"strings"
	"io"
)
//create resource based on yaml file
func CreateResource(yamlFile string, errOut io.Writer, out io.Writer) error {
	arrMapYaml, err := GetFromYaml(yamlFile)
	if(err != nil){
		return err
	}

	if len(arrMapYaml) > 0{
		//first excute creating tenant
		errs := []error{}
		for _, mapYaml := range arrMapYaml{
			kind, found := mapYaml["kind"]
			if found && kind == "tenant"{
				err := CreateTN(mapYaml, out)
				if err != nil{
					//only print ,do not return
					errs = append(errs, err)
				}
			}
		}

		for _, mapYaml := range arrMapYaml{
			kind, found := mapYaml["kind"]
			if found && kind == "network" {
				//it means that this is a network resource
				err := CreateNW(mapYaml, out)
				if err != nil {
					//only print ,do not return
					errs = append(errs, err)
				}
			}
		}

		for _, mapyamlNw := range arrMapYaml{
			kind, found := mapyamlNw["kind"]
			if !found {
				err :=fmt.Errorf("no found resource kind in file %s.\n", yamlFile)
				errs = append(errs, err)
			}else if kind == "ipgroup" {
				err := CreateIPG(mapyamlNw, out)
				if err != nil {
					errs = append(errs, err)
				}
			}else if kind != "tenant" && kind != "network"{
				err := fmt.Errorf("unknow resource kind %s in %s. supported kind: 'network' or 'tenant'.\n", mapyamlNw["kind"], yamlFile)
				errs = append(errs, err)
			}
		}
		ExcuteErrorArr(errs, errOut)

	}
	return nil
}

//delete resource based on yaml file
func DeleteResource(yamlFile string, errOut io.Writer, out io.Writer) error {
	arrMapYaml, err := GetFromYaml(yamlFile)
	if(err != nil){
		return err
	}

	if len(arrMapYaml) > 0{
		errs := []error{}
		//first excute deleting network
		for _, mapYaml := range arrMapYaml {
			kind, found := mapYaml["kind"]
			if found && kind == "ipgroup" {
				//it means that this is a ipgroup resource
				erro := DeleteIPGFromMap(mapYaml, out)
				if erro != nil {
					//only print ,do not return
					errs = append(errs, erro)
				}
			}
		}

		for _, mapYaml := range arrMapYaml{
			kind, found := mapYaml["kind"]
			if found && kind == "network"{
				//it means that this is a network resource
				erro := DeleteNetwork(mapYaml, out)
				if erro != nil{
					//only print ,do not return
					errs = append(errs, erro)
				}
			}
		}

		for _, mapYaml := range arrMapYaml{
			kind, found := mapYaml["kind"]
			if !found{
				errs = append(errs, fmt.Errorf("no found resource kind in file %s.\n", yamlFile))

			}else if kind == "tenant"{
				//it means that this is a tenant resource
				err0 := DeleteTN(mapYaml, out)
				if err0 != nil{
					//only print ,do not return
					errs = append(errs, err0)
				}
			}else if kind != "network" && kind != "ipgroup" {
				err := fmt.Errorf("unknow resource kind %s in file %s. supported type: 'network', 'tenant', and 'ipgroup'.\n", mapYaml["kind"], yamlFile)
				errs = append(errs, err)
			}
		}

		ExcuteErrorArr(errs, errOut)
	}
	return nil
}

//////////////////////////////// network function
func CreateNW(mapYaml map[string]interface{}, out io.Writer) error {
	if kind, found := mapYaml["kind"]; found && kind == "network"{
		metadata, fMetadata := mapYaml["metadata"]
		spec, fSpec := mapYaml["spec"]

		necessaryField := ""
		if !fSpec || spec == nil{
			necessaryField += "spec: {gateway}, spec: {cidr}, "
		}
		if !fMetadata || metadata == nil{
			necessaryField += "metadata: {name}, metadata: {tenant}, "
		}
		if fMetadata && metadata != nil {
			name, fName := metadata.(map[string]interface{})["name"]
			tenant, fTenant := (metadata.(map[string]interface{}))["tenant"]
			if !fName || name == nil{
				necessaryField += "metadata: {name}, "
			}
			if !fTenant || tenant ==nil {
				necessaryField += "metadata: {tenant}, "
			}
		}
		if fSpec && spec != nil {
			gateway, fGateway := (spec.(map[string]interface{}))["gateway"]
			cidr, fCidr := (spec.(map[string]interface{}))["cidr"]
			if !fGateway || gateway == nil{
				necessaryField += "spec: {gateway}, "
			}
			if !fCidr || cidr == nil {
				necessaryField += "spec: {cidr}, "
			}
		}
		if necessaryField != ""{
			necessaryField = strings.TrimRight(necessaryField, ", ")
			return  fmt.Errorf("no found necessary field '%s'.\n", necessaryField)
		}

		strNetwork := NewStruNetwork(mapYaml)
		strRe, err :=  strNetwork.Post()
		ExcInfo(strRe, out)
		return err
	}

	return fmt.Errorf("no found resource kind network")
}

func GetNetwork(tenant string, nameOrId string, outputKind string, out io.Writer) (error) {
	mapMetaData := make(map[string]interface{})
	mapMetaData["tenant"] = tenant
	mapMetaData["name"] = nameOrId

	mapYaml := make(map[string]interface{})
	mapYaml["metadata"] = mapMetaData
	mapYaml["kind"] = "network"
	mapYaml["apiversion"] = "v1"

	strNW := NewStruNetwork(mapYaml)
	_, strRe, err := strNW.Get(outputKind)
	ExcInfo(strRe, out)

	return err
}

// a interface
func getNetworkIdByTenantAndName(tenant string, nameOrId string) (string) {
	mapMetaData := make(map[string]interface{})
	mapMetaData["tenant"] = tenant
	mapMetaData["name"] = nameOrId

	mapYaml := make(map[string]interface{})
	mapYaml["metadata"] = mapMetaData
	mapYaml["kind"] = "network"
	mapYaml["apiversion"] = "v1"

	strNW := NewStruNetwork(mapYaml)
	arrRe, _, _ := strNW.Get("");
	if arrRe == nil || len(arrRe) == 0 {
		return ""
	}
	return arrRe[0].Network_id
}

func PutNetwork(tenant string, nameOrId string, hostRoute []StruHostRoute, out io.Writer) error {
	mapMetaData := make(map[string]interface{})
	mapMetaData["tenant"] = tenant
	mapMetaData["name"] = nameOrId

	mapYaml := make(map[string]interface{})
	mapYaml["metadata"] = mapMetaData

	strNW := NewStruNetwork(mapYaml)
	strNW.Host_route = hostRoute
	strRe, err := strNW.Put()

	ExcInfo(strRe, out)

	return err
}

func DeleteNetworkFromFiled(tenant string, nameOrId string, out io.Writer) error {
	mapMetaData := make(map[string]interface{})
	mapMetaData["tenant"] = tenant
	mapMetaData["name"] = nameOrId

	mapYaml := make(map[string]interface{})
	mapYaml["metadata"] = mapMetaData

	strNW := NewStruNetwork(mapYaml)
	str, err := strNW.Delete()
	ExcInfo(str, out)
	return err
}

func DeleteNetwork(mapYaml map[string]interface{}, out io.Writer) error {

	if kind, found := mapYaml["kind"]; found && kind == "network"{

		metadata, fMetadata := mapYaml["metadata"]

		necessaryField := ""
		if !fMetadata || metadata == nil{
			necessaryField += "metadata: {tenant}, metadata: {name}, "
		}else {
			name, fName := metadata.(map[string]interface{})["name"]
			tenant, fTenant := metadata.(map[string]interface{})["tenant"]
			if !fName || name == nil{
				necessaryField += "metadata: {name}, "
			}
			if !fTenant || tenant == nil{
				necessaryField += "metadata: {tenant}, "
			}
		}
		if necessaryField != ""{
			return  fmt.Errorf("no found necessary network field: '%s'.\n", necessaryField)
		}
		strNetwork := NewStruNetwork(mapYaml)
		strRe, err := strNetwork.Delete()
		ExcInfo(strRe, out)
		return err
	}

	return fmt.Errorf("no found resource kind network")
}


//////////////////////////////// Tenant function
func CreateTN(mapYaml map[string]interface{}, out io.Writer) error{
	if kind, found := mapYaml["kind"]; found && kind == "tenant"{
		metadata, fMetadata := mapYaml["metadata"]

		necessaryField := ""
		if !fMetadata || metadata == nil{
			necessaryField += "metadata: {name}"
			return fmt.Errorf("no found necessary field: '%s'.\n", necessaryField)
		}

		strTN := NewStruTenant(mapYaml)
		strRe , err := strTN.Post()
		ExcInfo(strRe, out)
		return err
	}

	return fmt.Errorf("no found resource kind tenant")
}

func GetTN(tenant string, outputKind string, out io.Writer) error {
	mapMeta := make(map[string]interface{})
	mapMeta["name"] = tenant

	mapYaml := make(map[string]interface{})
	mapYaml["metadata"] = mapMeta

	strTN := NewStruTenant(mapYaml)
	_, strRe, err := strTN.Get(outputKind)
	ExcInfo(strRe, out)
	return err
}

func DeleteTNFromFiled(tenant string, out io.Writer) error {
	mapTenant := make(map[string]interface{})
	mapTenant["name"] = tenant

	mapYaml := make(map[string]interface{})
	mapYaml["metadata"] = mapTenant

	strTN := NewStruTenant(mapYaml)
	strRe, err := strTN.Delete()
	ExcInfo(strRe, out)
	return err
}

func PutTN(name string, quota string, out io.Writer) error{
	mapMeta := make(map[string]interface{})
	mapMeta["name"] = name

	mapSpec := make(map[string]interface{})
	mapSpec["quota"] = quota

	mapPut := make(map[string]interface{})
	mapPut["metadata"] = mapMeta
	mapPut["spec"] = mapSpec

	strTN := NewStruTenant(mapPut)
	strRe, err := strTN.Put()
	ExcInfo(strRe, out)
	return err
}

func DeleteTN(mapYaml map[string]interface{}, out io.Writer) error {
	if kind, found := mapYaml["kind"]; found && kind == "tenant"{

		if metadata, fMeta := mapYaml["metadata"]; !fMeta || metadata == nil{
			necessary := "metadata: {name}"
			return  fmt.Errorf("no found necessary tenant field: '%s'.\n", necessary)
		}

		name, fName := mapYaml["metadata"].(map[string]interface{})["name"]

		necessaryField := "name"
		if !fName || name == nil{
			return  fmt.Errorf("no found necessary tenant field: '%s'.\n", necessaryField)
		}
		strTN := NewStruTenant(mapYaml)
		strRe, err := strTN.Delete()
		ExcInfo(strRe, out)
		return err
	}

	return fmt.Errorf("no found resource kind network")
}
//////////////////////////////////ipgroup
func CreateIPG(mapYaml map[string]interface{}, out io.Writer) error {
	if kind, found := mapYaml["kind"]; found && kind == "ipgroup"{
		metadata, fMetadata := mapYaml["metadata"]
		spec, fSpec := mapYaml["spec"]

		necessaryField := ""
		if !fSpec || spec == nil {
			necessaryField += "spec: {ips or size}, "
		}
		if !fMetadata || metadata == nil{
			necessaryField += "metadata: {name}, metadata: {tenant}, spec: {network}, "
		}
		if fMetadata && metadata != nil {
			_, fName := metadata.(map[string]interface{})["name"]
			_, fTenant := (metadata.(map[string]interface{}))["tenant"]
			_, fNw := (metadata.(map[string]interface{}))["network"]
			if !fName {
				necessaryField += "metadata: {name}, "
			}
			if !fTenant {
				necessaryField += "metadata: {tenant}, "
			}
			if !fNw{
				necessaryField += "metadata: {network}, "
			}
		}
		if fSpec && spec != nil {
			_, fIps := (spec.(map[string]interface{}))["ips"]
			_, fSize := (spec.(map[string]interface{}))["size"]
			if !fIps && !fSize {
				necessaryField += "spec: {ips or size}, "
			}
		}
		if necessaryField != ""{
			necessaryField = strings.TrimRight(necessaryField, ", ")
			return  fmt.Errorf("no found necessary field '%s'.\n", necessaryField)
		}

		strIPG := NewStruIpgroup(mapYaml)
		strRe, err := strIPG.Post()
		ExcInfo(strRe, out)
		return err
	}

	return fmt.Errorf("no found resource kind ipgroup")
}

func GetIPG(tenant string, nameOrId string, ouputKind string, out io.Writer) error {
	mapMetaData := make(map[string]interface{})
	mapMetaData["tenant"] = tenant
	mapMetaData["name"] = nameOrId

	mapYaml := make(map[string]interface{})
	mapYaml["metadata"] = mapMetaData

	strIpg := NewStruIpgroup(mapYaml)
	_, strRe, err := strIpg.Get(ouputKind)
	ExcInfo(strRe, out)
	return err
}

func PutIPG(name string, tenant string, ips string, nwid string, size string, out io.Writer) error{
	//check necessary field by options.
	mapPut := make(map[string]interface{})

	mapMeta := make(map[string]interface{})
	mapMeta["name"] = name
	mapMeta["tenant"] = tenant
	mapMeta["network"] = nwid

	mapSpec := make(map[string]interface{})
	mapSpec["ips"] = ips
	mapSpec["size"] = size

	mapPut["metadata"] = mapMeta
	mapPut["spec"] = mapSpec
	strIPG := NewStruIpgroup(mapPut)
	strRe, err := strIPG.Put()
	ExcInfo(strRe, out)
	return err
}

func DeleteIPG(tenant string, nameOrId string, out io.Writer) error {
	mapMetaData := make(map[string]interface{})
	mapMetaData["tenant"] = tenant
	mapMetaData["name"] = nameOrId

	mapYaml := make(map[string]interface{})
	mapYaml["metadata"] = mapMetaData

	mapYaml["kind"] = "ipgroup"
	mapYaml["apiVersion"] = "v1"

	strIpg := NewStruIpgroup(mapYaml)
	strRe, err := strIpg.Delete()
	ExcInfo(strRe, out)
	return err
}

func DeleteIPGFromMap(mapYaml map[string]interface{}, out io.Writer) error {
	if kind, found := mapYaml["kind"]; found && kind == "ipgroup"{

		metadata, fMetadata := mapYaml["metadata"]

		necessaryField := ""
		if !fMetadata || metadata == nil{
			necessaryField += "metadata: {tenant}, metadata: {name}, "
		}else {
			name, fName := metadata.(map[string]interface{})["name"]
			tenant, fTenant := metadata.(map[string]interface{})["tenant"]
			if !fName || name ==nil{
				necessaryField += "metadata: {name}, "
			}
			if !fTenant || tenant == nil{
				necessaryField += "metadata: {tenant}, "
			}
		}
		if necessaryField != ""{
			return  fmt.Errorf("no found necessary ipgroup field: '%s'.\n", necessaryField)
		}
		struIpg := NewStruIpgroup(mapYaml)
		strRe, err := struIpg.Delete()
		ExcInfo(strRe, out)
		return  err
	}

	return fmt.Errorf("no found resource kind ipgroup")
}
//////////////////////pod
func GetPod(tenant string, nameOrId string, out io.Writer) error {
	mapMetaData := make(map[string]interface{})
	mapMetaData["tenant"] = tenant
	mapMetaData["name"] = nameOrId

	mapYaml := make(map[string]interface{})
	mapYaml["metadata"] = mapMetaData

	struPod := NewStruPod(mapYaml)
	strRe, err := struPod.Get()
	ExcInfo(strRe, out)
	return err
}