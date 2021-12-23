package grpc

import (
	"context"
	"github.com/ariden83/blockchain/internal/wallet"
	"github.com/ariden83/blockchain/pkg/api"
)

func (EndPoint) GetWallet(_ context.Context, req *api.GetWalletInput) (*api.GetWalletOutput, error) {

	keys := wallet.GetKeys(req.GetSeed())

	return &api.GetWalletOutput{
		Address: keys.Address,
		PubKey:  keys.PubKey,
	}, nil
}
