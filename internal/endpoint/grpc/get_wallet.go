package grpc

import (
	"context"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/pkg/api"
	pkgErr "github.com/ariden83/blockchain/pkg/errors"
)

// GetWallet
func (e *EndPoint) GetWallet(_ context.Context, req *api.GetWalletInput) (*api.GetWalletOutput, error) {
	if req.GetMnemonic() == nil || req.GetPassword() == nil {
		err := pkgErr.ErrMissingFields
		e.log.Error(err.Error(), zap.String("mnemonic", string(req.Mnemonic)))
		return nil, pkgErr.GRPC(err)
	}

	seed, err := e.wallets.Seed(req.GetMnemonic(), req.GetPassword())
	if err != nil {
		return nil, pkgErr.GRPC(err)
	}

	return &api.GetWalletOutput{
		Address: seed.Address,
		PubKey:  seed.PubKey,
		PrivKey: seed.PrivKey,
	}, nil
}
