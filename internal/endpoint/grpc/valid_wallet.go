package grpc

import (
	"context"
	"github.com/ariden83/blockchain/pkg/api"
	pkgErr "github.com/ariden83/blockchain/pkg/errors"
)

func (e *EndPoint) ValidWallet(_ context.Context, req *api.ValidWalletInput) (*api.ValidWalletOutput, error) {
	if req.GetPubKey() == nil {
		err := pkgErr.ErrMissingFields
		e.log.Error(err.Error())
		return nil, pkgErr.GRPC(err)
	}

	valid := e.wallets.Validate(req.GetPubKey())
	if !valid {
		err := pkgErr.ErrSeedNotFound
		e.log.Error(err.Error())
		return nil, pkgErr.GRPC(err)
	}

	return &api.ValidWalletOutput{
		Valid: valid,
	}, nil
}
