package grpc

import (
	"context"
	"github.com/ariden83/blockchain/internal/transactions"
	pkgErr "github.com/ariden83/blockchain/pkg/errors"
	"go.uber.org/zap"

	"github.com/ariden83/blockchain/pkg/api"
)

func (e *EndPoint) CreateBlock(_ context.Context, req *api.CreateBlockInput) (*api.CreateBlockOutput, error) {
	if req.GetPubKey() == "" {
		err := pkgErr.ErrMissingFields
		e.log.Error("Empty pub key", zap.Error(err))
		return nil, pkgErr.GRPC(err)
	}

	if req.GetPrivateKey() == "" {
		err := pkgErr.ErrMissingFields
		e.log.Error("Empty private key", zap.Error(err))
		return nil, pkgErr.GRPC(err)
	}

	_, err := e.transaction.WriteBlock(
		transactions.WriteBlockInput{
			PubKey:     req.GetPubKey(),
			PrivateKey: req.GetPrivateKey(),
		})
	if err != nil {
		return nil, pkgErr.GRPC(err)
	}

	return &api.CreateBlockOutput{}, nil
}
