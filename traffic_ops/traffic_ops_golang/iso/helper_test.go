package iso

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// mockISOCmd returns a modified version of the given Cmd
// so that when run, the command actually invokes the
// TestHelperMockCmd test. See TestHelperMockCmd for
// more details on its behavior.
// If forceError is true, the command will exit a non-0 code and write
// nothing to STDOUT/output.
// If cmdOutput is blank, the command will write to STDOUT, otherwise
// it will write its output to the file specified by cmdOutput.
func mockISOCmd(cmd *exec.Cmd, forceError bool, cmdOutput string) *exec.Cmd {
	args := []string{
		"-test.run=TestHelperMockCmd",
		"--",
	}
	args = append(args, cmd.Args...)

	// os.Args[0] is the invokation of this test binary
	mocked := exec.Command(os.Args[0], args...)

	env := cmd.Env
	env = append(cmd.Env, "GO_HELPER_CMD=1")
	if forceError {
		env = append(env, "GO_HELPER_CMD_FORCE_ERROR=1")
	}
	if cmdOutput != "" {
		env = append(env, fmt.Sprintf("GO_HELPER_CMD_OUTPUT=%s", cmdOutput))
	}
	mocked.Env = env

	return mocked
}

func TestHelperMockCmd(t *testing.T) {
	if os.Getenv("GO_HELPER_CMD") != "1" {
		return
	}

	var respCode int
	if os.Getenv("GO_HELPER_CMD_FORCE_ERROR") == "1" {
		respCode = 1
	}

	dest := os.Stdout
	if cmdOutput := os.Getenv("GO_HELPER_CMD_OUTPUT"); cmdOutput != "" {
		fd, err := os.Create(cmdOutput)
		if err == nil {
			defer fd.Close()
			dest = fd
		}
	}

	// Set args to all arguments past '--'.
	var args []string
	for i, v := range os.Args {
		if v == "--" {
			args = os.Args[i+1:]
			break
		}
	}

	if respCode == 0 {
		fmt.Fprintf(dest, strings.Join(args, " "))
	}
	os.Exit(respCode)
}
