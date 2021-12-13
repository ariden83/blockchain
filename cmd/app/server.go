package main

import (
	"context"
	"fmt"
	"github.com/ariden83/blockchain/config"
	grpcEndpoint "github.com/ariden83/blockchain/internal/endpoint/grpc"
	httpEndpoint "github.com/ariden83/blockchain/internal/endpoint/http"
	metricsEndpoint "github.com/ariden83/blockchain/internal/endpoint/metrics"
)

type Server struct {
	cfg           *config.Config
	httpServer    *httpEndpoint.EndPoint
	grpcServer    *grpcEndpoint.EndPoint
	metricsServer *metricsEndpoint.EndPoint
}

func (s *Server) Start(stop chan error) {
	if s.cfg.GRPC.Enabled {
		s.startGRPCServer(stop)
	}
	if s.cfg.API.Enabled {
		s.startHTTPServer(stop)
	}
	s.startMetricsServer(stop)
}

// Start Set classic routes
func (s *Server) startHTTPServer(stop chan error) {
	go func() {
		if err := s.httpServer.Listen(); err != nil {
			stop <- fmt.Errorf("cannot start server HTTP %s", err)
		}
	}()
}

// StartGRPC Start GRPC server
func (s *Server) startGRPCServer(stop chan error) {
	go func() {
		if err := s.grpcServer.Listen(); err != nil {
			stop <- err
		}
	}()
}

// StartGRPC Start GRPC server
func (s *Server) startMetricsServer(stop chan error) {
	go func() {
		if err := s.metricsServer.Listen(); err != nil {
			stop <- fmt.Errorf("cannot start healthz server %s", err)
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) {
	s.httpServer.Shutdown(ctx)
	s.grpcServer.Shutdown()
	s.metricsServer.Shutdown(ctx)
}
