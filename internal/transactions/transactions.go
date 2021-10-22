package transactions

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/iterator"
	"github.com/ariden83/blockchain/internal/persistence"
	"go.uber.org/zap"
	"math/big"
	"time"
)

var ErrNotEnoughFunds = errors.New("Not enough funds")

type Transactions struct {
	Reward      *big.Int
	persistence persistence.IPersistence
	log         *zap.Logger
}

type ITransaction interface {
	New(string, string, *big.Int) (*blockchain.Transaction, error)
	CoinBaseTx(string, string) *blockchain.Transaction
	FindUserBalance(string) *big.Int
	FindUserTokensSend(string) *big.Int
	FindUserTokensReceived(pubKey string) *big.Int
}

func Init(conf config.Transactions, per persistence.IPersistence, log *zap.Logger) *Transactions {
	return &Transactions{
		Reward:      conf.Reward,
		persistence: per,
		log:         log.With(zap.String("service", "transactions")),
	}
}

// CoinBaseTx is the function that will run when someone on a node succesfully "mines" a block. The reward inside as it were.
func (t *Transactions) CoinBaseTx(toPubKey, data string) *blockchain.Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", toPubKey)
	}
	// Since this is the "first" transaction of the block, it has no previous output to reference.
	// This means that we initialize it with no ID, and it's OutputIndex is -1
	txIn := blockchain.TxInput{
		// ID will find the Transaction that a specific output is inside of
		ID: []byte{},
		// Out will be the index of the specific output we found within a transaction.
		// For example if a transaction has 4 outputs, we can use this "out" field to specify which output we are looking for
		Out: -1,
		// This would be a script that adds data to an outputs' PubKey
		// however for this tutorial the Sig will be indentical to the PubKey.
		Sig: data,
	}
	// txOut will represent the amount of tokens(reward) given to the person(toAddress) that executed CoinbaseTx

	// Value would be representative of the amount of coins in a transaction
	txOut := blockchain.TxOutput{
		// Value would be representative of the amount of coins in a transaction
		Value: t.Reward,
		// La Pubkey est nécessaire pour "déverrouiller" toutes les pièces dans une sortie. Cela indique que VOUS êtes celui qui l'a envoyé.
		// Vous êtes identifiable par votre PubKey
		// PubKey dans cette itération sera très simple, mais dans une application réelle, il s'agit d'un algorithme plus complexe
		PubKey: toPubKey,
	} // You can see it follows {value, PubKey}

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	tx := blockchain.Transaction{nil, []blockchain.TxInput{txIn}, []blockchain.TxOutput{txOut}, timestamp}

	return &tx

}

