package create

import (
	"bytes"
	"testing"
	"strings"
	"github.com/ZTE/knitctl/command/delete"
	"time"
)

func TestCreatenw_success(t *testing.T) {

	buf := bytes.NewBuffer([]byte{})
	cmd := NewCmdCreateNW(buf, buf)
	cmd.Flags().Set("tenant", "default")
	cmd.Flags().Set("gateway", "10.144.220.1")
	cmd.Flags().Set("cidr", "10.144.220.0/24")
	cmd.Run(cmd, []string{"create_cmd_network_test"})

	cmdResult := buf.String()
	expectResult := []string{"create network 'create_cmd_network_test' in tenant 'default' success"}

	for _, oneExpect := range expectResult {
		if !strings.Contains(cmdResult, oneExpect) {
			t.Errorf("expect result: %s, but got: %s", oneExpect, cmdResult)
		}
	}

	time.Sleep(time.Second *1)

	cmdDelete := delete.NewCmdDeleteNW(buf, buf)
	cmdDelete.Flags().Set("tenant", "default")
	cmdDelete.Run(cmdDelete, []string{"create_cmd_network_test"})
}

func TestCreatenw_failure(t *testing.T) {

	buf := bytes.NewBuffer([]byte{})
	cmd := NewCmdCreateNW(buf, buf)
	cmd.Flags().Set("tenant", "default")
	cmd.Flags().Set("gateway", "10.144.220.1")
	cmd.Flags().Set("cidr", "10.144.220.0/24")

	cmd.Run(cmd, []string{})

	cmdResult := buf.String()
	expectResult := []string{"no set network name to create"}

	for _, oneExpect := range expectResult {
		if !strings.Contains(cmdResult, oneExpect) {
			t.Errorf("expect result: %s, but got: %s", oneExpect, cmdResult)
		}
	}
}