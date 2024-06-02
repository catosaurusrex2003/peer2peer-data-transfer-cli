package sendHandler

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
)

func HandleSend() {
	ctx := context.Background()

	// Generate a key pair for this host
	priv, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new libp2p Host
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

	// Print the host's Multiaddress
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("libp2p node address:", addrs[0])

	// Define the target peer multiaddress (Replace with your target peer's multiaddress)
	targetAddrStr := "/ip4/127.0.0.1/tcp/10000/p2p/12D3KooWEPPE5ZMZCDkS9eAYrVDXuiFief1iwZjVbvny3rm7syqE"
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
	s, err := node.NewStream(ctx, peerinfo.ID, "/p2p-event/1.0.0")
	if err != nil {
		log.Fatal(err)
	}

	// Send events (messages) to the target peer
	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("Event #%d", i)
		_, err := s.Write([]byte(msg))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Sent:", msg)
		time.Sleep(2 * time.Second)
	}
}
