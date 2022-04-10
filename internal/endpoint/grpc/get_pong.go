package grpc

import (
	"golang.org/x/net/context"

	protoAPI "github.com/ariden83/blockchain/pkg/api"
)

// GetPong service
func (EndPoint) GetPong(_ context.Context, _ *protoAPI.Ping) (*protoAPI.Pong, error) {
	return &protoAPI.Pong{Message: "pong"}, nil
}
