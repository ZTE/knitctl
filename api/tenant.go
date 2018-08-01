package api

import (
	"bytes"
	"fmt"
	"net/http"
	"encoding/json"
	"strings"
)

type struTenant struct{
	BaseInfo *BaseObject          //base information

}

type struTenantRe struct{
	Name string
	Id string
	Net_number int
	Created_at string
	Quotas map[string]interface{}
	Status string
	Iaas_tenant map[string]interface{}
}

var(
	pTenantUri = "nw/v1/tenants/{tenant-name}"
	gTenantUri = "nw/v1/tenants/{tenant-name}"
	putTenantUri = "nw/v1/tenants/{tenant-name}/quota?value="
	dTenantUri = "nw/v1/tenants/{tenant-name}"
)

var (
	Head_tenant = []string{"NAME", "ID", "NET_NUMBER", "CREATED_AT", "STATUS", "QUOTA"}
)

//construct strutenant structure
func NewStruTenant(mapData map[string]interface{}) *struTenant{
	kind := "tenant"
	version, found1 := mapData["apiversion"]
	if !found1{
		version = "v1"
	}

	metaData := mapData["metadata"]
	mapMetaData := make(map[string]interface{})
	if metaData != nil{
		mapMetaData = metaData.(map[string]interface{})
	}
	spec := mapData["spec"]
	mapSpec := make(map[string]interface{})
	if spec != nil{
		mapSpec = spec.(map[string]interface{})
	}

	baseInfo := BaseObject{
		Kind: kind,
		Apiversion: version.(string),
		MetaData: mapMetaData,
		Spec: mapSpec,
	}


	return &struTenant{
		BaseInfo: &baseInfo,

	}
}

func (t *struTenant) Post() (string, error) {
	hostport, err := GetHostPortFromPerporty()
	if (err != nil) {
		return "", err
	}
	//if count, _ := t.Get(false); (count != nil && len(count) >= 1 && count[0].Status == "ACTIVE") {
	//	return fmt.Errorf("tenant %s is exist and active.\n", t.BaseInfo.MetaData["name"].(string))
	//}
	uri := strings.Replace(pTenantUri, "{tenant-name}", t.BaseInfo.MetaData["name"].(string), -1)
	postUrl := hostport + "/" + uri
	postData := ""
	res, err := http.Post(postUrl, "application/json", bytes.NewReader([]byte(postData)))
	if (err != nil) {
		return "", err
	}

	_, erro := excuteTenantResponse(res, postUrl, "POST")
	if (erro != nil) {
		return "", erro
	}

	strRe := fmt.Sprintf("create tenant '%s' successfully.\n", t.BaseInfo.MetaData["name"].(string))

	return strRe, nil
}

func (t *struTenant) Get(outputKind string) ([]struTenantRe, string, error) {
	hostport, err := GetHostPortFromPerporty()
	if (err != nil) {
		return nil, "", err
	}
	nameOrId := ""
	if name, found := t.BaseInfo.MetaData["name"]; found {
		nameOrId = name.(string)
	}
	uri := strings.Replace(gTenantUri, "{tenant-name}", "", -1)
	getUrl := hostport + "/" + uri
	response, err := http.Get(getUrl)
	if (err != nil) {
		return nil, "", err
	}

	count, erro := excuteTenantResponse(response, getUrl, "GET")
	if (erro != nil) {
		return count, "", erro
	}

	getResult := getResultTenant(count, nameOrId)


	return getResult, printTenants(getResult, outputKind), nil
}

func (t *struTenant) Put() (string, error) {
	hostport, err := GetHostPortFromPerporty()
	if (err != nil) {
		return "", err
	}
	//put url
	uri := strings.Replace(putTenantUri, "{tenant-name}", t.BaseInfo.MetaData["name"].(string), -1)

	if quota, found := t.BaseInfo.Spec["quota"]; found {
		uri = (uri + quota.(string))
	}
	putUrl := hostport + "/" + uri

	response, err := RequestCommon(putUrl, http.MethodPut, []byte(""))

	if (err != nil) {
		return "", err
	}

	_, erro := excuteTenantResponse(response, putUrl, http.MethodPut)
	if (erro != nil) {
		return "", erro
	}

	strRe := fmt.Sprintf("update tenant '%s' quota successfully.\n", t.BaseInfo.MetaData["name"].(string))

	return strRe, nil
}

