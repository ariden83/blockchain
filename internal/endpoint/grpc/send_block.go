package grpc

import (
	"context"
	"fmt"
	"math/big"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/internal/transaction"
	"github.com/ariden83/blockchain/pkg/api"
	pkgErr "github.com/ariden83/blockchain/pkg/errors"
)

func (e *EndPoint) SendBlock(_ context.Context, req *api.SendBlockInput) (*api.SendBlockOutput, error) {

	amount := new(big.Int)
	_, err := fmt.Sscan(req.GetAmount(), amount)
	if err != nil {
		return nil, pkgErr.GRPC(err)
	} else if amount.BitLen() == 0 {
		err := pkgErr.ErrMissingFields
		e.log.Error("amount is empty", zap.Error(err))
		return nil, pkgErr.GRPC(err)
	}

	if req.GetFrom() == "" {
		err := pkgErr.ErrMissingFields
		e.log.Error("missing from param", zap.Error(err))
		return nil, pkgErr.GRPC(err)
	}

	if req.GetTo() == "" {
		err := pkgErr.ErrMissingFields
		e.log.Error("missing to param", zap.Error(err))
		return nil, pkgErr.GRPC(err)
	}

	e.transaction.SendBlock(transaction.SendBlockInput{
		From:   []byte(req.GetFrom()),
		To:     []byte(req.GetTo()),
		Amount: amount,
	})

	return &api.SendBlockOutput{}, nil
}
