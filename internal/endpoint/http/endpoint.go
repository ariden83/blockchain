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
	"github.com/ariden83/blockchain/internal/transaction"
	"github.com/ariden83/blockchain/internal/wallet"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var mutex = &sync.Mutex{}

// EndPoint represent a HTTP endpoint adapter.
type EndPoint struct {
	cfg           config.API
	event         *event.Event
	log           *zap.Logger
	metrics       *metrics.Metrics
	metricsServer *http.Server
	persistence   persistenceadapter.Adapter
	server        *http.Server
	transaction   transaction.Adapter
	userAddress   string
	wallets       wallet.IWallets
}

// New define a new HTTP endpoint adapter.
func New(options ...func(*EndPoint)) *EndPoint {
	ep := &EndPoint{
		log: zap.NewNop(),
	}

	for _, o := range options {
		o(ep)
	}

	return ep
}

// WithConfig offer the possibility to set a config to the endpoint adapter.
func WithConfig(cfg config.API) func(*EndPoint) {
	return func(e *EndPoint) {
		e.cfg = cfg
	}
}

// WithPersistence offer the possibility to set a persistence adapter to the endpoint adapter.
func WithPersistence(p persistenceadapter.Adapter) func(*EndPoint) {
	return func(e *EndPoint) {
		e.persistence = p
	}
}

// WithTransactions offer the possibility to set a transactions adapter to the endpoint adapter.
func WithTransactions(t transaction.Adapter) func(*EndPoint) {
	return func(e *EndPoint) {
		e.transaction = t
	}
}

// WithWallets offer the possibility to set a wallet adapter to the endpoint adapter.
func WithWallets(w wallet.IWallets) func(*EndPoint) {
	return func(e *EndPoint) {
		e.wallets = w
	}
}

// WithMetrics offer the possibility to set a metric adapter to the endpoint adapter.
func WithMetrics(m *metrics.Metrics) func(*EndPoint) {
	return func(e *EndPoint) {
		e.metrics = m
	}
}

// WithLogs offer the possibility to set a logger to the endpoint adapter.
func WithLogs(logs *zap.Logger) func(*EndPoint) {
	return func(e *EndPoint) {
		e.log = logs.With(zap.String("service", "http"))
	}
}

// WithEvents offer the possibility to set an event adapter to the endpoint adapter.
func WithEvents(evt *event.Event) func(*EndPoint) {
	return func(e *EndPoint) {
		e.event = evt
	}
}

// WithUserAddress  offer the possibility to set an user address adapter to the endpoint adapter.
func WithUserAddress(a string) func(*EndPoint) {
	return func(e *EndPoint) {
		e.userAddress = a
	}
}

// IsEnabled define if we enable HTTP endpoint adapter.
func (e *EndPoint) IsEnabled() bool {
	return e.cfg.Enabled
}

// Listen
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

// MetricsMiddleware set a metrics middleware for HTTP endpoints.
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

// Shutdown the HTTP server.
func (e *EndPoint) Shutdown(ctx context.Context) {
	err := e.server.Shutdown(ctx)
	if err != nil {
		e.log.Error("fail to shutdown server", zap.Error(err))
	}
}
