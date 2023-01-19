// Package P2P represent a peer to peer linked on go-libp2p.
package p2p

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"go.uber.org/zap"
	"io"
	mrand "math/rand"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/core/routing"

	pstore "github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/satori/go.uuid"

	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/p2p/address"
	"github.com/ariden83/blockchain/internal/persistence"
	"github.com/ariden83/blockchain/internal/wallet"
	"github.com/ariden83/blockchain/internal/xcache"
)

type EndPoint struct {
	address     []string
	cfg         config.P2P
	dbLoad      bool
	event       *event.Event
	enabled     bool
	host        host.Host
	linked      bool
	log         *zap.Logger
	msgReceived []string
	persistence persistenceadapter.Adapter
	readerReady bool
	target      string
	wallets     wallet.IWallets
	xCache      *xcache.Cache
	writerReady bool
}

// Option is the type of option passed to the constructor.
type Option func(e *EndPoint)

// New represent a new peer to peer adapter.
func New(
	cfg config.P2P,
	per persistenceadapter.Adapter,
	wallets wallet.IWallets,
	logs *zap.Logger,
	evt *event.Event,
	opts ...Option,
) *EndPoint {

	e := &EndPoint{
		cfg:         cfg,
		enabled:     cfg.Enabled,
		event:       evt,
		log:         logs.With(zap.String("service", "p2p")),
		persistence: per,
		target:      cfg.Target,
		wallets:     wallets,
	}

	for _, o := range opts {
		o(e)
	}

	return e
}

// WithXCache offer the possibility to add a cache system to the peer to peer adapter.
func WithXCache(cfg config.XCache) Option {
	return func(e *EndPoint) {
		var err error

		e.xCache, err = xcache.New(
			xcache.WithSize(int32(cfg.Size)),
			xcache.WithTTL(time.Duration(cfg.TTL)*time.Second),
			xcache.WithNegSize(int32(cfg.NegSize)),
			xcache.WithNegTTL(time.Duration(cfg.NegTTL)*time.Second),
			xcache.WithStale(true),
			xcache.WithPruneSize(int32(cfg.Size/20)+1))

		if err != nil {
			e.log.Error("fail to init xCache for p2p service", zap.Error(err))
		}
	}
}

// Enabled lets you know if the peer to peer system is activated.
func (e *EndPoint) Enabled() bool {
	return e.enabled
}

