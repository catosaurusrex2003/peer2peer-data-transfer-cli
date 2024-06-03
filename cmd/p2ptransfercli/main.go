package main

import (
	"fmt"
	"log"

	"github.com/manifoldco/promptui"

	"main.go/pkg/cli"
	"main.go/pkg/fileio"
	"main.go/pkg/peer"
)

func main() {

	// PREVENT THE TERMINAL FROM FUCKIGN UP WHEN KERNEL PANICS
	defer cli.ResetCLI()

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
		filePath := cli.GetFilePathPrompt()
		fileio.GetAndLogFileProperties(filePath)
		peer.HandleSend()

	case "Receive":
		fmt.Println("You chose to receive.")
		peer.HandleRecieve()
	}
}
