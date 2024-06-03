package fileio

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"main.go/pkg/cli"
)

func getFileProperties(filePath string) fs.FileInfo {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatalf("Failed to get file properties: %v", err)
		return nil
	}

	if fileInfo.IsDir() {
		cli.LogError("It is a directory. You can compress it and send as a file. Right now compression is not supported in this program")
		os.Exit(cli.ResetCLI_Exit())
	}
	return fileInfo
}

func GetAndLogFileProperties(filePath string) {
	fileInfo := getFileProperties(filePath)
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
