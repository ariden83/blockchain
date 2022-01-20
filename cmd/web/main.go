package main

import (
	"context"
	"fmt"
	"github.com/ariden83/blockchain/cmd/web/config"
	"github.com/ariden83/blockchain/cmd/web/internal/auth"
	"github.com/ariden83/blockchain/cmd/web/internal/explorer"
	"github.com/ariden83/blockchain/cmd/web/internal/metrics"
	"github.com/ariden83/blockchain/cmd/web/internal/model"
	"github.com/ariden83/blockchain/internal/logger"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	defer cleanExit()
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("fail to init persistence %s", err)
	}

	logs := logger.InitLight(cfg.Log)
	logs = logs.With(zap.String("v", cfg.Version))
	defer logs.Sync()

	m, err := model.New(cfg.BlockchainAPI, logs)
	if err != nil {
		logs.Fatal("fail to init model", zap.Error(err))
	}

	e := explorer.New(
		explorer.WithConfig(cfg),
		explorer.WithLogs(logs),
		explorer.WithModel(m),
		explorer.WithMetadata(cfg.Metadata),
		explorer.WithRecaptcha(cfg.ReCaptcha, logs),
		explorer.WithAuth(auth.New(
			// auth.WithGoogleAPI(cfg.Auth.GoogleAPI),
			auth.WithClassic(cfg.Auth.Classic),
		)),
		explorer.WithMetrics(metrics.New(cfg.Name)),
		explorer.WithLocales(cfg.Locales),
	)

	stop := make(chan error, 1)

	e.StartMetricsServer(stop)
	e.Start(stop)

	logs.Info("Close web server")

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
	e.Shutdown(stopCtx)
}

func cleanExit() {
	os.Exit(0)
}
