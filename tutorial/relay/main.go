package main

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/libp2p/go-libp2p/p2p/net/swarm"
	circuit "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	run()
}

// https://docs.libp2p.io/concepts/circuit-relay/
func run() {
	// Créez trois hôtes libp2p, activez les capacités de client de relais sur chacun d'eux
	// Dites à l'hôte d'utiliser des relais
	h1, err := libp2p.New(libp2p.EnableRelay())
	if err != nil {
		log.Printf("Failed to create h1: %v", err)
		return
	}

	// Dites à l'hôte de relayer les connexions pour d'autres pairs (la possibilité d'*utiliser*
	// un relais vs la capacité d'*être* un relais)
	h2, err := libp2p.New(libp2p.DisableRelay())
	if err != nil {
		log.Printf("Failed to create h2: %v", err)
		return
	}
	_, err = circuit.New(h2, nil)
	if err != nil {
		log.Printf("Failed to instantiate h2 relay: %v", err)
		return
	}

	// Mettez à zéro les adresses d'écoute pour l'hôte, afin qu'il ne puisse communiquer que via le circuit p2p pour notre exemple
	h3, err := libp2p.New(libp2p.ListenAddrs(), libp2p.EnableRelay())
	if err != nil {
		log.Printf("Failed to create h3: %v", err)
		return
	}

	h2info := peer.AddrInfo{
		ID:    h2.ID(),
		Addrs: h2.Addrs(),
	}

	log.Printf("h2info: %v", h2info)
	// Connectez à la fois h1 et h3 à h2, mais pas l'un à l'autre
	if err := h1.Connect(context.Background(), h2info); err != nil {
		log.Printf("Failed to connect h1 and h2: %v", err)
		return
	}

	if err := h3.Connect(context.Background(), h2info); err != nil {
		log.Printf("Failed to connect h3 and h2: %v", err)
		return
	}

	// Maintenant, pour tester les choses, configurons un gestionnaire de protocole sur h3
	h3.SetStreamHandler("/cats", func(s network.Stream) {
		log.Println("Meow! It worked!")
		s.Close()
	})

	_, err = h1.NewStream(context.Background(), h3.ID(), "/cats")
	if err == nil {
		log.Println("Didnt actually expect to get a stream here. What happened?")
		return
	}
	log.Printf("Okay, no connection from h1 to h3: %v", err)
	log.Println("Just as we suspected")

	// Crée une adresse de relais vers h3 en utilisant h2 comme relais
	relayaddr, err := ma.NewMultiaddr("/p2p/" + h2.ID().Pretty() + "/p2p-circuit/ipfs/" + h3.ID().Pretty())
	if err != nil {
		log.Println(err)
		return
	}

	// Étant donné que nous venons d'essayer et que nous n'avons pas réussi à composer le numéro,
	// le système de numérotation nous empêchera par défaut de recomposer si rapidement.
	// Puisque nous savons ce que nous faisons, nous pouvons utiliser ce vilain hack
	// (il est sur notre liste <TODO pour le rendre un peu plus propre) pour dire au numéroteur "non, ça va, essayons à nouveau"
	h1.Network().(*swarm.Swarm).Backoff().Clear(h3.ID())

	h3relayInfo := peer.AddrInfo{
		ID:    h3.ID(),
		Addrs: []ma.Multiaddr{relayaddr},
	}
	log.Printf("h3relayInfo: %v", h3relayInfo)
	if err := h1.Connect(context.Background(), h3relayInfo); err != nil {
		log.Printf("Failed to connect h1 and h3: %v", err)
		return
	}

	// Woohoo! we're connected!
	s, err := h1.NewStream(context.Background(), h3.ID(), "/cats")
	if err != nil {
		log.Println("huh, this should have worked: ", err)
		return
	}

	s.Read(make([]byte, 1)) // block until the handler closes the stream
}
