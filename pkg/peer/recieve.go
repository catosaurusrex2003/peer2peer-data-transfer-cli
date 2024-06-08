package peer

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"main.go/pkg/fileio"
)

var writeChannel_Reciever = make(chan string, 10)

func readStream_Reciever(rw *bufio.ReadWriter) {
	defer wg_Reciever.Done()
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from stream:", err)
			}
			return
		}

		var genericMap map[string]interface{}
		err = json.Unmarshal([]byte(str), &genericMap)
		if err != nil {
			log.Println("Error unmarshaling JSON:", err)
			continue
		}
		fmt.Printf("\x1b[32m%s\x1b[0m> ", genericMap)

		actionValue, ok := genericMap["action"]
		// Determine type of message based on action
		if ok {
			if actionValue == "fileInfoSend" {

				fileInfoData := fileio.FileInfoData{
					FileName:         genericMap["fileName"].(string),
					FileSize:         int64(genericMap["size"].(float64)),
					FileLastModified: genericMap["lastModified"].(string),
					FileIsDir:        genericMap["isDir"].(bool),
				}

				fmt.Println("\n<<<<<<< Incoming File Info >>>>>>>")
				fileio.LogFileInfo(fileInfoData)

				fmt.Print("> Accept File ? (y/n) : ")

				var acceptChoice string
				_, err2 := fmt.Scanln(&acceptChoice)
				if err2 != nil {
					fmt.Println("Error reading input:", err)
					return
				}
				// fmt.Println("acceptChoice is ", strings.ToLower(acceptChoice))

				if strings.ToLower(acceptChoice) == "y" {
					fmt.Println("recieving yes pls wait")
					payload := FileAccept_MessageType{
						Action: "fileAcceptMessage",
						Value:  "accept",
					}
					payloadJsonString, err := json.Marshal(payload)
					fmt.Println(payloadJsonString)
					if err != nil {
						fmt.Println("error marshalling fileInfoMessage")
					}

					writeChannel_Reciever <- string(payloadJsonString)
				}
			}
		} else {
			// This Means that action is not present in json.
			// most prob means the new stream is of file data
			fmt.Println("Unknown JSON structure")
		}
	}
}

func writeStream_Reciever(rw *bufio.ReadWriter) {
	defer wg_Reciever.Done()
	for {
		for msg := range writeChannel_Reciever {
			// the msg is a json string

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

var wg_Reciever sync.WaitGroup

func RecieverMain() {
	// ctx := context.Background()
	wg_Reciever.Add(2)

	// Generate a key pair for this host
	priv, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new libp2p Host
	node, err := libp2p.New(
		libp2p.Identity(priv),
		// libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/10000"),
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
		panic(err)
	}

	fmt.Println("Your libp2p node address:", addrs[0])

	// Set a stream handler on the host
	// node.SetStreamHandler("/p2p-event/1.0.0", handleStream_onRecieverSide)
	node.SetStreamHandler("/p2p-event/1.0.0", func(s network.Stream) {
		log.Println("Established connection to destination")
		// Create a buffered stream so that read and writes are non-blocking.
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		go writeStream_Reciever(rw)
		go readStream_Reciever(rw)
	})

	wg_Reciever.Wait()
}
