package peer

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
	"main.go/pkg/cli"
	"main.go/pkg/fileio"
)

func handlSendFileInfo(nodeStream network.Stream, filePath string) {

	fileInfo := fileio.GetFileProperties(filePath)

	fileInfoMessage := FileMetadata_MessageType{
		Action:       "fileInfoSend",
		FileName:     fileInfo.Name(),
		Size:         fileInfo.Size(),
		IsDir:        fileInfo.IsDir(),
		LastModified: fileInfo.ModTime().Format(time.RFC1123)}
	jsonString, err := json.Marshal(fileInfoMessage)
	if err != nil {
		fmt.Println("error marshalling fileInfoMessage")
	}

	cache, err := nodeStream.Write([]byte(jsonString))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File message Sent return int is", cache)
}

func handleStream_onSenderSide(s network.Stream) {

	buf := make([]byte, 256)
	n, err := s.Read(buf)
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
		if actionValue == "fileAcceptMessage" {
			if genericMap["value"] == "accept" {
				fmt.Println("SENDING FILE")
			} else {
				fmt.Println("Ah Shit i cant send the file")
			}
		}
	} else {
		// This Means that action is not present in json.
		// most prob means the new stream is of file data
		fmt.Println("Unknown JSON structure")
	}

}

func HandleSend(filePath string) {
	ctx := context.Background()

	priv, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		log.Fatal("1", err)
	}
	node, err := libp2p.New(
		libp2p.Identity(priv),
	)
	if err != nil {
		log.Fatal("2", err)
	}
	node.SetStreamHandler("/p2p-event/1.0.0", handleStream_onSenderSide)

	peerInfo := peerstore.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}

	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		cli.LogError("%w", err)
		cli.ResetCLI()
	}
	fmt.Println("Your libp2p node address:", addrs[0])

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
	nodeStream, err := node.NewStream(ctx, peerinfo.ID, "/p2p-event/1.0.0")
	if err != nil {
		log.Fatal("6", err)
	}

	// Send the initial file Info Message
	handlSendFileInfo(nodeStream, filePath)
}
