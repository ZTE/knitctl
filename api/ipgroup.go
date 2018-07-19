package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

//operator struct
type struIpgroup struct {
	base *BaseObject
}

//post struct
type struPostIpgroup struct {
	Ipgroup struIpgroupSpec
}
type struIpgroupSpec struct {
	Name       string
	Network_id string
	Ips        string
}

//response struct
type struIpgroupsRe struct {
	Name       string
	Id         string
	Network_id string
	Size       int
	Ips        []struIps
}

type struIps struct {
	Ip_addr string
	Used    bool
}

var (
	Ipgroup       = "ipgroup"
	postIpgroup   = "nw/v1/tenants/{user}/ipgroups"
	getIpgroup    = "nw/v1/tenants/{user}/ipgroups/{ipgroup-uuid}"
	putIpgroup    = "nw/v1/tenants/{user}/ipgroups/{ipgroup-uuid}"
	deleteIpgroup = "nw/v1/tenants/{user}/ipgroups/{ipgroup-uuid}"

	Head_ipg = []string{"NAME", "ID", "IPS", "NETWORK_ID", "SIZE"}
)

func NewStruIpgroup(mapYaml map[string]interface{}) *struIpgroup {
	return &struIpgroup{
		base: NewBaseInfo(mapYaml),
	}
}

func (ipgroup *struIpgroup) Post() (string, error) {
	hostport, err := GetHostPortFromPerporty()
	if err != nil {
		return "", err
	}
	tenant := ipgroup.base.MetaData["tenant"].(string)

	uri := strings.Replace(postIpgroup, "{user}", tenant, -1)
	postUrl := hostport + "/" + uri

	postJson := ipgroup.getPostPutJson(tenant)

	res, err := http.Post(postUrl, "application/json", bytes.NewReader([]byte(postJson)))
	if err != nil {
		return "", err
	}

	if _, err := executeIpgroupResponse(res, postUrl, http.MethodPost); err != nil {
		return "", err
	}
	strRe := fmt.Sprintf("create ipgroup '%s' in tenant '%s' successfully.\n", ipgroup.base.MetaData["name"].(string), tenant)
	//fmt.Printf("create ipgroup '%s' in tenant '%s' success.\n", ipgroup.base.MetaData["name"].(string), tenant)
	return strRe, nil

}

func (i *struIpgroup) Get(outputKind string) ([]struIpgroupsRe, string, error) {
	hostport, err := GetHostPortFromPerporty()
	if err != nil {
		return nil, "", err
	}
	strTenant := ""
	nameOrId := ""
	if name, found := i.base.MetaData["name"]; found {
		nameOrId = name.(string)
	}
	if tenant, found := i.base.MetaData["tenant"]; found {
		strTenant = tenant.(string)
	}
	geAlltUri := strings.Replace(getIpgroup, "{user}", strTenant, -1)

	geAlltUri = strings.Replace(geAlltUri, "/{ipgroup-uuid}", "", -1)

	getAllUrl := hostport + "/" + geAlltUri
	res, err := http.Get(getAllUrl)
	if err != nil {
		return nil, "", err
	}
	arrRe, err := executeIpgroupResponse(res, getAllUrl, http.MethodGet)
	if err != nil {
		return arrRe, "", err
	}

	getResult := getIpgroupRe(arrRe, nameOrId)

	return getResult, printIpgroup(getResult, outputKind), nil
}

func (i *struIpgroup) Put() (string, error) {
	hostport, err := GetHostPortFromPerporty()
	if err != nil {
		return "", err
	}
	strTenant := ""
	if tenant, found := i.base.MetaData["tenant"]; found {
		strTenant = tenant.(string)
	}
	putUri := strings.Replace(getIpgroup, "{user}", strTenant, -1)

	getArr, _, _ := i.Get("")
	if getArr == nil || len(getArr) == 0 {
		return "", fmt.Errorf("no found resouce to update.")
	}

	putUri = strings.Replace(putUri, "{ipgroup-uuid}", getArr[0].Id, -1)

	putUrl := hostport + "/" + putUri

	putData := i.getPostPutJson(strTenant)

	res, err := RequestCommon(putUrl, http.MethodPut, []byte(putData))
	if err != nil {
		return "", err
	}
	if _, err := executeIpgroupResponse(res, putUrl, http.MethodPut); err != nil {
		return "", err
	}
	strRe := fmt.Sprintf("update ipgroup '%s' in tenant '%s' successfully.\n", i.base.MetaData["name"].(string), strTenant)
	return strRe, nil
}

