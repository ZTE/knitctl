package api
/*
Copyright 2018 The ZTE Authors.

Deal with error.
*/

import (
	"fmt"
	"strings"
	"io"
)

func ExcInfo(info string, out io.Writer) {
	if len(info) > 0 {
		// add line in the end of errInfo
		if !strings.HasSuffix(info, "\n") {
			info += "\n"
		}
		out.Write([]byte(info))
	}
}

//Excute one error
func ExcuteError(err error, errOut io.Writer){
	if err == nil {
		return
	}
	errInfo := err.Error()
	if !strings.HasPrefix(errInfo, "error: ") {
		errInfo = fmt.Sprintf("error: %s", errInfo)
		if !strings.HasSuffix(errInfo, "\n") {
			errInfo += "\n"
		}
		errOut.Write([]byte(errInfo))
	}
	//ErrorPrint(errInfo)
}

//Excute error array
func ExcuteErrorArr(errs []error, errOut io.Writer){
	if len(errs) > 0 {
		for _, err := range errs{
			ExcuteError(err, errOut)
		}
	}
}
