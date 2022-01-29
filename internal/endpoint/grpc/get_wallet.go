package grpc

import (
	"context"
	"errors"
	"github.com/ariden83/blockchain/pkg/api"
	"go.uber.org/zap"
)

func (e *EndPoint) GetWallet(_ context.Context, req *api.GetWalletInput) (*api.GetWalletOutput, error) {
	if req.GetMnemonic() == nil || req.GetPassword() == nil {
		err := errors.New("missing fields")
		e.log.Error("fail to get wallet", zap.Error(err))
		return nil, err
	}

	seed, err := e.wallets.GetSeed(req.GetMnemonic(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &api.GetWalletOutput{
		Address:  seed.Address,
		PubKey:   seed.PubKey,
		Mnemonic: seed.Mnemonic,
	}, nil
}
