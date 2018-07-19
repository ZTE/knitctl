package delete

import (
	"bytes"
	"knitter.io/knitter/command/create"
	"strings"
	"testing"
	"time"
)

func TestDeleteipg_success(t *testing.T) {
	// first create -f test-create.yaml
	bufCreate := bytes.NewBuffer([]byte{})
	cmdCreateIpg := create.NewCmdCreateIpgroup(bufCreate, bufCreate)
	cmdCreateIpg.Flags().Set("tenant", "default")
	cmdCreateIpg.Flags().Set("ips", "100.100.221.33")
	cmdCreateIpg.Flags().Set("network", "lan")
	cmdCreateIpg.Run(cmdCreateIpg, []string{"delete_cmd_ipgroup_test"})

	time.Sleep(time.Second * 1)

	buf := bytes.NewBuffer([]byte{})
	cmd := NewCmdDeleteIPG(buf, buf)
	cmd.Flags().Set("tenant", "default")
	cmd.Run(cmd, []string{"delete_cmd_ipgroup_test"})

	cmdResult := buf.String()
	expectResult := []string{"delete ipgroup name: 'delete_cmd_ipgroup_test'"}

	for _, oneExpect := range expectResult {
		if !strings.Contains(cmdResult, oneExpect) {
			t.Errorf("expect result: %s, but got: %s", oneExpect, cmdResult)
		}
	}

}

func TestDeleteipg_failure(t *testing.T) {

	buf := bytes.NewBuffer([]byte{})
	cmd := NewCmdDeleteIPG(buf, buf)
	cmd.Flags().Set("tenant", "default")
	cmd.Run(cmd, []string{})
	cmdResult := buf.String()
	expectResult := []string{"no set ipgroup name to delete."}

	for _, oneExpect := range expectResult {
		if !strings.Contains(cmdResult, oneExpect) {
			t.Errorf("expect result: %s, but got: %s", oneExpect, cmdResult)
		}
	}

}
