package cli

import (
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"
)

func GetFilePathPrompt() string {
	//
	// Prompt to choose file using autocomplete
	//

	// prevent go-prompt for fucking the terminal
	defer ResetCLI()

	fmt.Println("Enter the file path to send (Tab for autocomplete):")
	filePath := prompt.Input("> ", Completer)

	if filePath == "" {
		_ = fmt.Errorf("file path cannot be empty")
		os.Exit(ResetCLI_Exit())
	}
	return filePath
}
