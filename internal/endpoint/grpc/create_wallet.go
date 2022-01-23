package grpc

import (
	"context"
	"errors"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/pkg/api"
	"go.uber.org/zap"
)

func (e *EndPoint) CreateWallet(_ context.Context, input *api.CreateWalletInput) (*api.CreateWalletOutput, error) {
	if input.Password != "" {
		err := errors.New("missing parameters")
		e.log.Error("Fail to create wallet", zap.Error(err))
		return nil, err
	}
	seed, err := e.wallets.Create([]byte(input.Password))
	if err != nil {
		e.log.Error("Fail to create wallet", zap.Error(err))
		return nil, err
	}

	e.event.Push(event.Message{Type: event.Wallet})

	return &api.CreateWalletOutput{
		Mnemonic: seed.Mnemonic,
		Address:  seed.Address,
		PubKey:   seed.PubKey,
	}, nil
}
