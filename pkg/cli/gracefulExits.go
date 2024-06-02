package cli

import (
	"os"
	"os/exec"
)

func ResetCLI() {
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()

}

func ResetCLI_Exit() int {
	ResetCLI()
	return 1
}
