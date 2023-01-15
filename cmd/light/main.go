package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/config"
	httpEndpoint "github.com/ariden83/blockchain/internal/endpoint/http"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/logger"
	"github.com/ariden83/blockchain/internal/metrics"
	"github.com/ariden83/blockchain/internal/p2p"
	persistencefactory "github.com/ariden83/blockchain/internal/persistence/factory"
	transactionfactory "github.com/ariden83/blockchain/internal/transaction/factory"
	"github.com/ariden83/blockchain/internal/transaction/impl/transaction"
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

	wallets, err := wallet.New(cfg.Wallet, logs)
	if err != nil {
		logs.Fatal("fail to init wallet", zap.Error(err))
	}

	evt := event.New()
	stop := make(chan error, 1)

	per, err := persistencefactory.New(persistencefactory.Config{
		Implementation: persistencefactory.ImplementationStub,
	})
	if err != nil {
		logs.Fatal("fail to init persistence", zap.Error(err))
	}

	trans, err := transactionfactory.New(transactionfactory.Config{Implementation: transactionfactory.ImplementationTransaction},
		transaction.WithPersistence(per),
		transaction.WithLogs(logs),
		transaction.WithEvents(evt),
		transaction.WithConfig(cfg.Transactions))
	if err != nil {
		logs.Fatal("fail to init transaction", zap.Error(err))
	}

	mtc := metrics.New(cfg.Metrics, nil)

	server := httpEndpoint.New(httpEndpoint.WithPersistence(per),
		httpEndpoint.WithTransactions(trans),
		httpEndpoint.WithMetrics(mtc),
		httpEndpoint.WithLogs(logs),
		httpEndpoint.WithWallets(wallets),
		httpEndpoint.WithEvents(evt),
		httpEndpoint.WithUserAddress(cfg.Address),
		httpEndpoint.WithConfig(cfg.API))

	var p *p2p.EndPoint
	p = p2p.New(cfg.P2P, per, wallets, logs, evt, p2p.WithXCache(cfg.XCache))
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
