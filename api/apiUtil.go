package api

/*
Copyright 2018 The ZTE Authors.

global function and varible.
*/

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"strings"
	"io"
)

const (
	APPName = "knitctl"
	ApplicationPropertyFile = "/etc/knitctl/config.yaml"
	ApplicationPropertyFolder = "/etc/knitctl/"
)
var (
	bodyCache = 1024 * 4
)

func GetFromYaml(yamlFile string) ([]map[string]interface{}, error){

	mapYaml := make(map[string]interface{})

	arrMapYaml := []map[string]interface{}{}

	find, err := FileExists(yamlFile)
	if err != nil {
		return []map[string]interface{}{mapYaml},err
	}
	if find == false {
		return []map[string]interface{}{mapYaml}, fmt.Errorf("no such file or directory: %s.\n", yamlFile)
	}
	//load yaml file
	data, err := ioutil.ReadFile(yamlFile)
	if (err != nil) {
		return []map[string]interface{}{mapYaml}, err
	}
	//analysis yaml file to map
	strSplit := "---"
	strFileInfo := strings.Split(string(data), strSplit)
	for _, oneInfo := range strFileInfo{
		mapJson := make(map[string]interface{})
		erro := yaml.Unmarshal([]byte(oneInfo), &mapJson)
		if(erro != nil){
			return arrMapYaml, erro
		}
		arrMapYaml = append(arrMapYaml, mapJson)
	}

	return arrMapYaml, nil
}

func GetHostPortFromPerporty() (string, error) {
	arrMapYaml, err := GetFromYaml(ApplicationPropertyFile)
	if(err != nil){
		return "", err
	}
	if len(arrMapYaml) <= 0{
		return "", fmt.Errorf("no found knitterurl field in file %s.\n", ApplicationPropertyFile)
	}
	if hostport, found := arrMapYaml[0]["knitterurl"]; found{
		return hostport.(string), nil
	}
	return "", fmt.Errorf("no found knitterurl field in file %s.\n", ApplicationPropertyFile)
}

func FileExists(filename string) (bool, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func GetAllFileFormat(dir string, format string) []string{
	if find, _ := FileExists(dir); find && strings.Contains(strings.ToLower(dir), format) {
		//if dir is a file supported format, return dir
		return []string{dir}
	}
	dir = strings.Trim(dir, " ")
	if len(dir) == 0{
		return []string{}
	}
	fileOut := []string{}
	files, _ := ioutil.ReadDir(dir)
	for _, oneFile := range files{
		if oneFile.IsDir(){
			subFile := GetAllFileFormat(dir + "/" + oneFile.Name(), format)
			for _, oneSubFile := range subFile {
				fileOut = append(fileOut, oneSubFile)
			}
		}
		if strings.Contains(strings.ToLower(oneFile.Name()), format){
			fileOut = append(fileOut, dir + "/" + oneFile.Name())
		}
	}
	return fileOut
}

func GenerateConfigFile(mapYaml map[string]string) error{

	if find,_ := FileExists(ApplicationPropertyFolder); !find {
		err := os.Mkdir(ApplicationPropertyFolder, 0777)
		if err != nil {
			return err
		}
	}
	find, err := FileExists(ApplicationPropertyFile)
	if err != nil{
		fmt.Errorf("generate config file %s failure.\n", ApplicationPropertyFile)
		return err
	}
	if find{
		return nil
	}
	return ModifyConfigFile(mapYaml)
}

func ModifyConfigFile(mapYaml map[string]string) error{
	if len(mapYaml) == 0{
		return nil
	}
	if find,_ := FileExists(ApplicationPropertyFolder); !find {
		err := os.Mkdir(ApplicationPropertyFolder, 0777)
		if err != nil {
			return err
		}
	}
	data, err:= yaml.Marshal(mapYaml)
	if err != nil{
		return err
	}
	file, err := os.OpenFile(ApplicationPropertyFile, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0777)
	if err != nil{
		return err
	}
	if _, err :=file.Write(data); err != nil{
		return err
	}

	return nil
}


func GetBodyFromResponse(reader io.ReadCloser) string {
	strRe := "";
	temp := make([]byte, bodyCache)
	byteCount, _ := reader.Read(temp)

	for ; byteCount > 0; {
		strRe += byteToString(temp, byteCount)
		temp = make([]byte, bodyCache)
		byteCount, _ = reader.Read(temp)
	}
	strRe = strings.Trim(strRe, "\"")
	strRe = strings.Trim(strRe, " ")
	strRe = strings.Trim(strRe, "\t")
	strRe = strings.Trim(strRe, "\n")
	return  strRe
}

func byteToString(temp []byte, byteCount int) string {
	jsonRe := []byte{}
	for index, oneByte := range temp{
		if index == byteCount{
			break
		}
		jsonRe = append(jsonRe, oneByte)
	}
	return string(jsonRe)
}

func ContainsElem(arr []string, elem string) bool {
	if len(arr) == 0 {
		return false
	}
	for _, oneElem := range arr{
		if oneElem == elem{
			return true
		}
	}
	return false
}