func (i *struIpgroup) Delete() (string, error) {
	hostport, err := GetHostPortFromPerporty()
	if err != nil {
		return "", err
	}
	strTenant := ""
	nameOrId := ""
	if name, found := i.base.MetaData["name"]; found {
		nameOrId = name.(string)
	}
	if tenant, found := i.base.MetaData["tenant"]; found {
		strTenant = tenant.(string)
	}
	tenantUri := strings.Replace(getIpgroup, "{user}", strTenant, -1)

	arrGet, _, err := i.Get("")
	if arrGet == nil || nameOrId == "" {
		return "", fmt.Errorf("no found ipgroup resource to delete. ")
	}
	strRe := ""
	for _, oneRe := range arrGet {
		deleteUri := strings.Replace(tenantUri, "{ipgroup-uuid}", oneRe.Id, -1)
		deleteUrl := hostport + "/" + deleteUri

		res, err := RequestCommon(deleteUrl, http.MethodDelete, []byte(""))
		if err != nil {
			return "", err
		}
		if _, err := executeIpgroupResponse(res, deleteUrl, http.MethodDelete); err != nil {
			return "", err
		}
		strRe += fmt.Sprintf("delete ipgroup name: '%s' in tenant '%s' successfully.\n", oneRe.Name, oneRe.Id, strTenant)
	}
	return strRe, nil
}

func (ipgroup *struIpgroup) getPostPutJson(tenant string) string {
	//get network_id from id or name
	ipgroupName := ""
	ipgroupNetworkid := ""
	ipgroupIps := ""
	if name, found := ipgroup.base.MetaData["name"]; found {
		ipgroupName = name.(string)
	}
	if networkid, found := ipgroup.base.MetaData["network"]; found {
		ipgNw := networkid.(string)
		ipgroupNetworkid = getNetworkIdByTenantAndName(tenant, ipgNw)
	}
	if ips, found := ipgroup.base.Spec["ips"]; found {
		ipgroupIps = "[" + ips.(string) + "]"
	}

	ipgroupSpec := struIpgroupSpec{
		Name:       ipgroupName,
		Network_id: ipgroupNetworkid,
		Ips:        ipgroupIps,
	}
	postMode := struPostIpgroup{
		Ipgroup: ipgroupSpec,
	}

	arr, _ := json.Marshal(postMode)
	return string(arr)
}

func executeIpgroupResponse(res *http.Response, requestUrl string, method string) ([]struIpgroupsRe, error) {
	statusCode := res.StatusCode
	strJson := GetBodyFromResponse(res.Body)

	if statusCode > 300 || statusCode < 200 {
		//response code error. get error message
		errRes := ErrResponse{}
		err := json.Unmarshal([]byte(strJson), &errRes)
		errMessage := ""
		if err == nil {
			errMessage = errRes.Message
		}
		return nil, fmt.Errorf("%s request %s failure. status code: %d. message: %s", method, requestUrl, res.StatusCode, errMessage)
	}

	defer res.Body.Close()

	if strJson == "" {
		return nil, nil
	}

	mapFind := make(map[string]interface{})
	erro := json.Unmarshal([]byte(strJson), &mapFind)
	if erro != nil {
		return nil, erro
	}

	arrMapResult := []struIpgroupsRe{}
	if _, found := mapFind["ipgroups"]; found {
		//it means that one array result
		mapResult := make(map[string][]struIpgroupsRe)
		json.Unmarshal([]byte(strJson), &mapResult)
		arrMapResult = mapResult["ipgroups"]

	} else if _, found := mapFind["ipgroup"]; found {
		//it means that only one result
		mapResult := make(map[string]struIpgroupsRe)
		json.Unmarshal([]byte(strJson), &mapResult)
		arrMapResult = []struIpgroupsRe{mapResult["ipgroup"]}
	}

	return arrMapResult, nil
}

func getIpgroupRe(result []struIpgroupsRe, nameOrId string) []struIpgroupsRe {
	if len(result) == 0 {
		return result
	}
	if nameOrId == "" {
		// print all result array
		return result
	}
	for _, ipgroup := range result {
		if ipgroup.Name == nameOrId || ipgroup.Id == nameOrId {
			//print this ipgroup
			return []struIpgroupsRe{ipgroup}
		}
	}
	return nil

}
func printIpgroup(result []struIpgroupsRe, outputKind string) string {
	if result == nil || len(result) == 0 {
		return ("the server doesn't have a resource type.")
	}

	resultArr := [][]string{Head_ipg}

	//"NAME", "ID", "IPS", "NETWORK_ID", "SIZE"
	for _, oneIpg := range result {
		jsonIps, _ := json.Marshal(oneIpg.Ips)
		oneArr := []string{oneIpg.Name, oneIpg.Id, string(jsonIps),
			oneIpg.Network_id, fmt.Sprintf("%d", oneIpg.Size)}
		resultArr = append(resultArr, oneArr)
	}
	printHead := []int{0, 2, 4}
	if outputKind != "wide" {
		return PrintArray(resultArr, printHead)
	}
	return PrintArray(resultArr, nil)

}
