package grpc

import (
	"context"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/pkg/api"
	pkgErr "github.com/ariden83/blockchain/pkg/errors"
)

func (e *EndPoint) CreateBlock(_ context.Context, req *api.CreateBlockInput) (*api.CreateBlockOutput, error) {
	if req.GetPrivKey() == nil {
		err := pkgErr.ErrMissingFields
		e.log.Error("Empty private key", zap.Error(err))
		return nil, pkgErr.GRPC(err)
	}

	go func() {
		_, err := e.transaction.WriteBlock(req.GetPrivKey())
		if err != nil {
			e.log.Error("invalid block", zap.Error(err))
		}
	}()

	return &api.CreateBlockOutput{}, nil
}
