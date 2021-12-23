package grpc

import (
	protoAPI "github.com/ariden83/blockchain/pkg/api"
	"golang.org/x/net/context"
)

// GetPong service
func (EndPoint) GetPong(_ context.Context, in *protoAPI.Ping) (*protoAPI.Pong, error) {
	return &protoAPI.Pong{Message: "pong"}, nil
}
