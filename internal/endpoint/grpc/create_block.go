package grpc

import (
	"context"
	pkgErr "github.com/ariden83/blockchain/pkg/errors"
	"go.uber.org/zap"

	"github.com/ariden83/blockchain/pkg/api"
)

func (e *EndPoint) CreateBlock(_ context.Context, req *api.CreateBlockInput) (*api.CreateBlockOutput, error) {
	if req.PubKey == "" {
		err := pkgErr.ErrMissingFields
		e.log.Error("Fail to create block", zap.Error(err))
		return nil, pkgErr.GRPC(err)
	}

	if req.PrivateKey == "" {
		err := pkgErr.ErrMissingFields
		e.log.Error("Fail to create block", zap.Error(err))
		return nil, pkgErr.GRPC(err)
	}

	e.wallets.WriteBlock(*req)
	return &api.CreateBlockOutput{}, nil
}
