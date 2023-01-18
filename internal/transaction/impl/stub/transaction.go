package stub

import (
	"math/big"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/persistence"
	transactionadapter "github.com/ariden83/blockchain/internal/transaction"
)

// Transactions represent a transactions adapter.
type Transactions struct {
	Reward          *big.Int
	serverPublicKey string
	persistence     persistenceadapter.Adapter
	event           *event.Event
	log             *zap.Logger
}

func New() *Transactions {
	return &Transactions{}
}

func (p *Transactions) New([]byte, []byte, *big.Int) (*blockchain.Transaction, error) {
	return &blockchain.Transaction{}, nil
}

func (p *Transactions) CoinBaseTx([]byte) *blockchain.Transaction {
	return &blockchain.Transaction{}
}

func (p *Transactions) FindUserBalance([]byte) *big.Int {
	balance := new(big.Int)
	balance.SetInt64(1)
	return balance
}

func (p *Transactions) FindUserTokensSend([]byte) *big.Int {
	balance := new(big.Int)
	balance.SetInt64(1)
	return balance
}

func (p *Transactions) FindUserTokensReceived([]byte) *big.Int {
	balance := new(big.Int)
	balance.SetInt64(1)
	return balance
}

func (p *Transactions) WriteBlock([]byte) (*blockchain.Block, error) {
	return &blockchain.Block{}, nil
}

func (p *Transactions) GetLastBlock() ([]byte, *big.Int, error) {
	balance := new(big.Int)
	balance.SetInt64(1)
	return []byte("last block"), balance, nil
}

func (p *Transactions) SendBlock(input transactionadapter.SendBlockInput) error {
	return nil
}
