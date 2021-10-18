package p2p

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/wallet"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"go.uber.org/zap"
	"io"
	mrand "math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/persistence"
	net "github.com/libp2p/go-libp2p-core"
	host "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pstore "github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/protocol"
	ma "github.com/multiformats/go-multiaddr"
	"strings"
)

var mutex = &sync.Mutex{}

type EndPoint struct {
	cfg         config.P2P
	persistence persistence.IPersistence
	host        host.Host
	wallets     wallet.IWallets
	log         *zap.Logger
	event       *event.Event
	enabled     bool
	msgReceived []string
}

type message struct {
	Name  event.EventType
	ID    string
	Value []byte
}

func Init(
	cfg config.P2P,
	per persistence.IPersistence,
	wallets wallet.IWallets,
	logs *zap.Logger,
	evt *event.Event,
) *EndPoint {

	e := &EndPoint{
		cfg:         cfg,
		persistence: per,
		wallets:     wallets,
		log:         logs.With(zap.String("service", "p2p")),
		event:       evt,
		enabled:     cfg.Enabled,
	}

	return e
}

func (e *EndPoint) Enabled() bool {
	return e.enabled
}

func (e *EndPoint) Listen(stop chan error) {
	e.hasRequiredPort()
	go func() {
		if err := e.makeBasicHost(); err != nil {
			e.log.Error("fail to listen p2p", zap.Error(err))
			stop <- err
			return
		}

		if !e.HasTarget() {
			e.log.Info("listening for connections")
			// Set a stream handler on host A. /p2p/1.0.0 is
			// a user-defined protocol name.
			e.host.SetStreamHandler("/p2p/1.0.0", e.handleStream)
			// select {} // hang forever
			/**** This is where the listener code ends ****/
		} else {
			// e.log.Fatal("fail, no peer address target found")
			e.host.SetStreamHandler("/p2p/1.0.0", e.handleStream)
			e.connectToIPFS(stop, e.host)
		}

		sigCh := make(chan os.Signal)
		signal.Notify(sigCh, syscall.SIGKILL, syscall.SIGINT)
		<-sigCh
	}()
}

func (e *EndPoint) hasRequiredPort() {
	if e.cfg.Port == 0 {
		e.log.Fatal("Please provide a port to bind on with -l")
	}
}

func (e *EndPoint) HasTarget() bool {
	if e.cfg.Target == "" {
		// call default genesis
		return false
	}
	return true
}

func (e *EndPoint) connectToIPFS(stop chan error, ha host.Host) {
	// The following code extracts target's peer ID from the
	// given multiaddress
	ipfsAddr, err := ma.NewMultiaddr(e.cfg.Target)
	if err != nil {
		e.log.Error("fail to set new multi address", zap.Error(err), zap.String("target", e.cfg.Target))
		stop <- err
		return
	}

	/*
		peerAddrInfo, err := peer.AddrInfoFromP2pAddr(ipfsAddr)
		if err != nil {
			panic(err)
		}

		// Connect to the node at the given address.
		if err := host.Connect(context.Background(), *peerAddrInfo); err != nil {
			panic(err)
		}*/

	pid, err := ipfsAddr.ValueForProtocol(ma.P_IPFS)
	if err != nil {
		e.log.Error("fail to get protocol to ipfs", zap.Error(err))
		stop <- err
		return
	}
	// Nous nous retrouvons avec le peerID et l'adresse cible targetAddr de l'hôte auquel nous voulons nous connecter
	// et ajoutons cet enregistrement dans notre "magasin"
	// afin que nous puissions garder une trace de qui nous sommes connectés.
	// Nous le faisons avec ha.Peerstore().AddAddr
	peerid, err := peer.Decode(pid)
	if err != nil {
		e.log.Error("fail to get decode peer", zap.Error(err))
		stop <- err
		return
	}
	// Decapsulate the /ipfs/<peerID> part from the target
	// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
	targetPeerAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", peer.Encode(peerid)))
	if err != nil {
		e.log.Error("fail to set new multi address", zap.Error(err), zap.String("target", e.cfg.Target))
		stop <- err
		return
	}

	targetAddr := ipfsAddr.Decapsulate(targetPeerAddr)

	// We have a peer ID and a targetAddr so we add it to the peerstore
	// so LibP2P knows how to contact it
	ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

	e.log.Info("opening stream p2p", zap.String("target", e.cfg.Target))
	// make a new stream from host B to host A
	// it should be handled on host A by the handler we set above because
	// we use the same /p2p/1.0.0 protocol
	protocolID := protocol.ID(e.cfg.ProtocolID)
	s, err := ha.NewStream(context.Background(), peerid, protocolID)
	if err != nil {
		e.log.Error("fail to set new stream", zap.Error(err), zap.Any("peer_id", peerid), zap.Any("protocol_id", protocolID))
		stop <- err
		return
	}
	// Create a buffered stream so that read and writes are non blocking.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	// Create a thread to read and write data.
	go e.writeData(rw)
	go e.readData(rw)
}

func (e *EndPoint) makeBasicHost() error {
	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var (
		r   io.Reader = e.setIoReader()
		err error
	)

	// Generate a key pair for this host. We will use it
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", e.cfg.Port)),
		libp2p.Identity(priv),
	}

	e.host, err = libp2p.New(context.Background(), opts...)
	if err != nil {
		return err
	}

	e.log.Info("P2P start:", zap.Any("address", e.host.Addrs()), zap.Any("host_id", e.host.ID()))
	e.setStreamHandler()

	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return err
	}

	// Parse the multiaddr string.
	peerMA, err := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))
	//peerMA, err := ma.NewMultiaddr(e.cfg.Target)
	if err != nil {
		return err
	}

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addrs := basicHost.Addrs()
	var addr ma.Multiaddr
	// select the address starting with "ip4"
	for _, i := range addrs {
		if strings.HasPrefix(i.String(), "/ip4") {
			addr = i
			break
		}
	}

	fullAddr := addr.Encapsulate(peerMA)
	e.log.Info(fmt.Sprintf("I am %s\n", fullAddr))
	if e.cfg.Secio {
		e.log.Info(fmt.Sprintf("Now run \"go run main.go -api_port %d -p2p_port %d -p2p_target %s -secio\" on a different terminal", e.cfg.Port+15, e.cfg.Port+16, fullAddr))
	} else {
		e.log.Info(fmt.Sprintf("Now run \"go run main.go -api_port %d -p2p_port %d -p2p_target %s\" on a different terminal", e.cfg.Port+15, e.cfg.Port+16, fullAddr))
	}

	return nil
}

// Setup a stream handler.
//
// This gets called every time a peer connects and opens a stream to this node.
func (e *EndPoint) setStreamHandler() {
	protocolID := protocol.ID(e.cfg.ProtocolID)
	e.host.SetStreamHandler(protocolID, func(s network.Stream) {
		go e.writeCounter(s)
		go e.readCounter(s)
	})
}

func (e *EndPoint) setIoReader() io.Reader {
	var r io.Reader
	if e.cfg.Seed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(e.cfg.Seed))
	}
	return r
}

func (e *EndPoint) handleStream(s net.Stream) {

	e.log.Info("Got a new stream p2p")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go e.readData(rw)
	go e.writeData(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}

func (e *EndPoint) Shutdown() {
	e.host.Close()
}
