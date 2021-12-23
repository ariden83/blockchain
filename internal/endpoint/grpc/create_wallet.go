package grpc

import (
	"context"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/pkg/api"
	"go.uber.org/zap"
	"time"
)

func (e *EndPoint) CreateWallet(_ context.Context, _ *api.CreateWalletInput) (*api.CreateWalletOutput, error) {
	newSeed, err := e.wallets.Create()
	if err != nil {
		e.log.Error("Fail to create wallet", zap.Error(err))
		return nil, err
	}

	e.event.Push(event.Message{Type: event.Wallet})

	return &api.CreateWalletOutput{
		Address:   newSeed.Address,
		Timestamp: time.Unix(newSeed.Timestamp, 0).String(),
		PubKey:    newSeed.PubKey,
		Mnemonic:  newSeed.Mnemonic,
	}, nil
}
