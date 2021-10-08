package transactions

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/iterator"
	"github.com/ariden83/blockchain/internal/persistence"
)

//reward is the amnount of tokens given to someone that "mines" a new block
const reward = 100

var ErrNotEnoughFunds = errors.New("Not enough funds")

type Transactions struct {
	config      *config.Config
	persistence *persistence.Persistence
}

func Init(conf *config.Config, per *persistence.Persistence) *Transactions {
	return &Transactions{
		config:      conf,
		persistence: per,
	}
}

//CoinbaseTx is the function that will run when someone on a node succesfully "mines" a block. The reward inside as it were.
func CoinbaseTx(toAddress, publicKey string) *blockchain.Transaction {
	if publicKey == "" {
		publicKey = fmt.Sprintf("Coins to %s", toAddress)
	}
	//Since this is the "first" transaction of the block, it has no previous output to reference.
	//This means that we initialize it with no ID, and it's OutputIndex is -1
	txIn := blockchain.TxInput{[]byte{}, -1, publicKey}
	//txOut will represent the amount of tokens(reward) given to the person(toAddress) that executed CoinbaseTx
	txOut := blockchain.TxOutput{reward, toAddress} // You can see it follows {value, PubKey}

	tx := blockchain.Transaction{nil, []blockchain.TxInput{txIn}, []blockchain.TxOutput{txOut}}

	return &tx

}

func (t *Transactions) New(from, to string, amount int) (*blockchain.Transaction, error) {
	var inputs []blockchain.TxInput
	var outputs []blockchain.TxOutput

	acc, validOutputs := t.FindSpendableOutputs(from, amount)

	if acc < amount {
		return nil, ErrNotEnoughFunds
	}
	for txID, outs := range validOutputs {
		decodeTxID, err := hex.DecodeString(txID)
		return nil, err

		for _, out := range outs {
			input := blockchain.TxInput{decodeTxID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, blockchain.TxOutput{amount, to})

	if acc > amount {
		outputs = append(outputs, blockchain.TxOutput{acc - amount, from})
	}

	tx := blockchain.Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx, nil
}

func (t *Transactions) FindUnspentTransactions(address string) []blockchain.Transaction {
	var unspentTxs []blockchain.Transaction

	spentTXOs := make(map[string][]int)

	iter := iterator.BlockChainIterator{
		CurrentHash: t.persistence.LastHash,
		Persistence: t.persistence,
	}

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			if !tx.IsCoinBase() {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}
	return unspentTxs
}

func (t *Transactions) FindUTXO(address string) []blockchain.TxOutput {
	var UTXOs []blockchain.TxOutput
	unspentTransactions := t.FindUnspentTransactions(address)
	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (t *Transactions) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := t.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOuts
}
