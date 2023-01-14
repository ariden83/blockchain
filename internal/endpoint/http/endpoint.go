package http

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/metrics"
	"github.com/ariden83/blockchain/internal/persistence"
	"github.com/ariden83/blockchain/internal/transactions"
	"github.com/ariden83/blockchain/internal/wallet"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var mutex = &sync.Mutex{}

type EndPoint struct {
	cfg           config.API
	event         *event.Event
	log           *zap.Logger
	metrics       *metrics.Metrics
	metricsServer *http.Server
	persistence   persistenceadapter.Adapter
	server        *http.Server
	transaction   transactions.ITransaction
	userAddress   string
	wallets       wallet.IWallets
}

func New(options ...func(*EndPoint)) *EndPoint {
	ep := &EndPoint{}

	for _, o := range options {
		o(ep)
	}

	return ep
}

func WithConfig(cfg config.API) func(*EndPoint) {
	return func(e *EndPoint) {
		e.cfg = cfg
	}
}

func WithPersistence(p persistenceadapter.Adapter) func(*EndPoint) {
	return func(e *EndPoint) {
		e.persistence = p
	}
}

func WithTransactions(t transactions.ITransaction) func(*EndPoint) {
	return func(e *EndPoint) {
		e.transaction = t
	}
}

func WithWallets(w wallet.IWallets) func(*EndPoint) {
	return func(e *EndPoint) {
		e.wallets = w
	}
}

func WithMetrics(m *metrics.Metrics) func(*EndPoint) {
	return func(e *EndPoint) {
		e.metrics = m
	}
}

func WithLogs(logs *zap.Logger) func(*EndPoint) {
	return func(e *EndPoint) {
		e.log = logs.With(zap.String("service", "http"))
	}
}

func WithEvents(evt *event.Event) func(*EndPoint) {
	return func(e *EndPoint) {
		e.event = evt
	}
}

func WithUserAddress(a string) func(*EndPoint) {
	return func(e *EndPoint) {
		e.userAddress = a
	}
}

func (e *EndPoint) Enabled() bool {
	return e.cfg.Enabled
}

func (e *EndPoint) Listen() error {
	e.log.Info("Start listening HTTP Server", zap.Int("port", e.cfg.Port))

	mux := e.makeMuxRouter()

	e.server = &http.Server{
		Addr:           e.cfg.Host + ":" + strconv.Itoa(e.cfg.Port),
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
