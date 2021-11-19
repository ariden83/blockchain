package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/metrics"
	"github.com/ariden83/blockchain/internal/middleware"
	"github.com/ariden83/blockchain/internal/persistence"
	"github.com/ariden83/blockchain/internal/transactions"
	"github.com/ariden83/blockchain/internal/utils"
	"github.com/ariden83/blockchain/internal/wallet"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

type EndPoint struct {
	cfg           *config.Config
	persistence   persistence.IPersistence
	transaction   transactions.ITransaction
	server        *http.Server
	metricsServer *http.Server
	wallets       wallet.IWallets
	metrics       *metrics.Metrics
	log           *zap.Logger
	event         *event.Event
}

type Healthz struct {
	Result   bool     `json:"result"`
	Messages []string `json:"messages"`
	Version  string   `json:"version"`
}

func Init(
	cfg *config.Config,
	per persistence.IPersistence,
	trans transactions.ITransaction,
	wallets wallet.IWallets,
	mtcs *metrics.Metrics,
	logs *zap.Logger,
	evt *event.Event,
) *EndPoint {
	e := &EndPoint{
		cfg:         cfg,
		persistence: per,
		transaction: trans,
		wallets:     wallets,
		metrics:     mtcs,
		log:         logs.With(zap.String("service", "http")),
		event:       evt,
	}

	return e
}

func (e *EndPoint) Genesis() {
	go func() {
		var lastHash []byte

		// si les fichiers locaux n'existent pas
		if !e.persistence.DBExists() {
			e.Handle(fmt.Errorf("fail to open DB files"))
		}

		lastHash, err := e.persistence.GetLastHash()
		e.Handle(err)

		if lastHash == nil {
			lastHash = e.createGenesis()

		} else {

			val, err := e.persistence.GetCurrentHashSerialize(lastHash)
			e.Handle(err)
			block, err := utils.Deserialize(val)
			e.Handle(err)

			e.persistence.SetLastHash(lastHash)

			mutex.Lock()
			blockchain.BlockChain = append(blockchain.BlockChain, *block)
			mutex.Unlock()

			spew.Dump(blockchain.BlockChain)
		}

	}()
}

func (e *EndPoint) createGenesis() []byte {
	var genesisData string = "First Transaction from Genesis" // This is arbitrary public key for our genesis data
	cbtx := e.transaction.CoinBaseTx(e.cfg.Address, genesisData)
	genesis := blockchain.Genesis(cbtx)
	fmt.Println("Genesis proved")

	firstHash := genesis.Hash

	serializeBLock, err := utils.Serialize(genesis)
	e.Handle(err)

	err = e.persistence.Update(firstHash, serializeBLock)
	e.Handle(err)
	return firstHash
}

func (e *EndPoint) ListenHTTP(stop chan error) {
	go func() {
		if err := e.StartServer(); err != nil {
			stop <- fmt.Errorf("cannot start server HTTP %s", err)
		}
	}()
}

func (e *EndPoint) StartServer() error {
	e.log.Info("Start listening HTTP Server", zap.Int("port", e.cfg.API.Port))

	mux := e.makeMuxRouter()

	e.server = &http.Server{
		Addr:           ":" + strconv.Itoa(e.cfg.API.Port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := e.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		e.log.Error("fail to listen http", zap.Error(err))
		e.cfg.API.Port = e.cfg.API.Port + 1
		return e.StartServer()
	}
	return err
}

func (e *EndPoint) makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/blockchain", e.handleGetBlockChain).Methods("GET")
	muxRouter.HandleFunc("/balance", e.handleGetBalance).Methods("POST")
	muxRouter.HandleFunc("/write", e.handleWriteBlock).Methods("POST")
	muxRouter.HandleFunc("/send", e.handleSendBlock).Methods("POST")
	muxRouter.HandleFunc("/wallets", e.handlePrintWallets).Methods("GET")
	muxRouter.HandleFunc("/wallet", e.handleCreateWallet).Methods("POST")
	muxRouter.HandleFunc("/mywallet", e.handleMyWallet).Methods("POST")
	muxRouter.HandleFunc("/address", e.handleGetServersAddress).Methods("GET")

	muxRouter.Use(middleware.DefaultHeader)
	muxRouter.Use(e.MetricsMiddleware)

	return muxRouter
}

func (e *EndPoint) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := strings.ToLower(r.Method)

		jsonHandler := promhttp.InstrumentHandlerInFlight(
			e.metrics.InFlight,

			promhttp.InstrumentHandlerResponseSize(
				e.metrics.ResponseSize.MustCurryWith(prometheus.Labels{"service": route}),

				promhttp.InstrumentHandlerRequestSize(
					e.metrics.RequestSize.MustCurryWith(prometheus.Labels{"service": route}),

					promhttp.InstrumentHandlerCounter(
						e.metrics.RouteCountReqs.MustCurryWith(prometheus.Labels{"service": route}),

						promhttp.InstrumentHandlerDuration(
							e.metrics.ResponseDuration.MustCurryWith(prometheus.Labels{"service": route}),
							next)))))

		jsonHandler.ServeHTTP(w, r)
	})
}

func (e *EndPoint) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		e.log.Error("HTTP 500: Internal Server Error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write([]byte("HTTP 500: Internal Server Error")); err != nil {
			e.log.Error("fail to write response", zap.Error(err))
		}
		return
	}
	w.WriteHeader(code)
	if _, err = w.Write(response); err != nil {
		e.log.Error("fail to write response", zap.Error(err))
	}
}

func (e *EndPoint) Shutdown(ctx context.Context) {
	e.persistence.Close()
	err := e.server.Shutdown(ctx)
	if err != nil {
		e.log.Error("fail to shutdown server", zap.Error(err))
	}
}
