package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"context"
	"fmt"
	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/endpoint"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/genesis"
	"github.com/ariden83/blockchain/internal/logger"
	"github.com/ariden83/blockchain/internal/metrics"
	"github.com/ariden83/blockchain/internal/p2p"
	"github.com/ariden83/blockchain/internal/persistence"
	"github.com/ariden83/blockchain/internal/transactions"
	"github.com/ariden83/blockchain/internal/wallet"
	"go.uber.org/zap"
	"runtime"
	"syscall"
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

	trans := transactions.Init(cfg.Transactions, per, logs)

	wallets, err := wallet.Init(cfg.Wallet)
	if err != nil {
		logs.Fatal("fail to init wallet", zap.Error(err))
	}

	mtc := metrics.New(cfg.Metrics)

	evt := event.New()

	server := endpoint.Init(cfg, per, trans, wallets, mtc, logs, evt)
	stop := make(chan error, 1)

	server.ListenMetrics(stop)

	var p *p2p.EndPoint
	p = p2p.Init(cfg.P2P, per, wallets, logs, evt)
	if p.Enabled() {
		p.Listen(stop)
	}

	gen := genesis.New(cfg, per, trans, p, evt)
	gen.Load(stop)

	if cfg.API.Enabled {
		server.ListenHTTP(stop)
	}

	/**
	 * And wait for shutdown via signal or error.
	 */
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGKILL, syscall.SIGINT)
		<-sig
		stop <- fmt.Errorf("received Interrupt signal")
	}()

	err = <-stop
	logs.Error("end service", zap.Error(err))

	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(stopCtx)
	p.Shutdown()
}
