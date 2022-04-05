package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/config"
	grpcEndpoint "github.com/ariden83/blockchain/internal/endpoint/grpc"
	httpEndpoint "github.com/ariden83/blockchain/internal/endpoint/http"
	metricsEndpoint "github.com/ariden83/blockchain/internal/endpoint/metrics"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/genesis"
	"github.com/ariden83/blockchain/internal/logger"
	"github.com/ariden83/blockchain/internal/metrics"
	"github.com/ariden83/blockchain/internal/p2p"
	"github.com/ariden83/blockchain/internal/persistence"
	"github.com/ariden83/blockchain/internal/transactions"
	"github.com/ariden83/blockchain/internal/wallet"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("fail to init persistence %s", err)
	}

	if cfg.Threads > 0 {
		runtime.GOMAXPROCS(cfg.Threads)
		log.Printf("Running with %v threads", cfg.Threads)

	} else {
		n := runtime.NumCPU()
		runtime.GOMAXPROCS(n)
		log.Printf("Running with default %v threads", n)
	}

	logs := logger.Init(cfg.Log)
	logs = logs.With(zap.String("v", cfg.Version))
	defer logs.Sync()

	per, err := persistence.Init(cfg.Database)
	if err != nil {
		logs.Fatal("fail to init persistence", zap.Error(err))
	}
	defer per.Close()

	evt := event.New()

	trans := transactions.New(
		transactions.WithPersistence(per),
		transactions.WithLogs(logs),
		transactions.WithEvents(evt),
		transactions.WithTraces(cfg.Transactions.Trace),
		transactions.WithConfig(cfg.Transactions))

	wallets, err := wallet.Init(cfg.Wallet, logs)
	if err != nil {
		logs.Fatal("fail to init wallet", zap.Error(err))
	}
	defer wallets.Close()

	mtc := metrics.New(cfg.Metrics)

	s := Server{}

	s.httpServer = httpEndpoint.New(httpEndpoint.WithPersistence(per),
		httpEndpoint.WithTransactions(trans),
		httpEndpoint.WithMetrics(mtc),
		httpEndpoint.WithLogs(logs),
		httpEndpoint.WithWallets(wallets),
		httpEndpoint.WithEvents(evt),
		httpEndpoint.WithUserAddress(cfg.Address),
		httpEndpoint.WithConfig(cfg.API))

	s.grpcServer = grpcEndpoint.New(grpcEndpoint.WithPersistence(per),
		grpcEndpoint.WithTransactions(trans),
		grpcEndpoint.WithMetrics(mtc),
		grpcEndpoint.WithLogs(logs),
		grpcEndpoint.WithWallets(wallets),
		grpcEndpoint.WithEvents(evt),
		grpcEndpoint.WithUserAddress(cfg.Address),
		grpcEndpoint.WithConfig(cfg.GRPC))

	s.metricsServer = metricsEndpoint.New(cfg.Metrics, logs)

	stop := make(chan error, 1)

	var p *p2p.EndPoint
	p = p2p.Init(cfg.P2P, per, wallets, logs, evt)
	if p.Enabled() {
		p.Listen()
	}

	gen := genesis.New(cfg, per, trans, p, evt)
	gen.Load(stop)

	s.Start(stop)
	/**
	 * And wait for shutdown via signal or error.
	 */
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGKILL, syscall.SIGINT)
		err := <-sig
		stop <- fmt.Errorf("received Interrupt signal %v", err)
	}()

	err = <-stop

	logs.Info("kill service", zap.Error(err))

	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.Shutdown(stopCtx)
}
