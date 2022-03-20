package grpc

import (
	"context"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/pkg/api"
	pkgErr "github.com/ariden83/blockchain/pkg/errors"
)

func (e *EndPoint) CreateBlock(_ context.Context, req *api.CreateBlockInput) (*api.CreateBlockOutput, error) {
	if req.GetPrivateKey() == "" {
		err := pkgErr.ErrMissingFields
		e.log.Error("Empty private key", zap.Error(err))
		return nil, pkgErr.GRPC(err)
	}

	_, err := e.transaction.WriteBlock([]byte(req.GetPrivateKey()))
	if err != nil {
		return nil, pkgErr.GRPC(err)
	}

	return &api.CreateBlockOutput{}, nil
}
