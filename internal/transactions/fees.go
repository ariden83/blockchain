package transactions

import (
	"github.com/ariden83/blockchain/internal/blockchain"
	"go.uber.org/zap"
	"math/big"
)

func (t *Transactions) canPayTransactionFees(amount *big.Int) bool {
	if amount.Cmp(t.Reward) <= 0 {
		t.log.Info("amount to send less than rewards, must send more token", zap.Any("reward", t.Reward), zap.Any("amount", amount))
		return false
	}
	return true
}

func (t *Transactions) setTransactionFees(amount *big.Int) *big.Int {
	newAmount := amount.Sub(amount, t.Reward)
	return newAmount
}

func (t *Transactions) payTransactionFees(outputs []blockchain.TxOutput) []blockchain.TxOutput {
	return append(outputs, blockchain.TxOutput{t.Reward, t.serverPublicKey})
}
