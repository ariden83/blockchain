package grpc

import (
	"context"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/pkg/api"
	pkgErr "github.com/ariden83/blockchain/pkg/errors"
)

func (e *EndPoint) CreateWallet(_ context.Context, input *api.CreateWalletInput) (*api.CreateWalletOutput, error) {
	if input.GetPassword() == nil {
		err := pkgErr.ErrMissingPassword
		e.log.Error("Fail to create wallet", zap.Error(err))
		return nil, pkgErr.GRPC(err)
	}

	seed, err := e.wallets.Create(input.Password)
	if err != nil {
		e.log.Error("Fail to create wallet", zap.Error(err))
		return nil, pkgErr.GRPC(err)
	}

	e.event.Push(event.Message{Type: event.Wallet})

	return &api.CreateWalletOutput{
		Mnemonic: seed.Mnemonic,
		Address:  seed.Address,
		PubKey:   seed.PubKey,
		PrivKey:  seed.PrivKey,
	}, nil
}
