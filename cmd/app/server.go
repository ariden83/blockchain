package main

import (
	"context"
	"fmt"

	grpcEndpoint "github.com/ariden83/blockchain/internal/endpoint/grpc"
	httpEndpoint "github.com/ariden83/blockchain/internal/endpoint/http"
	metricsEndpoint "github.com/ariden83/blockchain/internal/endpoint/metrics"
)

// Server represents all the servers that the service can launch.
type Server struct {
	httpServer    *httpEndpoint.EndPoint
	grpcServer    *grpcEndpoint.EndPoint
	metricsServer *metricsEndpoint.EndPoint
}

// Start HTTP, GRPC and metric servers.
func (s *Server) Start(stop chan error) {
	if s.grpcServer != nil && s.grpcServer.Enabled() {
		s.startGRPCServer(stop)
	}
	if s.httpServer != nil && s.httpServer.IsEnabled() {
		s.startHTTPServer(stop)
	}
	s.startMetricsServer(stop)
}

// startHTTPServer start the HTTP server.
func (s *Server) startHTTPServer(stop chan error) {
	go func() {
		if err := s.httpServer.Listen(); err != nil {
			stop <- fmt.Errorf("cannot start server HTTP %s", err)
		}
	}()
}

// startGRPCServer start the GRPC server.
func (s *Server) startGRPCServer(stop chan error) {
	go func() {
		if err := s.grpcServer.Listen(); err != nil {
			stop <- err
		}
	}()
}

// startMetricsServer start the metrics server.
func (s *Server) startMetricsServer(stop chan error) {
	go func() {
		if err := s.metricsServer.Listen(); err != nil {
			stop <- fmt.Errorf("cannot start healthz server %s", err)
		}
	}()
}

// Shutdown HTTP, GRPC and metric server.
func (s *Server) Shutdown(ctx context.Context) {
	s.httpServer.Shutdown(ctx)
	s.grpcServer.Shutdown()
	s.metricsServer.Shutdown(ctx)
}
