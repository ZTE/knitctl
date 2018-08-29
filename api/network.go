package api

import (
	"fmt"
	"net/http"
	"bytes"
	"strings"
	"encoding/json"
)

type struNetwork struct{
    BaseNetwork *BaseObject          //base information


    //Public bool                        //marked network is public ,private or protected
    //Gateway string                     //gateway of pod
    //Cidr string                        //pod ip
    Host_route []StruHostRoute

}
type StruHostRoute struct {
	Destination string
	Nexthop string
}

type putNwAll struct {
	Network putNetwork
}
type putNetwork struct {
	Host_routes []StruHostRoute
}

//get request network reponse structure
type struNetworkRe struct{
	Name string
	Network_id string
	Gateway string
	Cidr string
	Create_time string
	State string
	Public bool
	External bool
	Owner string
	Description string
	Subnet_id string
	Allocation_pools string

}

var(
	postNWUri = "nw/v1/tenants/{user}/networks"
	getNWUri = "nw/v1/tenants/{user}/networks/{network_uuid}"
	getNWByNameUri = "nw/v1/tenants/{user}/networks?name="
	putNWUri = "nw/v1/tenants/{user}/networks/{network_uuid}"
	deleteNWUri = "nw/v1/tenants/{user}/networks/{network_uuid}"
)

var (
	Head_nw = []string{"NAME", "CREATE_TIME", "GATEWAY", "CIDR", "PUBLIC", "OWNER", "STATE",
		"NETWORK_ID", "SUBNET_ID", "EXTERNAL", "DESCRIPTION", "ALLOCATION_POOLS"}
)

func NewStruNetwork(mapData map[string]interface{}) *struNetwork{

	kind := "network"
	version, found1 := mapData["apiversion"]
	if !found1{
		version = "v1"
	}

	mapMetaData := mapData["metadata"].(map[string]interface{})

	mapSpec := make(map[string]interface{})
	if spec, found := mapData["spec"]; found {
		mapSpec = spec.(map[string]interface{})
	}

	baseInfo := BaseObject{
		Kind: kind,
		Apiversion: version.(string),
		MetaData: mapMetaData,
		Spec: mapSpec,
	}

	return &struNetwork{
		BaseNetwork: &baseInfo,
	}
}

func (network *struNetwork) Post() (string, error) {
	hostport, err := GetHostPortFromPerporty()
	if (err != nil) {
		return "", err
	}
	uri := strings.Replace(postNWUri, "{user}", network.BaseNetwork.MetaData["tenant"].(string), -1)
	postUrl := hostport + "/" + uri

	//if number, _ :=network.Get(false); (number != nil && len(number) > 0) {
	//	return fmt.Errorf("network %s in tenant %s has been exist.\n", network.BaseNetwork.MetaData["name"].(string), network.BaseNetwork.MetaData["tenant"].(string))
	//}
	nwPublic := "false"
	if public, found := network.BaseNetwork.Spec["public"]; found && public == "true"{
		nwPublic = "true"
	}
	//{"network": {"name": "network_bar","public":false, "gateway": "123.124.125.1", "cidr": "123.124.125.0/24"}}
	//construct post body
	postData := `{"network": {"name": "`+ network.BaseNetwork.MetaData["name"].(string) + `","public": ` + nwPublic +
		`, "gateway": "` + network.BaseNetwork.Spec["gateway"].(string) + `", "cidr": "` + network.BaseNetwork.Spec["cidr"].(string) + `"}}`

	res, err := http.Post(postUrl, "application/json", bytes.NewReader([]byte(postData)))
	if (err != nil) {
		return "", err
	}

	_, erro := ExcuteNWResponse(res, postUrl, http.MethodPost)
	if (erro != nil) {
		return "", erro
	}

	strRe := fmt.Sprintf("create network '%s' in tenant '%s' successfully.\n", network.BaseNetwork.MetaData["name"].(string), network.BaseNetwork.MetaData["tenant"].(string))

	return strRe, nil
}

