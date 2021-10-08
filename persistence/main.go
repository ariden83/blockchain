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
	"github.com/ariden83/blockchain/internal/persistence"
	"github.com/ariden83/blockchain/internal/transactions"
	"github.com/joho/godotenv"
)

func main() {
	defer os.Exit(0)

	conf := config.New()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	per := persistence.Init(conf)
	trans := transactions.Init(conf, per)
	server := endpoint.Init(conf, per, trans)

	stop := make(chan error, 1)
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

	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(stopCtx)
}
