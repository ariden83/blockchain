package grpc

import (
	"fmt"
	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/metrics"
	"github.com/ariden83/blockchain/internal/persistence"
	"github.com/ariden83/blockchain/internal/transactions"
	"github.com/ariden83/blockchain/internal/wallet"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-prometheus"

	protoAPI "github.com/ariden83/blockchain/pkg/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
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
}

func New(
	cfg config.GRPC,
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
		log:         logs.With(zap.String("service", "grpc")),
		event:       evt,
		userAddress: userAddress,
	}

	return e
}

// Listen start the server.
func (s *EndPoint) Listen() error {
	address := fmt.Sprintf(":%d", s.cfg.Port)

	s.server = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
		)),
	)
	//Register the server :
	protoAPI.RegisterApiServer(s.server, s)

	reflection.Register(s.server)
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(s.server)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen to port %s : %s", address, err.Error())
	}

	s.log.Info("Listening GRPC server", zap.String("address", address))
	s.ready = true

	if err := s.server.Serve(lis); err != nil {
		s.ready = false
		s.log.Error("failed to serve : %s", zap.Error(err))
	}

	return nil
}

func (s *EndPoint) Shutdown() {
	s.log.Debug("Gracefully pausing down the GRPC server")
	s.server.GracefulStop()
}
