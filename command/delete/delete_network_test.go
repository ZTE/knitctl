package delete

import (
	"bytes"
	"knitter.io/knitter/command/create"
	"strings"
	"testing"
	"time"
)

func TestDeletenw_success(t *testing.T) {
	// first create -f test-create.yaml
	bufCreate := bytes.NewBuffer([]byte{})
	cmdCreate := create.NewCmdCreateNW(bufCreate, bufCreate)
	cmdCreate.Flags().Set("tenant", "default")
	cmdCreate.Flags().Set("gateway", "100.100.10.1")
	cmdCreate.Flags().Set("cidr", "100.100.10.0/24")
	cmdCreate.Run(cmdCreate, []string{"delete_cmd_network_test"})

	time.Sleep(time.Second * 1)

	buf := bytes.NewBuffer([]byte{})
	cmd := NewCmdDeleteNW(buf, buf)
	cmd.Flags().Set("tenant", "default")
	cmd.Run(cmd, []string{"delete_cmd_network_test"})

	cmdResult := buf.String()
	expectResult := []string{"delete network 'delete_cmd_network_test' in tenant 'default' success"}

	for _, oneExpect := range expectResult {
		if !strings.Contains(cmdResult, oneExpect) {
			t.Errorf("expect result: %s, but got: %s", oneExpect, cmdResult)
		}
	}

}

func TestDeletenw_failure(t *testing.T) {

	buf := bytes.NewBuffer([]byte{})
	cmd := NewCmdDeleteNW(buf, buf)
	cmd.Flags().Set("tenant", "default")
	cmd.Run(cmd, []string{})
	cmdResult := buf.String()
	expectResult := []string{"no set network name to delete."}

	for _, oneExpect := range expectResult {
		if !strings.Contains(cmdResult, oneExpect) {
			t.Errorf("expect result: %s, but got: %s", oneExpect, cmdResult)
		}
	}

}