func GetNameById(tenant string, network_id string) (string, error) {
	hostport, err := GetHostPortFromPerporty()
	if (err != nil) {
		return "", err
	}

	if tenant == "" {
		return "", fmt.Errorf("get network_name by netowrk_id failure, tenant is empty.")
	}
	if network_id == "" {
		return "", fmt.Errorf("get network_name by netowrk_id failure, network_id is empty.")
	}
	uri := strings.Replace(getNWUri, "{user}", tenant, -1)
	uri = strings.Replace(uri, "{network_uuid}", network_id, -1)

	getUrl := hostport + "/" + uri

	response, err := http.Get(getUrl)
	// get all result
	arr, erro := ExcuteNWResponse(response, getUrl, http.MethodGet)
	if (erro != nil) {
		return "", erro
	}
	if len(arr) == 1 {
		return  arr[0].Name, nil
	}
	return "", fmt.Errorf("get network_name by netowrk_id failure, no found network_name by netowrk_id.")
}

func (network *struNetwork) Get(outputKind string) ([]struNetworkRe, string, error) {
	hostport, err := GetHostPortFromPerporty()
	if (err != nil) {
		return nil, "", err
	}

	uri := strings.Replace(getNWUri, "{user}", network.BaseNetwork.MetaData["tenant"].(string), -1)
	networkIdOrName := network.BaseNetwork.MetaData["name"].(string)
	getUriAll := strings.Replace(uri, "/{network_uuid}", "", -1)
	//support id and name


	getUrlAll := hostport + "/" + getUriAll
	response, err := http.Get(getUrlAll)
	// get all result
	arr, erro := ExcuteNWResponse(response, getUrlAll, http.MethodGet)
	if (erro != nil) {
		return arr, "", erro
	}
	// find one result by name or id
	arrRe := getNwResult(arr, networkIdOrName)

	strRe := printNws(arrRe, outputKind)


	return arrRe, strRe, nil
}

func (network *struNetwork) Put() (string, error) {
	if network.Host_route == nil || len(network.Host_route) == 0 {
		return "", fmt.Errorf("no set host_route to update.")
	}

	//no completed
	hostport, err := GetHostPortFromPerporty()
	if (err != nil) {
		return "", err
	}
	uri := strings.Replace(putNWUri, "{user}", network.BaseNetwork.MetaData["tenant"].(string), -1)
	getResult, _, _ := network.Get("")
	if getResult == nil || len(getResult) == 0{
		return "", fmt.Errorf("no found network resource to update.")
	}
	putUri := strings.Replace(uri, "{network_uuid}", getResult[0].Network_id, -1)
	putUrl := hostport + "/" + putUri

	putNw := putNetwork{
		Host_routes: network.Host_route,
	}
	putNwAll := putNwAll{
		Network: putNw,
	}
	putBody, _ := json.Marshal(putNwAll)

	res, err := RequestCommon(putUrl, http.MethodPut, putBody)
	if _, err := ExcuteNWResponse(res, putUrl, http.MethodPut); err != nil{
		return "", err
	}
	strRe := fmt.Sprintln("update network '%s' in tenant '%s' successfully.", network.BaseNetwork.MetaData["name"].(string), network.BaseNetwork.MetaData["tenant"].(string))
	return strRe, nil
}