func (t *struTenant) Delete() (string, error) {
	hostport, err := GetHostPortFromPerporty()
	if (err != nil) {
		return "", err
	}

	//judge whether this tenant is exist
	//if count, _ := t.Get(false); (count == nil || len(count) == 0 ){
	//	return fmt.Errorf("tenant %s is not exist, which can not be deleted.\n", t.BaseInfo.MetaData["name"].(string))
	//}

	uri := strings.Replace(dTenantUri, "{tenant-name}", t.BaseInfo.MetaData["name"].(string), -1)
	deleteUrl := hostport + "/" + uri

	response, err := RequestCommon(deleteUrl, http.MethodDelete, []byte(""))
	if err != nil {
		return "", err
	}

	_, erro := excuteTenantResponse(response, deleteUrl, http.MethodDelete)
	if (erro != nil) {
		return "", erro
	}

	strRe := fmt.Sprintf("delete tenant '%s' successfully.\n", t.BaseInfo.MetaData["name"].(string))

	return strRe, nil
}

func excuteTenantResponse(res *http.Response, requestUrl string, method string) ([]struTenantRe, error) {
	strJson := GetBodyFromResponse(res.Body)
	if res.StatusCode > 300 || res.StatusCode < 200{
		//response code error
		errRes := ErrResponse{}
		err := json.Unmarshal([]byte(strJson), &errRes)
		errMessage := ""
		if err == nil {
			errMessage = errRes.Message
		}
		return nil, fmt.Errorf("%s request %s failure. status code: %d. message: %s", method, requestUrl, res.StatusCode, errMessage)
	}

	if strJson == "" && method == http.MethodGet{
		return nil, fmt.Errorf("no resource found.")
	}

	if method != http.MethodGet{
		//no getting method, no printing result
		return  nil, nil
	}

	mapFind := make(map[string]interface{})
	erro := json.Unmarshal([]byte(strJson), &mapFind)
	if erro != nil{
		return nil, erro
	}

	defer res.Body.Close()

	arrTenantRe := []struTenantRe{}
	if _, found := mapFind["tenant"]; found{
		//it means that only one result
		mapResult := make(map[string]struTenantRe)
		err := json.Unmarshal([]byte(strJson), &mapResult)
		if err == nil {
			arrTenantRe = []struTenantRe{mapResult["tenant"]}
		}else {
			mapResults := make(map[string][]struTenantRe)
			err := json.Unmarshal([]byte(strJson), &mapResults)
			if err != nil{
				return nil, err
			}
			arrTenantRe = mapResults["tenant"]
		}

	}

	return arrTenantRe, nil

}


func printTenants(result []struTenantRe, outputKind string) string {
	if(result == nil || len(result) <= 0){
		return ("no resource found.")

	}
	arrResult := [][]string{Head_tenant}
	//"NAME", "ID", "NET_NUMBER", "CREATED_AT", "STATUS", "QUOTA", "IAAS_TENANT"
	for _, oneRe := range result{
		strNumber := fmt.Sprintf("%d", oneRe.Net_number)
		strQuota, _ := json.Marshal(oneRe.Quotas)
		//strIaaSTenant, _ := json.Marshal(oneRe.Iaas_tenant)
		oneArr := []string{oneRe.Name, oneRe.Id, strNumber, oneRe.Created_at,
			oneRe.Status, string(strQuota)}
		arrResult = append(arrResult, oneArr)
	}
	//only print result NAME(column 0), ID(column 1), NET_NUMBER(column 2), CREATED_AT(column 3), STATUS(column 4)
	printHeadIndex := []int{0, 1, 2, 3, 4}
	if outputKind != "wide" {
		return  PrintArray(arrResult, printHeadIndex)
	}
	return  PrintArray(arrResult, nil)

}

func getResultTenant(result []struTenantRe, nameOrId string) []struTenantRe {
	if result == nil {
		return  nil
	}
	if len(result) == 0 {
		return result
	}
	if nameOrId == "" {
		// print all result array
		return result
	}
	for _, t := range result{
		if t.Name == nameOrId || t.Id == nameOrId{
			//print this ipgroup
			return []struTenantRe{t}
		}
	}
	return nil
}