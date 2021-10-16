package main

import (
	"fmt"
	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/logger"
	"github.com/ariden83/blockchain/internal/p2p"
	"github.com/ariden83/blockchain/internal/wallet"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
}

type Persistence struct{}

func (p *Persistence) GetLastHash() ([]byte, error) {
	return []byte{}, nil
}
func (p *Persistence) Update(lastHash []byte, hashSerialize []byte) error {
	return nil
}

func main() {

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("fail to init persistence %s", err)
	}

	cfg.Wallet.WithFile = false
	cfg.Log.WithFile = false

	logs := logger.InitLight(cfg.Log)
	logs = logs.With(zap.String("v", cfg.Version))
	defer logs.Sync()

	wallets, err := wallet.Init(cfg.Wallet)
	if err != nil {
		logs.Fatal("fail to init wallet", zap.Error(err))
	}

	evt := event.New()
	stop := make(chan error, 1)

	per := &Persistence{}

	var p *p2p.EndPoint
	p = p2p.Init(cfg.P2P, per, wallets, logs, evt)
	p.Listen(stop)

	/**
	 * And wait for shutdown via signal or error.
	 */
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig
		stop <- fmt.Errorf("received Interrupt signal")
	}()

	err = <-stop
}

type message struct {
	Name  event.EventType
	Value []byte
}

type Seed struct {
	Address   string
	Timestamp string
	PubKey    string
	PrivKey   string
	Mnemonic  string
}
