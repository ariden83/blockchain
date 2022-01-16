package grpc

import (
	"context"
	"github.com/ariden83/blockchain/pkg/api"
)

func (e *EndPoint) GetBalance(_ context.Context, req *api.GetBalanceInput) (*api.GetBalanceOutput, error) {
	balance := e.transaction.FindUserBalance(req.PubKey)
	tokensSend := e.transaction.FindUserTokensSend(req.PubKey)
	tokensReceived := e.transaction.FindUserTokensReceived(req.PubKey)

	return &api.GetBalanceOutput{
		Address:       req.PubKey,
		Balance:       balance.String(),
		TotalReceived: tokensReceived.String(),
		TotalSent:     tokensSend.String(),
	}, nil
}
