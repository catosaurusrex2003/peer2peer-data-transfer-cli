package peer

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	network "github.com/libp2p/go-libp2p/core/network"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
	"main.go/pkg/cli"
	"main.go/pkg/fileio"
)

func handlSendFileInfo(nodeStream network.Stream, filePath string) {

	fileInfo := fileio.GetFileProperties(filePath)

	fileInfoMessage := FileMetadata_MessageType{Action: "fileInfoSend", FileName: fileInfo.Name(), Size: fileInfo.Size()}
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

// func handleSendStream(s net.Stream) {
// 	fmt.Println("Got a new stream!")

// 	buf := make([]byte, 256)
// 	for {
// 		n, err := s.Read(buf)
// 		if err == io.EOF {
// 			break
// 		} else if err != nil {
// 			log.Fatal(err)
// 		}

// 		fmt.Println("Received:", string(buf[:n]))
// 	}
// }

func HandleSend(filePath string) {
	ctx := context.Background()

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
	fmt.Println("Your libp2p node address:", addrs[0])

	fmt.Print("> Enter Recievers Address: ")

	var targetAddrStr string
	_, err2 := fmt.Scanln(&targetAddrStr)
	if err2 != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	targetAddr, err := multiaddr.NewMultiaddr(targetAddrStr)
	if err != nil {
		log.Fatal(err)
	}

	// Extract the target peer ID from the multiaddress
	peerinfo, err := peerstore.AddrInfoFromP2pAddr(targetAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the target peer
	if err := node.Connect(ctx, *peerinfo); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to", targetAddrStr)

	// Create a stream to the target peer
	nodeStream, err := node.NewStream(ctx, peerinfo.ID, "/p2p-event/1.0.0")
	if err != nil {
		log.Fatal(err)
	}

	// Send the initial file Info Message
	handlSendFileInfo(nodeStream, filePath)

	// Send events (messages) to the target peer
	// for i := 0; i < 5; i++ {
	// 	msg := fmt.Sprintf("Event #%d", i)
	// 	_, err := s.Write([]byte(msg))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println("Sent:", msg)
	// 	time.Sleep(2 * time.Second)
	// }
}
