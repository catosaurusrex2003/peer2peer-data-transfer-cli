package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/tabwriter"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/manifoldco/promptui"

	"main.go/utils"
)

func handleExitCLI() int {
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
	return 1
}

func getFileProperties(filePath string) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatalf("Failed to get file properties: %v", err)
	}

	if fileInfo.IsDir() {
		utils.LogError("It is a directory. You can compress it and send as a file. Right now compression is not supported in this program")
		os.Exit(handleExitCLI())
	}

	// Print the File Info.
	// NEED: colour change. make more good looking
	fmt.Println("\n<<<<<<<File Info>>>>>>>")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "File Name\t%s\n", fileInfo.Name())
	fmt.Fprintf(w, "Size\t%d bytes\n", fileInfo.Size())
	fmt.Fprintf(w, "Permissions\t%s\n", fileInfo.Mode().String())
	fmt.Fprintf(w, "Last Modified\t%s\n", fileInfo.ModTime().Format(time.RFC1123))
	fmt.Fprintf(w, "Is Directory\t%t\n", fileInfo.IsDir())
	w.Flush()
}

func main() {

	// oldTermimnalState, err := term.MakeRaw(int(os.Stdin.Fd()))
	// if err != nil {
	// 	fmt.Println("%w", err)
	// }

	// spew.Dump(oldTermimnalState)

	// prevent the terminal from fucking up when go-prompt exits
	defer handleExitCLI()

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
		filePath := prompt.Input("> ", utils.Completer)

		if filePath == "" {
			_ = fmt.Errorf("file path cannot be empty")
			return
		}

		getFileProperties(filePath)

		// Add your send logic here with the filePath

	case "Receive":
		fmt.Println("You chose to receive.")
		// Add your receive logic here
	}
}
