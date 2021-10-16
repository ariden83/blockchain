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
	"strings"
	"sync"

	"github.com/ariden83/blockchain/internal/event"
	net "github.com/libp2p/go-libp2p-core"
	host "github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pstore "github.com/libp2p/go-libp2p-core/peerstore"
	ma "github.com/multiformats/go-multiaddr"
)

var mutex = &sync.Mutex{}

type EndPoint struct {
	cfg         config.P2P
	persistence iPersistence
	host        host.Host
	wallets     iWallets
	log         *zap.Logger
	event       *event.Event
	enabled     bool
}

type iPersistence interface {
	GetLastHash() ([]byte, error)
	Update([]byte, []byte) error
}

type iWallets interface {
	GetSeeds() *[]wallet.Seed
	GetAllPublicSeeds() []wallet.SeedNoPrivKey
}

func Init(
	cfg config.P2P,
	per iPersistence,
	wallets iWallets,
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
	go func() {
		ha, err := e.makeBasicHost()
		if err != nil {
			e.log.Error("fail to listen p2p", zap.Error(err))
			stop <- err
			return
		}

		if e.cfg.Target == "" {
			e.log.Info("listening for connections")
			// Set a stream handler on host A. /p2p/1.0.0 is
			// a user-defined protocol name.
			ha.SetStreamHandler("/p2p/1.0.0", e.handleStream)

			select {} // hang forever
			/**** This is where the listener code ends ****/
		} else {
			ha.SetStreamHandler("/p2p/1.0.0", e.handleStream)
			e.connectToIPFS(stop, ha)
		}
	}()
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
		e.log.Error("fail to set new multi address to ipfs", zap.Error(err), zap.String("target", e.cfg.Target))
		stop <- err
		return
	}

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
	targetPeerAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", peer.Encode(peerid)))
	targetAddr := ipfsAddr.Decapsulate(targetPeerAddr)

	// We have a peer ID and a targetAddr so we add it to the peerstore
	// so LibP2P knows how to contact it
	ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

	e.log.Info("opening stream p2p", zap.String("target", e.cfg.Target))
	// make a new stream from host B to host A
	// it should be handled on host A by the handler we set above because
	// we use the same /p2p/1.0.0 protocol
	s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
	if err != nil {
		e.log.Error("fail to set new stream", zap.Error(err))
		stop <- err
		return
	}
	// Create a buffered stream so that read and writes are non blocking.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	// Create a thread to read and write data.
	go e.writeData(rw)
	go e.readData(rw)

	// select vide afin que notre programme ne se contente pas de se terminer et de s'arrêter
	select {} // hang forever
}

func (e *EndPoint) makeBasicHost() (host.Host, error) {

	if e.cfg.Port == 0 {
		e.log.Fatal("Please provide a port to bind on with -l")
	}

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	if e.cfg.Seed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(e.cfg.Seed))
	}

	// Generate a key pair for this host. We will use it
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", e.cfg.Port)),
		libp2p.Identity(priv),
	}

	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	// Build host multiaddress
	hostAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))
	if err != nil {
		return nil, err
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

	fullAddr := addr.Encapsulate(hostAddr)
	e.log.Info(fmt.Sprintf("I am %s\n", fullAddr))
	if e.cfg.Secio {
		e.log.Info(fmt.Sprintf("Now run \"go run main.go -p2p_port %d -p2p_target %s -secio\" on a different terminal", e.cfg.Port+101, fullAddr))
	} else {
		e.log.Info(fmt.Sprintf("Now run \"go run main.go -p2p_port %d -p2p_target %s\" on a different terminal", e.cfg.Port+101, fullAddr))
	}

	e.host = basicHost
	return basicHost, nil
}

func (e *EndPoint) handleStream(s net.Stream) {

	e.log.Info("Got a new stream p2p")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go e.readData(rw)
	go e.writeData(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}

type message struct {
	Name  event.EventType
	Value []byte
}
