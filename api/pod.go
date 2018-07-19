package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type struPod struct {
	base *BaseObject
}

type struPodRe struct {
	Ips  []struPodIps
	Name string
}

type struPodIps struct {
	Ip_address         string
	Network_plane      string
	Network_plane_name string
}

var (
	getPodUri = "nw/v1/tenants/{user}/pods/{pod-uuid}"
	Head_pod  = []string{"IPS", "NAME"}
)

func NewStruPod(mapData map[string]interface{}) *struPod {
	return &struPod{
		base: NewBaseInfo(mapData),
	}

}

func (p *struPod) Get() (string, error) {
	hostport, err := GetHostPortFromPerporty()
	if err != nil {
		return "", err
	}
	strTenant := ""
	nameOrId := ""
	if name, found := p.base.MetaData["name"]; found {
		nameOrId = name.(string)
	}
	if tenant, found := p.base.MetaData["tenant"]; found {
		strTenant = tenant.(string)
	}
	getPodUri := strings.Replace(getPodUri, "{user}", strTenant, -1)
	getPodUri = strings.Replace(getPodUri, "/{pod-uuid}", "", -1)

	getUrl := hostport + "/" + getPodUri
	resAll, err := http.Get(getUrl)
	if err != nil {
		return "", err
	}
	arrReAll, err := executePodResponse(resAll, getUrl, http.MethodGet)
	if err != nil {
		return "", err
	}
	arrRe := getPodResult(arrReAll, nameOrId)

	return printPod(arrRe), nil
}

func executePodResponse(res *http.Response, requestUrl string, method string) ([]struPodRe, error) {
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

	mapFind := make(map[string]interface{})
	erro := json.Unmarshal([]byte(strJson), &mapFind)
	if erro != nil {
		return nil, erro
	}

	arrMapResult := []struPodRe{}
	if _, found := mapFind["pods"]; found {
		//it means that one array result
		mapResult := make(map[string][]struPodRe)
		json.Unmarshal([]byte(strJson), &mapResult)
		arrMapResult = mapResult["pods"]

	} else if _, found := mapFind["pod"]; found {
		//it means that only one result
		mapResult := make(map[string]struPodRe)
		json.Unmarshal([]byte(strJson), &mapResult)
		arrMapResult = []struPodRe{mapResult["pod"]}
	}

	return arrMapResult, nil
}

func getPodResult(result []struPodRe, nameOrId string) []struPodRe {
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
	for _, pod := range result {
		if pod.Name == nameOrId {
			//print this ipgroup
			return []struPodRe{pod}
		}
	}
	return nil
}

func printPod(result []struPodRe) string {
	if result == nil || len(result) == 0 {
		return ("the server doesn't have a resource type.")

	}
	resultArr := [][]string{Head_pod}
	for _, onePod := range result {
		json, _ := json.Marshal(onePod.Ips)
		oneArr := []string{string(json), onePod.Name}
		resultArr = append(resultArr, oneArr)
	}
	printHead := []int{0, 1}
	return PrintArray(resultArr, printHead)

}
