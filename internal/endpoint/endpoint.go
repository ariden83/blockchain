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
	"strconv"
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
		handle.Handle(fmt.Errorf("fail to open DB files"))
	}

	lastHash, err := e.persistence.GetLastHash()
	handle.Handle(err)

	if lastHash == nil {
		lastHash = e.createGenesis()

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

func (e *EndPoint) createGenesis() []byte {
	var genesisData = "First Transaction from Genesis" // This is arbitrary public key for our genesis data
	cbtx := e.transaction.CoinBaseTx(e.config.Address, genesisData)
	genesis := blockchain.Genesis(cbtx)
	fmt.Println("Genesis proved")

	firstHash := genesis.Hash

	serializeBLock, err := utils.Serialize(genesis)
	handle.Handle(err)

	err = e.persistence.Update(firstHash, serializeBLock)
	handle.Handle(err)
	return firstHash
}

func (e *EndPoint) ListenHTTP(stop chan error) {
	mux := e.makeMuxRouter()
	log.Println("HTTP Server Listening on port :", strconv.Itoa(e.config.Port))
	e.server = &http.Server{
		Addr:           ":" + strconv.Itoa(e.config.Port),
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
	muxRouter.HandleFunc("/blockchain", e.handlePrintBlockChain).Methods("GET")
	muxRouter.HandleFunc("/balance", e.handleGetBalance).Methods("POST")
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
