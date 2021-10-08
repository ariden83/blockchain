package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/handle"
	"github.com/ariden83/blockchain/internal/persistence"
	"github.com/ariden83/blockchain/internal/transactions"
	"github.com/ariden83/blockchain/internal/utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

type EndPoint struct {
	config      *config.Config
	persistence *persistence.Persistence
	transaction *transactions.Transactions
	server      *http.Server
}

func Init(conf *config.Config, per *persistence.Persistence, trans *transactions.Transactions) *EndPoint {
	e := &EndPoint{
		config:      conf,
		persistence: per,
		transaction: trans,
	}
	go func() {
		e.Genesis()
	}()

	return e
}

func (e *EndPoint) Genesis() {
	var lastHash []byte

	// si les fichiers locaux n'existent pas
	if !e.persistence.DBExists() {

		var genesisData = "First Transaction from Genesis" // This is arbitrary public key for our genesis data

		cbtx := transactions.CoinbaseTx(e.config.Address, genesisData)
		genesis := blockchain.Genesis(cbtx)
		fmt.Println("Genesis proved")

		lastHash = genesis.Hash

		serializeBLock, err := utils.Serialize(genesis)
		handle.Handle(err)

		err = e.persistence.Update(lastHash, serializeBLock)
		handle.Handle(err)
	}

	lastHash, err := e.persistence.GetLastHash()
	handle.Handle(err)

	if lastHash == nil {
		handle.Handle(fmt.Errorf("No blockchain found, please create one first"))

	} else {

		val, err := e.persistence.GetCurrentHashSerialize(lastHash)
		handle.Handle(err)
		block, err := utils.Deserialize(val)
		handle.Handle(err)

		e.persistence.SetLastHash(lastHash)

		mutex.Lock()
		blockchain.BlockChain = append(blockchain.BlockChain, *block)
		mutex.Unlock()

		spew.Dump(blockchain.BlockChain)
	}

	return
}

func (e *EndPoint) ListenHTTP(stop chan error) {
	mux := e.makeMuxRouter()
	httpPort := os.Getenv("PORT")
	log.Println("HTTP Server Listening on port :", httpPort)
	e.server = &http.Server{
		Addr:           ":" + httpPort,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := e.server.ListenAndServe(); err != nil {
		stop <- fmt.Errorf("cannot start server HTTP %s", err)
	}
}

func (e *EndPoint) makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/create", e.handlePrintBlockChain).Methods("GET")
	muxRouter.HandleFunc("/balance", e.handleGetBalance).Methods("GET")
	muxRouter.HandleFunc("/write", e.handleWriteBlock).Methods("POST")
	muxRouter.HandleFunc("/send", e.handleSendBlock).Methods("POST")
	return muxRouter
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func (e *EndPoint) Shutdown(ctx context.Context) {
	e.persistence.Close()
	e.server.Shutdown(ctx)
}