func (t *Transactions) New(from, to string, amount *big.Int) (*blockchain.Transaction, error) {
	var inputs []blockchain.TxInput
	var outputs []blockchain.TxOutput

	// Trouver des sorties utilisables
	acc, validOutputs := t.FindSpendableOutputs(from, amount)

	// Vérifiez si nous avons assez d'argent pour envoyer le montant que nous demandons
	if acc.Cmp(amount) < 0 {
		return nil, ErrNotEnoughFunds
	}
	// Si nous le faisons, créez des inputs qui indiquent les outputs que nous dépensons
	for txID, outs := range validOutputs {
		decodeTxID, err := hex.DecodeString(txID)
		if err != nil {
			return nil, err
		}

		for _, out := range outs {
			input := blockchain.TxInput{decodeTxID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, blockchain.TxOutput{amount, to})

	// S'il reste de l'argent, faites de nouvelles sorties à partir de la différence.
	if acc.Cmp(amount) > 0 {
		newAmount := acc.Sub(acc, amount)
		outputs = append(outputs, blockchain.TxOutput{newAmount, from})
	}

	// Initialiser une nouvelle transaction avec toutes les nouvelles entrées et sorties que nous avons effectuées
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	tx := blockchain.Transaction{nil, inputs, outputs, timestamp}

	// Définissez un nouvel identifiant et renvoyez-le.
	tx.SetID()

	return &tx, nil
}

// un TxOutput est représentatif d'une action entre deux adresses,
// TxInput est simplement une référence à un TxOutput précédent.
// nous sommes en mesure de déterminer le solde d'un compte.
// nous pouvons vérifier toutes les sorties qu'un compte lui a liées, puis vérifier toutes les entrées.
// Les sorties qui n'ont pas d'entrée pointant vers elles seront dépensables.
/**
Transactions: ([]*blockchain.Transaction) (len=1 cap=1) {
	(*blockchain.Transaction)(0xc00556c230)({
		ID: ([]uint8) <nil>,
		Inputs: ([]blockchain.TxInput) (len=1 cap=1) {
			(blockchain.TxInput) {
				ID: ([]uint8) {},
				Out: (int) -1,
				Sig: (string) (len=111) "xpub661MyMwAqRbcFTZYiEcSv4Qj2Qr2NzQ7rjYc3iv9c6VSTxoYsqA9AA6nNbp8e9nVR9hRARXz5CApP6j5BxUnohyj89oSg3zZdDuKmGhdSFF"
			}
		},
		Outputs: ([]blockchain.TxOutput) (len=1 cap=1) {
			(blockchain.TxOutput) {
				Value: (int) 100,
				PubKey: (string) (len=34) "1P1aBegXRiTinJhhEYHHiMALfG26Wu9sG3"
			}
		}
	})
}
*/
func (t *Transactions) FindUnspentTransactions(pubKey string) []blockchain.Transaction {
	var unspentTxs []blockchain.Transaction

	spentTXOs := make(map[string][]int)

	iter := iterator.BlockChainIterator{
		CurrentHash: t.persistence.LastHash(),
		Persistence: t.persistence,
	}

	// pour chaque bloc
	for {
		block, err := iter.Next()
		if err != nil {
			t.log.Fatal("fail to iterate next block", zap.Error(err))
		}

		// pour chaque transaction
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		Outputs:
			// pour chaque noeud (pub key) de chaque transaction
			for outIdx, out := range tx.Outputs {
				// Si, nous trouvons un txID (ID de transaction)
				// dans toutes les transactions déjà trouvées pour le pub key actuel,
				// nous savons que cette sortie a été dépensée plus tard dans le temps et doit être ignorée
				// (rappel, nous sommes ici dans une boucle inversée, on remonte le temps)
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				// Nous vérifions si le pub key en cours correspond au pubkey de la transaction
				if out.CanBeUnlocked(pubKey) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			// nous vérifions si la transaction est une transaction coinbase (donc nouveau block).
			// Si ce n'est pas le cas, nous pouvons parcourir toutes l'historique de transaction lié au block en cours.
			if !tx.IsCoinBase() {
				// Pour chaque historique de transaction
				for _, in := range tx.Inputs {
					// si le sig enregistré correspond à notre pubkey
					if in.CanUnlock(pubKey) {
						inTxID := hex.EncodeToString(in.ID)
						// alors nous enregistrons l'ID de transaction pour les ignorer immédiatement par la suite
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

// Maintenant que nous avons un moyen de trouver toutes les transactions non dépensées,
// nous pouvons rechercher les sorties non dépensées.
// Parcourez toutes les transactions non dépensées et voyez si nous pouvons déverrouiller les sorties.
// Ajoutez-les tous à un tableau et retournez ce tableau de TxOutputs
func (t *Transactions) FindUTXO(pubKey string) []blockchain.TxOutput {
	var UTXOs []blockchain.TxOutput
	unspentTransactions := t.FindUnspentTransactions(pubKey)
	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(pubKey) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (t *Transactions) FindUserBalance(pubKey string) *big.Int {
	var balance *big.Int = new(big.Int)
	UTXOs := t.FindUTXO(pubKey)

	for _, out := range UTXOs {
		balance = balance.Add(balance, out.Value)
	}
	return balance
}

func (t *Transactions) FindUserTokensSend(pubKey string) *big.Int {
	var tokensSend *big.Int = new(big.Int)

	iter := iterator.BlockChainIterator{
		CurrentHash: t.persistence.LastHash(),
		Persistence: t.persistence,
	}

	for {
		block, err := iter.Next()
		if err != nil {
			t.log.Fatal("fail to iterate next block", zap.Error(err))
		}
		for _, tx := range block.Transactions {

			var isSending bool
			for _, in := range tx.Inputs {
				// si le sig enregistré correspond à notre pubkey
				if in.CanUnlock(pubKey) {
					isSending = true
				}
			}

			if !isSending {
				continue
			}

			for _, out := range tx.Outputs {
				if !out.CanBeUnlocked(pubKey) && isSending {
					tokensSend = tokensSend.Add(tokensSend, out.Value)
				}
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}

	return tokensSend
}

func (t *Transactions) FindUserTokensReceived(pubKey string) *big.Int {
	var tokensReceived *big.Int = new(big.Int)

	iter := iterator.BlockChainIterator{
		CurrentHash: t.persistence.LastHash(),
		Persistence: t.persistence,
	}

	for {
		block, err := iter.Next()
		if err != nil {
			t.log.Fatal("fail to iterate next block", zap.Error(err))
		}
	Outputs:
		for _, tx := range block.Transactions {
			if tx.IsCoinBase() {
				continue
			}
			for _, in := range tx.Inputs {
				if in.CanUnlock(pubKey) {
					continue Outputs
				}
			}
			for _, out := range tx.Outputs {
				if out.CanBeUnlocked(pubKey) {
					tokensReceived = tokensReceived.Add(tokensReceived, out.Value)
				}
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}

	return tokensReceived
}

// Ce qui suit sera une fonction qui prend l'adresse d'un compte et le montant que nous aimerions dépenser.
// Il renvoie un tuple qui contient le montant que nous pouvons dépenser et une carte des sorties agrégées qui peuvent y arriver.
func (t *Transactions) FindSpendableOutputs(pubKey string, amount *big.Int) (*big.Int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := t.FindUnspentTransactions(pubKey)
	var accumulated *big.Int = new(big.Int)

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(pubKey) && accumulated.Cmp(amount) < 0 {

				accumulated = accumulated.Add(accumulated, out.Value)
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated.Cmp(amount) >= 0 {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOuts
}
