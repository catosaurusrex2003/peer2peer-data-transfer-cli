package main

import (
	"fmt"
	"log"

	"github.com/c-bata/go-prompt"
	"github.com/manifoldco/promptui"

	"main.go/pkg/cli"
	"main.go/pkg/fileio"
)

func main() {

	// prevent the terminal from fucking up when go-prompt exits
	defer cli.HandleExitCLI()

	menu := promptui.Select{
		Label: "Select Action",
		Items: []string{"Send", "Receive"},
	}

	_, result, err := menu.Run()

	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	// defer term.Restore(int(os.Stdin.Fd()), oldTermimnalState)

	switch result {

	case "Send":
		fmt.Println("Enter the file path to send (Tab for autocomplete):")
		filePath := prompt.Input("> ", cli.Completer)

		if filePath == "" {
			_ = fmt.Errorf("file path cannot be empty")
			return
		}

		fileio.GetFileProperties(filePath)

	case "Receive":
		fmt.Println("You chose to receive.")
	}
}
