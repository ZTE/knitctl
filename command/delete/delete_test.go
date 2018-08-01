package delete

import (
	"bytes"
	"testing"
	"strings"
	"github.com/ZTE/knitctl/command/create"
	"time"
)

func TestDeleteByFile_success(t *testing.T) {
	// first create -f test-create.yaml
	bufCreate := bytes.NewBuffer([]byte{})
	cmdCreate := create.NewCommandCreate(bufCreate, bufCreate)
	cmdCreate.Flags().Set("filename", "../../test/test-delete.yaml")
	cmdCreate.Run(cmdCreate, []string{})

	time.Sleep(time.Second *5)

	buf := bytes.NewBuffer([]byte{})
	cmd := NewCmdDelete(buf, buf)
	cmd.Flags().Set("filename", "../../test/test-delete.yaml")
	cmd.Run(cmd, []string{})
	cmdResult := buf.String()
	expectResult := []string{"delete tenant 'delete_yaml_tenant_test' success",
		"delete network 'delete_yaml_network_test' in tenant 'delete_yaml_tenant_test' success.",
		"delete ipgroup name: 'delete_yaml_ipgroup_test'"}

	for _, oneExpect := range expectResult {
		if !strings.Contains(cmdResult, oneExpect) {
			t.Errorf("expect result: %s, but got: %s", oneExpect, cmdResult)
		}
	}

}


func TestDeleteByFile_failure(t *testing.T) {

	buf := bytes.NewBuffer([]byte{})
	cmd := NewCmdDelete(buf, buf)
	cmd.Flags().Set("filename_failure", "../../test/test-create.yaml")

	cmd.Run(cmd, []string{"aa"})

	cmdResult := buf.String()

	expectResult := []string{"invalid argument '[aa]' to delete resource."}

	for _, oneExpect := range expectResult {
		if !strings.Contains(cmdResult, oneExpect) {
			t.Errorf("expect result: %s, but got: %s", oneExpect, cmdResult)
		}
	}

	// error : no file
	cmd.Flags().Set("filename", "a.yaml")
	cmd.Run(cmd, []string{})
	cmdResult = buf.String()
	expectRe := "no such file or directory: a.yaml"
	if !strings.Contains(cmdResult, expectRe){
		t.Errorf("expect result: %s, but got: %s", expectRe, cmdResult)
	}

}
