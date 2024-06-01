package main

import (
	"fmt"
	"log"

	"github.com/c-bata/go-prompt"
	"github.com/manifoldco/promptui"

	"main.go/utils"
)

func main() {

	// oldTermimnalState, err := term.MakeRaw(int(os.Stdin.Fd()))
	// if err != nil {
	// 	fmt.Println("%w", err)
	// }

	// spew.Dump(oldTermimnalState)

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
		fmt.Println("You chose to send.")
		fmt.Println("Enter the file path to send (Tab for autocomplete):")
		filePath := prompt.Input("> ", utils.Completer)

		if filePath == "" {
			log.Fatalf("File path cannot be empty")
		}

		fmt.Printf("You chose to send the file: %s\n", filePath)

		// Add your send logic here with the filePath

	case "Receive":
		fmt.Println("You chose to receive.")
		// Add your receive logic here
	}
}
