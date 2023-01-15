package transaction

import (
	"math/big"

	"github.com/ariden83/blockchain/internal/blockchain"
)

// Adapter is an interface which describes all the methods to interact with persistence.
type Adapter interface {
	New([]byte, []byte, *big.Int) (*blockchain.Transaction, error)
	CoinBaseTx([]byte) *blockchain.Transaction
	FindUserBalance([]byte) *big.Int
	FindUserTokensSend([]byte) *big.Int
	FindUserTokensReceived([]byte) *big.Int
	WriteBlock([]byte) (*blockchain.Block, error)
	GetLastBlock() ([]byte, *big.Int, error)
	SendBlock(input SendBlockInput) error
}

type SendBlockInput struct {
	From   []byte
	To     []byte
	Amount *big.Int
}
