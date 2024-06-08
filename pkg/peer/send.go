package peer

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
	"main.go/pkg/cli"
	"main.go/pkg/fileio"
)

// channel to handle the goroutine for writing data un the butffer stream
var writeChannel_Sender = make(chan string, 10)

func handleSendFileInfo(filePath string) {
	fileInfo := fileio.GetFileProperties(filePath)

	fileInfoMessage := FileMetadata_MessageType{
		Action:       "fileInfoSend",
		FileName:     fileInfo.Name(),
		Size:         fileInfo.Size(),
		IsDir:        fileInfo.IsDir(),
		LastModified: fileInfo.ModTime().Format(time.RFC1123),
	}

	jsonString, err := json.Marshal(fileInfoMessage)
	if err != nil {
		fmt.Println("Error marshalling fileInfoMessage:", err)
		return
	}

	writeChannel_Sender <- string(jsonString)
}

func readStream_Sender(rw *bufio.ReadWriter) {
	defer wg_Sender.Done()
	for {
		fmt.Println("reading")
		str, err := rw.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from stream:", err)
			}
			return
		}

		var msg map[string]interface{}
		err = json.Unmarshal([]byte(str), &msg)
		if err != nil {
			log.Println("Error unmarshaling JSON:", err)
			continue
		}

		fmt.Println("\x1b[32m", msg, "\x1b[0m>")
	}
}

func writeStream_Sender(rw *bufio.ReadWriter) {
	defer wg_Sender.Done()
	for {
		for msg := range writeChannel_Sender {
			// the msg is a json string
			// jsonData, err := json.Marshal(msg)
			// if err != nil {
			// 	log.Println("Error marshaling JSON:", err)
			// 	continue
			// }

			fmt.Println("writting : ", msg)

			_, err := rw.WriteString(fmt.Sprintf("%s\n", msg))
			if err != nil {
				log.Println("Error writing to stream:", err)
				return
			}
			rw.Flush()
		}

	}

}

var wg_Sender sync.WaitGroup

func SenderMain(filePath string) {
	//
	// Main function of a sender which
	// Handles the p2p connection and coordinates the event messages
	// uses 2 gorotines. one to read and one to write in the network.Stream
	//

	ctx := context.Background()
	wg_Sender.Add(2)

	// Assigning identity to this Host
	priv, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	node, err := libp2p.New(
		libp2p.Identity(priv),
	)
	if err != nil {
		log.Fatal(err)
	}

	peerInfo := peerstore.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}

	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		cli.LogError("%w", err)
		cli.ResetCLI()
	}
	fmt.Println("")
	fmt.Println("Your libp2p node address:", addrs[0])
	fmt.Println("")

	// get recievers address
	fmt.Print("> Enter Recievers Address: ")
	var targetAddrStr string
	_, err2 := fmt.Scanln(&targetAddrStr)
	if err2 != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Extract the target peer ID from the multiaddress
	targetAddr, err := multiaddr.NewMultiaddr(targetAddrStr)
	if err != nil {
		log.Fatal("3", err)
	}
	peerinfo, err := peerstore.AddrInfoFromP2pAddr(targetAddr)
	if err != nil {
		log.Fatal("4", err)
	}

	// Connect to the target peer
	if err := node.Connect(ctx, *peerinfo); err != nil {
		log.Fatal("5", err)
	}
	fmt.Println("Connected .....")

	// Create a stream to the target peer
	s, err := node.NewStream(ctx, peerinfo.ID, "/p2p-event/1.0.0")
	if err != nil {
		log.Fatal("6", err)
	}

	log.Println("Established connection to destination")

	// Create a buffered stream so that read and writes are non-blocking.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	// Send the initial file Info Message
	// handlSendFileInfo(nodeStream, filePath)

	// HERE SETUP A GOROUTINE
	go writeStream_Sender(rw)
	go readStream_Sender(rw)

	handleSendFileInfo(filePath)
	wg_Sender.Wait()
}

// func handleStream_onSenderSide(s network.Stream) {

// 	buf := make([]byte, 256)
// 	n, err := s.Read(buf)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	log.Printf("Received message: %s", string(buf[:n]))

// 	var genericMap map[string]interface{}
// 	if err := json.Unmarshal(buf[:n], &genericMap); err != nil {
// 		cli.LogError("Error Marshalling JSON", err)
// 	}

// 	actionValue, ok := genericMap["action"]
// 	// Determine type of message based on action
// 	if ok {
// 		if actionValue == "fileAcceptMessage" {
// 			if genericMap["value"] == "accept" {
// 				fmt.Println("SENDING FILE")
// 			} else {
// 				fmt.Println("Ah Shit i cant send the file")
// 			}
// 		}
// 	} else {
// 		// This Means that action is not present in json.
// 		// most prob means the new stream is of file data
// 		fmt.Println("Unknown JSON structure")
// 	}

// }
