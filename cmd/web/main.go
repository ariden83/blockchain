package main

import (
	"github.com/ariden83/blockchain/cmd/web/internal/explorer"
	"os"
)

func main() {
	defer cleanExit()
	explorer.Start()
}

func cleanExit() {
	os.Exit(0)
}
