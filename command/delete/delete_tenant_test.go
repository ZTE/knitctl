package delete

import (
	"bytes"
	"knitter.io/knitter/command/create"
	"strings"
	"testing"
	"time"
)

func TestDeleteTenant_success(t *testing.T) {

	bufCreateTenant := bytes.NewBuffer([]byte{})
	cmdCreateTenant := create.NewCmdCreateTenant(bufCreateTenant, bufCreateTenant)
	cmdCreateTenant.Run(cmdCreateTenant, []string{"delete_cmd_tenant_test"})

	time.Sleep(time.Second * 1)

	buf := bytes.NewBuffer([]byte{})
	cmd := NewCmdDeleteTenant(buf, buf)
	cmd.Run(cmd, []string{"delete_cmd_tenant_test"})

	cmdResult := buf.String()

	expectResult := []string{"delete tenant 'delete_cmd_tenant_test' success"}

	for _, oneExpect := range expectResult {
		if !strings.Contains(cmdResult, oneExpect) {
			t.Errorf("expect result: %s, but got: %s", oneExpect, cmdResult)
		}
	}

}

func TestDeleteTenant_failure(t *testing.T) {

	buf := bytes.NewBuffer([]byte{})
	cmd := NewCmdDeleteTenant(buf, buf)
	cmd.Run(cmd, []string{})
	cmdResult := buf.String()
	expectResult := []string{"no set tenant name to delete."}

	for _, oneExpect := range expectResult {
		if !strings.Contains(cmdResult, oneExpect) {
			t.Errorf("expect result: %s, but got: %s", oneExpect, cmdResult)
		}
	}

}
