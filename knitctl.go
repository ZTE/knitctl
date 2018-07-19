package main

import (
	"os"
	//"fmt"
	"github.com/ZTE/knitctl/api"
	"github.com/ZTE/knitctl/command"
	"fmt"
)
const (
	ErrorExitCode = 1
)

func main(){

	mapConfig := make(map[string]string)
	mapConfig["knitterurl"] = "http://127.0.0.1:9527"
	//set default knitter server
	if err := api.GenerateConfigFile(mapConfig); err != nil{
		fmt.Fprint(os.Stderr, (err.Error() + "\n"))
		os.Exit(ErrorExitCode)
	}
	cmd := command.NewKnitctlCommand(os.Stdout, os.Stderr)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(ErrorExitCode)
	}
}

