package cli

import (
	"os"
	"os/exec"
)

func HandleExitCLI() int {
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
	return 1
}
