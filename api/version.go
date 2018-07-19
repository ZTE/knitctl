package api

const (
	versionUri = "nw/version"
)

func GetKnitterVersion() (string, error) {
	//hostport, err := GetHostPortFromPerporty()
	//if (err != nil) {
	//	return err
	//}
	//
	////versionUrl := hostport + "/" + versionUri
	////
	////res, err := http.Get(versionUrl)
	////if (err != nil) {
	////	return err
	////}
	////if res.StatusCode > 299 || res.StatusCode < 200{
	////	return fmt.Errorf("GET request %s failure, status code is %d.", versionUrl, res.StatusCode)
	////}
	////
	////info := make([]byte, 1024)
	////res.Body.Read(info)
	strRe := ("apiVersion: v1\n" + "go version:1.9.5")

	return strRe, nil
}
