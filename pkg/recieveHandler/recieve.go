package recieveHandler

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	net "github.com/libp2p/go-libp2p/core/network"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
)

func handleStream(s net.Stream) {
	fmt.Println("Got a new stream!")

	buf := make([]byte, 256)
	for {
		n, err := s.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Received:", string(buf[:n]))
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

	// Print the host's Multiaddress
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Your libp2p node address:", addrs[0])

	// Set a stream handler on the host
	node.SetStreamHandler("/p2p-event/1.0.0", handleStream)

	// Keep the host running
	select {}
}
