package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/cmd/light/internal/persistance"
	"github.com/ariden83/blockchain/config"
	httpEndpoint "github.com/ariden83/blockchain/internal/endpoint/http"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/logger"
	"github.com/ariden83/blockchain/internal/metrics"
	"github.com/ariden83/blockchain/internal/p2p"
	"github.com/ariden83/blockchain/internal/transactions"
	"github.com/ariden83/blockchain/internal/wallet"
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

	wallets, err := wallet.Init(cfg.Wallet, logs)
	if err != nil {
		logs.Fatal("fail to init wallet", zap.Error(err))
	}

	evt := event.New()
	stop := make(chan error, 1)

	per := &persistance.Persistence{}
	trans := transactions.New(
		transactions.WithPersistence(per),
		transactions.WithLogs(logs),
		transactions.WithEvents(evt),
		transactions.WithConfig(cfg.Transactions))

	mtc := metrics.New(cfg.Metrics)

	server := httpEndpoint.New(httpEndpoint.WithPersistence(per),
		httpEndpoint.WithTransactions(trans),
		httpEndpoint.WithMetrics(mtc),
		httpEndpoint.WithLogs(logs),
		httpEndpoint.WithWallets(wallets),
		httpEndpoint.WithEvents(evt),
		httpEndpoint.WithUserAddress(cfg.Address),
		httpEndpoint.WithConfig(cfg.API))

	var p *p2p.EndPoint
	p = p2p.Init(cfg.P2P, per, wallets, logs, evt, p2p.WithXCache(cfg.XCache))
	p.Listen(stop)
	p.PushMsgForFiles(stop)

	go server.Listen()

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
