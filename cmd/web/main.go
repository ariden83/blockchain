package main

import (
	"github.com/ariden83/blockchain/cmd/web/config"
	"github.com/ariden83/blockchain/cmd/web/internal/explorer"
	"github.com/ariden83/blockchain/cmd/web/internal/model"
	"github.com/ariden83/blockchain/cmd/web/internal/token"
	"github.com/ariden83/blockchain/internal/logger"
	"go.uber.org/zap"
	"log"
	"os"
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

	m := model.New(cfg, logs)

	t := token.New(cfg.Token)

	explorer.New(cfg, logs, m, t).Start()
}

func cleanExit() {
	os.Exit(0)
}
