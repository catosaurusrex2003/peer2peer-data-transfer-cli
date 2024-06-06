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

type FileInfoData struct {
	FileName         string
	FileSize         int64
	FilePerms        string
	FileLastModified string
	FileIsDir        bool
}

func GetFileProperties(filePath string) fs.FileInfo {
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

func LogFileInfo(info FileInfoData) {

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "File Name\t%s\n", info.FileName)
	fmt.Fprintf(w, "Size\t%d bytes\n", info.FileSize)
	fmt.Fprintf(w, "Permissions\t%s\n", info.FilePerms)
	fmt.Fprintf(w, "Last Modified\t%s\n", info.FileLastModified)
	fmt.Fprintf(w, "Is Directory\t%t\n", info.FileIsDir)
	w.Flush()
}

func GetAndLogFileProperties(filePath string) {
	fileInfo := GetFileProperties(filePath)
	fileInfoData := FileInfoData{
		FileName:         fileInfo.Name(),
		FileSize:         fileInfo.Size(),
		FilePerms:        fileInfo.Mode().String(),
		FileLastModified: fileInfo.ModTime().Format(time.RFC1123),
		FileIsDir:        fileInfo.IsDir(),
	}
	fmt.Println("\n<<<<<<< File Info >>>>>>>")
	LogFileInfo(fileInfoData)

}
