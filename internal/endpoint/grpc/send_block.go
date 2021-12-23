package grpc

import (
	"context"
	"github.com/ariden83/blockchain/pkg/api"
)

func (e *EndPoint) SendBlock(_ context.Context, req *api.Ping) (*api.Pong, error) {
	return &api.Pong{}, nil
}
