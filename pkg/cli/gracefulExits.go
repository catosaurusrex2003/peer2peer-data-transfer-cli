package cli

import (
	"os"
	"os/exec"
)

func ResetCLI() {
	// reset the terminal to the original state
	// used when the go-prompt fucks up our terminal
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()

}

func ResetCLI_Exit() int {
	ResetCLI()
	return 1
}