// Listen start to listen the peer to peer connexion.
func (e *EndPoint) Listen(stop chan error) {
	e.hasRequiredPort()

	hasConnexion := false
	for !hasConnexion {
		hasBasicHost := false
		for !hasBasicHost {
			// try to connect to an existant host
			if err := e.makeBasicHost(); err != nil {
				/* if strings.Contains(fmt.Sprint(err), protocolError) || strings.Contains(fmt.Sprint(err), addressAMReadyUseError) {}*/
				time.Sleep(time.Millisecond * 10)
				e.log.Error("fail to listen p2p", zap.Int("port", e.cfg.Port))
				e.cfg.Port = e.cfg.Port + 1

			} else {
				e.log.Info("basic host connexion created")
				hasBasicHost = true
			}

			select {
			case <-stop: // closes when the caller cancels the ctx
				return
			default:
			}
		}

		err := e.retryConnectToIPFS()
		if err == nil {
			e.log.Info("connected to ipfs")
			hasConnexion = true
		} else {
			e.log.Warn("fail to connect to ipfs", zap.Error(err))
		}

		select {
		case <-stop: // closes when the caller cancels the ctx
			hasConnexion = true
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}

	go e.alertWaitFirstConnexion(stop)
}

func (e *EndPoint) retryConnectToIPFS() error {
	e.log.Info("listening for new connections", zap.Int("port", e.cfg.Port))
	// Set a stream handler on host A. /p2p/1.0.0 is
	// a user-defined protocol name.
	e.host.SetStreamHandler("/p2p/1.0.0", e.handleStream)

	if !e.HasTarget() {
		return nil
	}

	if err := e.connectToIPFS(e.host); err != nil {
		/* if strings.Contains(fmt.Sprint(err), failNegociateError) || strings.Contains(fmt.Sprint(err), protocolError) ||
			strings.Contains(fmt.Sprint(err), noGoodAddress) {
			e.log.Error("fail to listen p2p", zap.Int("port", e.cfg.Port))
			e.cfg.Port = e.cfg.Port + 1
			return err
		} */
		e.log.Error("fail to listen p2p", zap.Int("port", e.cfg.Port))
		e.cfg.Port = e.cfg.Port + 1
		return err
	}
	return nil
}

func (e *EndPoint) PushMsgForFiles(stop chan error) {
	go func() {
		for {
			select {
			case <-stop: // closes when the caller cancels the ctx
				return
			default:
			}

			// wait 10 millisecon before retry
			time.Sleep(10 * time.Millisecond)
			if e.writerReady && e.readerReady {
				id := uuid.NewV4().String()
				e.event.Push(event.Message{
					Type: event.Files,
					ID:   id,
				})
				e.saveMsgReceived(id)

				break
			}
		}
	}()
}

func (e *EndPoint) hasRequiredPort() {
	if e.cfg.Port == 0 {
		e.log.Fatal("Please provide a port to bind on with -l")
	}
}

func (e *EndPoint) SetTarget(target string) {
	e.target = target
}

func (e *EndPoint) Target() string {
	return e.target
}

func (e *EndPoint) HasTarget() bool {
	if e.target == "" {
		// call default genesis
		return false
	}
	return true
}

func (e *EndPoint) connectToIPFS(ha host.Host) error {
	// The following code extracts target's peer ID from the
	// given multiaddress
	ipfsAddr, err := ma.NewMultiaddr(e.target)
	if err != nil {
		e.log.Error("fail to set new multi address", zap.Error(err), zap.String("target", e.target))
		return err
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
		return err
	}
	// Nous nous retrouvons avec le peerID et l'adresse cible targetAddr de l'hôte auquel nous voulons nous connecter
	// et ajoutons cet enregistrement dans notre "magasin"
	// afin que nous puissions garder une trace de qui nous sommes connectés.
	// Nous le faisons avec ha.Peerstore().AddAddr
	peerid, err := peer.Decode(pid)
	if err != nil {
		e.log.Error("fail to get decode peer", zap.Error(err))
		return err
	}
	// Decapsulate the /ipfs/<peerID> part from the target
	// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
	targetPeerAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", peer.Encode(peerid)))
	if err != nil {
		e.log.Error("fail to set new multi address", zap.Error(err), zap.String("target", e.target))
		return err
	}

	targetAddr := ipfsAddr.Decapsulate(targetPeerAddr)

	// We have a peer ID and a targetAddr so we add it to the peerstore
	// so LibP2P knows how to contact it
	ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

	e.log.Info("opening stream p2p", zap.String("target", e.target))
	// make a new stream from host B to host A
	// it should be handled on host A by the handler we set above because
	// we use the same /p2p/1.0.0 protocol
	protocolID := protocol.ID(e.cfg.ProtocolID)
	s, err := ha.NewStream(context.Background(), peerid, protocolID)
	if err != nil {
		e.log.Error("fail to set new stream",
			zap.Error(err),
			zap.Any("peer_id", peerid),
			zap.Any("protocol_id", protocolID))
		return err
	}
	// Create a buffered stream so that read and writes are non blocking.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	// Create a thread to read and write data.
	e.writeData(rw)
	e.readData(rw)
	return nil
}

func (e *EndPoint) makeBasicHost() error {
	// The context governs the lifetime of the libp2p node.
	// Cancelling it will stop the the host.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
	// Set your own keypair
	/* priv, _, err := crypto.GenerateKeyPair(
		crypto.Ed25519, // Select your key type. Ed25519 are nice short
		-1,             // Select key length when possible (i.e. RSA).
	) */
	if err != nil {
		return err
	}

	var idht *dht.IpfsDHT

	connmgr, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		return err
	}

	e.host, err = libp2p.New( // Use the keypair we generated
		libp2p.Identity(priv),
		// Multiple listen addresses
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", e.cfg.Port),      // regular tcp connections
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/quic", e.cfg.Port), // a UDP endpoint for the QUIC transport
		),
		// support TLS connections
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		// support noise connections
		libp2p.Security(noise.ID, noise.New),
		// support any other default transports (TCP)
		libp2p.DefaultTransports,
		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(connmgr),
		// Attempt to open ports using uPNP for NATed hosts.
		libp2p.NATPortMap(),
		// Let this host use the DHT to find other hosts
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err = dht.New(ctx, h)
			return idht, err
		}),
		// If you want to help other peers to figure out if they are behind
		// NATs, you can launch the server-side of AutoNAT too (AutoRelay
		// already runs the client)
		//
		// This service is highly rate-limited and should not cause any
		// performance issues.
		libp2p.EnableNATService(),
	)
	if err != nil {
		e.log.Error("fail to set basic host", zap.Int("port", e.cfg.Port))
		return err
	}

	e.log.Info("P2P start:", zap.Any("address", e.host.Addrs()), zap.Any("host_id", e.host.ID()))
	e.setStreamHandler()

	// Parse the multiaddr string.
	peerMA, err := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", e.host.ID().Pretty()))
	//peerMA, err := ma.NewMultiaddr(e.target)
	if err != nil {
		return err
	}

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addrs := e.host.Addrs()
	var addr ma.Multiaddr
	// select the address starting with "ip4"
	for _, i := range addrs {
		if strings.HasPrefix(i.String(), "/ip4") {
			addr = i
			break
		}
	}

	address.IAM.SetAddress(addr.Encapsulate(peerMA).String())
	e.log.Info(fmt.Sprintf("I am %s\n", address.IAM.Address()))
	if e.cfg.Secio {
		e.log.Info(fmt.Sprintf("Now run \"go run main.go -p2p_target %s -secio\" on a different terminal", address.IAM.Address()))
	} else {
		e.log.Info(fmt.Sprintf("Now run \"go run main.go -p2p_target %s\" on a different terminal", address.IAM.Address()))
	}

	return nil
}

func (e *EndPoint) alertWaitFirstConnexion(stop chan error) {
	if e.target != "" {
		return
	}

	for {
		if e.linked {
			break
		}
		e.log.Warn("waiting for new P2P connexion ...")
		if e.cfg.Secio {
			e.log.Info(fmt.Sprintf("run \"go run main.go -p2p_target %s -secio\" on a different terminal", address.IAM.Address()))
		} else {
			e.log.Info(fmt.Sprintf("run \"go run main.go -p2p_target %s\" on a different terminal", address.IAM.Address()))
		}

		select {
		case <-stop: // closes when the caller cancels the ctx
			break
		default:
		}

		time.Sleep(30 * time.Minute)
	}
}

// Setup a stream handler.
//
// This gets called every time a peer connects and opens a stream to this node.
func (e *EndPoint) setStreamHandler() {
	protocolID := protocol.ID(e.cfg.ProtocolID)
	e.host.SetStreamHandler(protocolID, e.handleStream)
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

func (e *EndPoint) handleStream(stream network.Stream) {
	e.linked = true
	e.log.Info("Got a new stream p2p")
	e.address = append(e.address)
	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	e.readData(rw)
	e.writeData(rw)
	// stream 's' will stay open until you close it (or the other side closes it).
}

func (e *EndPoint) Shutdown() {
	e.host.Close()
}
