package grpc

import (
	"context"
	"github.com/ariden83/blockchain/pkg/api"
)

func (e *EndPoint) GetBalance(_ context.Context, req *api.GetBalanceInput) (*api.GetBalanceOutput, error) {

	balance := e.transaction.FindUserBalance(req.PrivKey)
	tokensSend := e.transaction.FindUserTokensSend(req.PrivKey)
	tokensReceived := e.transaction.FindUserTokensReceived(req.PrivKey)

	return &api.GetBalanceOutput{
		Address:       e.wallets.GetUserAddress(req.PrivKey),
		Balance:       balance.String(),
		TotalReceived: tokensReceived.String(),
		TotalSent:     tokensSend.String(),
	}, nil
}
