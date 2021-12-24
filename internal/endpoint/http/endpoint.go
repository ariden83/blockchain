package http

import (
	"context"
	"errors"
	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/metrics"
	"github.com/ariden83/blockchain/internal/persistence"
	"github.com/ariden83/blockchain/internal/transactions"
	"github.com/ariden83/blockchain/internal/utils"
	"github.com/ariden83/blockchain/internal/wallet"
	"github.com/davecgh/go-spew/spew"
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
	cfg           config.API
	persistence   persistence.IPersistence
	transaction   transactions.ITransaction
	server        *http.Server
	metricsServer *http.Server
	wallets       wallet.IWallets
	metrics       *metrics.Metrics
	log           *zap.Logger
	event         *event.Event
	userAddress   string
}

func New(
	cfg config.API,
	per persistence.IPersistence,
	trans transactions.ITransaction,
	wallets wallet.IWallets,
	mtcs *metrics.Metrics,
	logs *zap.Logger,
	evt *event.Event,
	userAddress string,
) *EndPoint {
	e := &EndPoint{
		cfg:         cfg,
		persistence: per,
		transaction: trans,
		wallets:     wallets,
		metrics:     mtcs,
		log:         logs.With(zap.String("service", "http")),
		event:       evt,
		userAddress: userAddress,
	}

	return e
}

func (e *EndPoint) Genesis() {
	go func() {
		var lastHash []byte

		// si les fichiers locaux n'existent pas
		if !e.persistence.DBExists() {
			e.Handle(errors.New("fail to open DB files"))
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
	cbtx := e.transaction.CoinBaseTx(e.userAddress, genesisData)
	genesis := blockchain.Genesis(cbtx)
	e.log.Info("Genesis proved")

	firstHash := genesis.Hash

	serializeBLock, err := utils.Serialize(genesis)
	e.Handle(err)

	err = e.persistence.Update(firstHash, serializeBLock)
	e.Handle(err)
	return firstHash
}

func (e *EndPoint) Listen() error {
	e.log.Info("Start listening HTTP Server", zap.Int("port", e.cfg.Port))

	mux := e.makeMuxRouter()

	e.server = &http.Server{
		Addr:           ":" + strconv.Itoa(e.cfg.Port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := e.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		e.log.Error("fail to listen http", zap.Error(err))
		e.cfg.Port = e.cfg.Port + 1
		return e.Listen()
	}
	return err
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

func (e *EndPoint) Shutdown(ctx context.Context) {
	err := e.server.Shutdown(ctx)
	if err != nil {
		e.log.Error("fail to shutdown server", zap.Error(err))
	}
}
