package api

import (
	"net/http"
	"bytes"
)

type BaseObject struct{
	Kind string                      //kind(network and tenant supported)
	Apiversion string                //version(it seems not used)
	MetaData map[string]interface{}       //meteData information. ex: name ,tenant and so on
	Spec map[string]interface{}            //detail information of resource
}
type ErrResponse struct{
	ERROR string
	Message string
	Code string
}


type IRestful interface{
	Post(postData string) error
	Get() error
	Put() error
	Delete() error
}

func NewBaseInfo(mapData map[string]interface{}) *BaseObject {
	kind := Ipgroup
	version, found1 := mapData["apiversion"]
	if !found1{
		version = "v1"
	}
	mapMetaData := make(map[string]interface{})
	if metaData, found := mapData["metadata"]; found && metaData != nil{
		mapMetaData = metaData.(map[string]interface{})
	}
	mapSpec := make(map[string]interface{})
	if spec, found := mapData["spec"]; found && spec != nil {
		mapSpec = spec.(map[string]interface{})
	}

	return &BaseObject{
		Kind: kind,
		Apiversion: version.(string),
		MetaData: mapMetaData,
		Spec: mapSpec,
	}
}

func RequestCommon(url string, method string, body []byte) (resp *http.Response, err error) {
	request, err := http.NewRequest(method, url, bytes.NewReader(body))
	if (err != nil) {
		return nil, err
	}

	Client := http.Client{}
	return Client.Do(request)
}

func PrintArray(resultArr []([]string), headIndex []int) string {
	if resultArr == nil || len(resultArr) == 0 {
		return ""
	}
	//get max column size
	columnSize := make([]int, len(resultArr[0]))
	for _, oneArr := range resultArr{
		for ix, oneStr := range oneArr{
			if len(oneStr) > columnSize[ix] {
				columnSize[ix] = len(oneStr)
			}
		}
	}
	//headIndex means what columns to print
	//print string array
	strRe := ""
	for _, oneArr := range resultArr{
		if headIndex != nil && len(headIndex) > 0{
			//print column based on headindex
			for _, headIx := range headIndex {
				if headIx < len(oneArr) {
					printStr := (oneArr[headIx] + getStrSpaceOfCount(columnSize[headIx]-len(oneArr[headIx])) + "    ")
					strRe += printStr
				}
			}
		}else{
			//print all column
			for ix, oneStr := range oneArr{
				printStr := (oneStr + getStrSpaceOfCount(columnSize[ix]-len(oneStr)) + "    ")
				strRe += printStr
			}
		}
		strRe += "\n"
	}
	return strRe

}

func getStrSpaceOfCount(count int) string{
	str := ""
	for ix := 0; ix < count; ix = ix +1{
		str += " "
	}
	return  str
}
