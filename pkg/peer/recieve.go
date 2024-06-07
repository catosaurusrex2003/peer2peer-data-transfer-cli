package peer

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"main.go/pkg/cli"
	"main.go/pkg/fileio"
)

func handleStream_onRecieverSide(nodeStream network.Stream) {
	fmt.Println("Got a new stream!")

	buf := make([]byte, 256)
	n, err := nodeStream.Read(buf)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Received message: %s", string(buf[:n]))

	var genericMap map[string]interface{}
	if err := json.Unmarshal(buf[:n], &genericMap); err != nil {
		cli.LogError("Error Marshalling JSON", err)
	}

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

				cache, err := nodeStream.Write([]byte(payloadJsonString))
				if err != nil {
					log.Println("error while sending yes", err)
				}
				fmt.Println("File message Sent return int is", cache)
			}
		}
	} else {
		// This Means that action is not present in json.
		// most prob means the new stream is of file data
		fmt.Println("Unknown JSON structure")
	}

}

func HandleRecieve() {
	// ctx := context.Background()

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
	node.SetStreamHandler("/p2p-event/1.0.0", handleStream_onRecieverSide)

	// Keep the host running
	select {}
}
