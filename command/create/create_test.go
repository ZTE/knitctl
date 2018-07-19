package create

import (
	"bytes"
	"knitter.io/knitter/command/delete"
	"strings"
	"testing"
	"time"
)

func TestCreateByFile_success(t *testing.T) {

	buf := bytes.NewBuffer([]byte{})
	cmd := NewCommandCreate(buf, buf)
	cmd.Flags().Set("filename", "../../test/test-create.yaml")

	cmd.Run(cmd, []string{})

	cmdResult := buf.String()
	expectResult := []string{"create tenant 'create_yaml_tenant_test1' success",
		"create network 'create_yaml_network_test1' in tenant 'create_yaml_tenant_test1' success.",
		"create ipgroup 'create_yaml_ipgroup_test1' in tenant 'create_yaml_tenant_test1' success"}

	for _, oneExpect := range expectResult {
		if !strings.Contains(cmdResult, oneExpect) {
			t.Errorf("expect result: %s, but got: %s", oneExpect, cmdResult)
		}
	}

	time.Sleep(time.Second * 5)

	bufDelete := bytes.NewBuffer([]byte{})
	cmdDelete := delete.NewCmdDelete(bufDelete, bufDelete)
	cmdDelete.Flags().Set("filename", "../../test/test-create.yaml")
	cmdDelete.Run(cmdDelete, []string{})

}
func TestCreateByFile_failure(t *testing.T) {

	buf := bytes.NewBuffer([]byte{})
	cmd := NewCommandCreate(buf, buf)
	cmd.Flags().Set("filename_failure", "../../test/test-create.yaml")

	cmd.Run(cmd, []string{"aa"})

	cmdResult := buf.String()

	expectResult := []string{"invalid argument '[aa]' to create resource."}

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
	if !strings.Contains(cmdResult, expectRe) {
		t.Errorf("expect result: %s, but got: %s", expectRe, cmdResult)
	}

}
