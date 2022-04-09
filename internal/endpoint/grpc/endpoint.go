package grpc

import (
	"fmt"
	"net"
	"strconv"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	protoAPI "github.com/ariden83/blockchain/pkg/api"

	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/metrics"
	"github.com/ariden83/blockchain/internal/persistence"
	"github.com/ariden83/blockchain/internal/transactions"
	"github.com/ariden83/blockchain/internal/wallet"
	"github.com/grpc-ecosystem/go-grpc-middleware"
)

type EndPoint struct {
	cfg         config.GRPC
	persistence persistence.IPersistence
	transaction transactions.ITransaction
	server      *grpc.Server
	wallets     wallet.IWallets
	metrics     *metrics.Metrics
	log         *zap.Logger
	event       *event.Event
	userAddress string
	ready       bool
	stop        chan error
}

func New(stop chan error, options ...func(*EndPoint)) *EndPoint {
	ep := &EndPoint{
		stop: stop,
	}

	for _, o := range options {
		o(ep)
	}

	return ep
}

func WithConfig(cfg config.GRPC) func(*EndPoint) {
	return func(e *EndPoint) {
		e.cfg = cfg
	}
}

func WithPersistence(p persistence.IPersistence) func(*EndPoint) {
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

// Listen start the server.
func (e *EndPoint) Listen() error {
	address := e.cfg.Host + ":" + strconv.Itoa(e.cfg.Port)

	optsMiddleWare := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			return status.Errorf(codes.Unknown, "panic triggered: %v", p)
		}),
	}

	e.server = grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_opentracing.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(e.log),
			// grpc_auth.StreamServerInterceptor(customFunc),
			grpc_recovery.StreamServerInterceptor(optsMiddleWare...),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_zap.UnaryServerInterceptor(e.log),
			// grpc_auth.StreamServerInterceptor(customFunc),
			grpc_recovery.UnaryServerInterceptor(optsMiddleWare...),
		)),
	)
	//Register the server :
	protoAPI.RegisterApiServer(e.server, e)

	reflection.Register(e.server)
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(e.server)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen to port %s : %s", address, err.Error())
	}

	e.log.Info("Listening GRPC server", zap.String("address", address))
	e.ready = true

	if err := e.server.Serve(lis); err != nil {
		e.ready = false
		e.log.Error("failed to serve : %w", zap.Error(err))
	}

	return nil
}

func (e *EndPoint) Shutdown() {
	e.log.Info("Gracefully pausing down the GRPC server")
	e.server.GracefulStop()
}