func (network *struNetwork) Delete() (string, error) {
	hostport, err := GetHostPortFromPerporty()
	if (err != nil) {
		return "", err
	}
	uri := strings.Replace(deleteNWUri, "{user}", network.BaseNetwork.MetaData["tenant"].(string), -1)
	networkIdOrName := network.BaseNetwork.MetaData["name"].(string)
	if networkIdOrName == ""{
		return "", fmt.Errorf("no set network name(id) to delete.")
	}
	// judge name if resource has been exist
	arrNetworkRe, _, _ :=network.Get("")
	if (arrNetworkRe == nil || len(arrNetworkRe) == 0) {
		return "", fmt.Errorf("network '%s' in tenant '%s' is not exist, which can not be deleted.\n", network.BaseNetwork.MetaData["name"].(string), network.BaseNetwork.MetaData["tenant"].(string))
	}
	//onle delete network for id, not name
	networkUUid := arrNetworkRe[0].Network_id
	uri = strings.Replace(uri, "{network_uuid}", networkUUid, -1)

	deleteUrl := hostport + "/" + uri
	request, err := http.NewRequest("DELETE", deleteUrl, nil)
	if (err != nil) {
		return "", err
	}
	//timeout := flag.Int("o", 5, "-o N");
	Client := http.Client{}
	response, erroDo := Client.Do(request)
	if erroDo != nil {
		return "", erroDo
	}
	_, erro := ExcuteNWResponse(response, deleteUrl, "DELETE")
	if (erro != nil) {
		return "", erro
	}

	strRe := fmt.Sprintf("delete network '%s' in tenant '%s' successfully.\n", network.BaseNetwork.MetaData["name"].(string), network.BaseNetwork.MetaData["tenant"].(string))

	return strRe, nil
}

func ExcuteNWResponse(res *http.Response, requestUrl string, method string) ([]struNetworkRe, error) {
	strJson := GetBodyFromResponse(res.Body)

	if res.StatusCode > 300 || res.StatusCode < 200{
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

	if strJson == "" || strJson == "\"\""{
		return nil, nil
	}

	mapFind := make(map[string]interface{})
	erro := json.Unmarshal([]byte(strJson), &mapFind)
	if erro != nil{
		return nil, erro
	}
	if method == "DELETE" || method == "POST"{
		//delete and post do not print result
		return nil, nil
	}

	arrMapResult := []struNetworkRe{}
	if _, found := mapFind["networks"]; found{
		//it means that only one array result
		mapResult := make(map[string][]struNetworkRe)
		json.Unmarshal([]byte(strJson), &mapResult)
		arrMapResult = mapResult["networks"]

	}else if _, found := mapFind["network"]; found{
		//it means that only one result
		mapResult := make(map[string]struNetworkRe)
		json.Unmarshal([]byte(strJson), &mapResult)
		arrMapResult = []struNetworkRe{mapResult["network"]}
	}

	return arrMapResult, nil

}

func getNwResult(result []struNetworkRe, nameOrId string) []struNetworkRe{
	if result == nil {
		return nil
	}
	if len(result) == 0 {
		return result
	}
	if nameOrId == "" {
		// print all result array
		return result
	}
	for _, nw := range result{
		if nw.Name == nameOrId || nw.Network_id == nameOrId{
			//print this ipgroup
			return []struNetworkRe{nw}
		}
	}
	return nil
}


func printNws(result []struNetworkRe, outputKind string) string {
	if result == nil || len(result) == 0{
		strRe := fmt.Sprintln("no resource found.")
		return strRe
	}
	strArrRe := [][]string{Head_nw}
	//"NAME", "CREATE_TIME", "GATEWAY", "CIDR", "PUBLIC", "OWNER", "STATE",
	//"NETWORK_ID", "SUBNET_ID", "EXTERNAL", "DESCRIPTION", "ALLOCATION_POOLS
	for _, oneRe := range result{
		public := "false"
		if oneRe.Public{
			public = "true"
		}
		external := "false"
		if oneRe.External{
			external = "true"
		}
		oneStrArr := []string{oneRe.Name, oneRe.Create_time, oneRe.Gateway, oneRe.Cidr, public, oneRe.Owner,
			oneRe.State, oneRe.Network_id, oneRe.Subnet_id, external, oneRe.Description, oneRe.Allocation_pools}
		strArrRe = append(strArrRe, oneStrArr)
	}
	printHead :=[]int{0,1,2,3,4,5,6}
	if outputKind != "wide" {
		return  PrintArray(strArrRe, printHead)
	}else{
		return  PrintArray(strArrRe, nil)
	}
}