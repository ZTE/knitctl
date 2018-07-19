package create

import (
	"bytes"
	"knitter.io/knitter/command/delete"
	"strings"
	"testing"
	"time"
)

func TestCreateipg_success(t *testing.T) {

	buf := bytes.NewBuffer([]byte{})
	cmd := NewCmdCreateIpgroup(buf, buf)
	cmd.Flags().Set("tenant", "default")
	cmd.Flags().Set("network", "lan")
	cmd.Flags().Set("ips", "100.100.10.3")
	cmd.Run(cmd, []string{"create_cmd_ipgroup_test"})
	cmdResult := buf.String()
	expectResult := []string{"create ipgroup 'create_cmd_ipgroup_test' in tenant 'default' success"}
	for _, oneExpect := range expectResult {
		if !strings.Contains(cmdResult, oneExpect) {
			t.Errorf("expect result: %s, but got: %s", oneExpect, cmdResult)
		}
	}

	time.Sleep(time.Second * 1)

	cmdDelete := delete.NewCmdDeleteIPG(buf, buf)
	cmdDelete.Flags().Set("tenant", "default")
	cmdDelete.Run(cmdDelete, []string{"create_cmd_ipgroup_test"})
}

func TestCreateipg_failure(t *testing.T) {

	buf := bytes.NewBuffer([]byte{})
	cmd := NewCmdCreateIpgroup(buf, buf)
	cmd.Flags().Set("tenant", "default")
	cmd.Flags().Set("network", "lan")
	cmd.Flags().Set("ips", "100.100.10.3")

	cmd.Run(cmd, []string{})

	cmdResult := buf.String()
	expectResult := []string{"no set ipgroup name to create"}

	for _, oneExpect := range expectResult {
		if !strings.Contains(cmdResult, oneExpect) {
			t.Errorf("expect result: %s, but got: %s", oneExpect, cmdResult)
		}
	}
}
