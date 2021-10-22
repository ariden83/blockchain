package main

import (
	"fmt"
	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/endpoint"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/logger"
	"github.com/ariden83/blockchain/internal/metrics"
	"github.com/ariden83/blockchain/internal/p2p"
	"github.com/ariden83/blockchain/internal/wallet"

	"github.com/ariden83/blockchain/internal/transactions"

	"github.com/ariden83/blockchain/example/light/internal/persistance"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
)

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

	per := &persistance.Persistence{}
	trans := transactions.Init(cfg.Transactions, per, logs)
	mtc := metrics.New(cfg.Metrics)

	server := endpoint.Init(cfg, per, trans, wallets, mtc, logs, evt)

	var p *p2p.EndPoint
	p = p2p.Init(cfg.P2P, per, wallets, logs, evt, p2p.WithXCache(cfg.XCache))
	p.Listen()
	p.PushMsgForFiles()

	server.ListenHTTP(stop)

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
	logs.Error("end service", zap.Error(err))
}
