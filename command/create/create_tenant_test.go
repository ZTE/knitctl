package create

import (
	"bytes"
	"testing"
	"strings"
	"github.com/ZTE/knitctl/command/delete"
	"time"
)

func TestCreateTenant_success(t *testing.T) {

	buf := bytes.NewBuffer([]byte{})
	cmd := NewCmdCreateTenant(buf, buf)
	cmd.Run(cmd, []string{"create_cmd_tenant_test"})
	cmdResult := buf.String()
	expectResult := []string{"create tenant 'create_cmd_tenant_test' success"}

	for _, oneExpect := range expectResult {
		if !strings.Contains(cmdResult, oneExpect) {
			t.Errorf("expect result: %s, but got: %s", oneExpect, cmdResult)
		}
	}

	time.Sleep(time.Second *1)

	cmdDelete := delete.NewCmdDeleteTenant(buf, buf)
	cmdDelete.Run(cmdDelete, []string{"create_cmd_tenant_test"})
}

func TestCreateTenant_failure(t *testing.T) {

	buf := bytes.NewBuffer([]byte{})
	cmd := NewCmdCreateTenant(buf, buf)

	cmd.Run(cmd, []string{})

	cmdResult := buf.String()
	expectResult := []string{"no set tenant name to create"}

	for _, oneExpect := range expectResult {
		if !strings.Contains(cmdResult, oneExpect) {
			t.Errorf("expect result: %s, but got: %s", oneExpect, cmdResult)
		}
	}
}
