package transaction

import (
	"encoding/hex"
	"math/big"
	"sync"
	"time"

	"github.com/gcash/bchd/bchec"
	"go.uber.org/zap"

	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/blockchain"
	signschnorr "github.com/ariden83/blockchain/internal/blockchain/signSchnorr"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/iterator"
	"github.com/ariden83/blockchain/internal/wallet"
	pkgError "github.com/ariden83/blockchain/pkg/errors"
)

// Transactions represent a transactions adapter.
type Transactions struct {
	Reward          *big.Int
	serverPublicKey string
	persistence     persistenceadapter.Adapter
	event           *event.Event
	log             *zap.Logger
}

var mutex = &sync.Mutex{}

// New define a new transaction adapter.
func New(options ...func(*Transactions)) *Transactions {
	t := &Transactions{}

	for _, o := range options {
		o(t)
	}

	return t
}

// WithConfig
func WithConfig(cfg config.Transactions) func(*Transactions) {
	return func(e *Transactions) {
		e.Reward = cfg.Reward
	}
}

// WithPersistence
func WithPersistence(p persistenceadapter.Adapter) func(*Transactions) {
	return func(e *Transactions) {
		e.persistence = p
	}
}

// WithEvents
func WithEvents(evt *event.Event) func(*Transactions) {
	return func(e *Transactions) {
		e.event = evt
	}
}

// WithLogs
func WithLogs(logs *zap.Logger) func(*Transactions) {
	return func(e *Transactions) {
		e.log = logs.With(zap.String("service", "transactions"))
	}
}

// getSchnorrKeys retrieve schnorr keys for a user.
//
// @see https://tlu.tarilabs.com/cryptography/introduction-schnorr-signatures
func (t *Transactions) getSchnorrKeys(pubKey, privKey []byte) ([]byte, []byte, error) {
	priv, pubKeySchnorr := bchec.PrivKeyFromBytes(bchec.S256(), privKey)

	sig, err := signschnorr.SignSchnorr(priv, pubKey)
	if err != nil {
		t.log.Error("fail to sign schnorr", zap.Error(err))
		return nil, nil, err
	}

	return pubKeySchnorr.SerializeCompressed(), sig.Serialize(), nil
}

// CoinBaseTx is the function that will run when someone on a node succesfully "mines" a block. The reward inside as it were.
func (t *Transactions) CoinBaseTx(privKey []byte) *blockchain.Transaction {
	pubKey, err := wallet.GetPubKey(privKey)
	if err != nil {
		t.log.Error("fail to get pub key from private key", zap.Error(err))
		return nil
	}

	pubKeySchnorr, sig, err := t.getSchnorrKeys(pubKey, privKey)
	if err != nil {
		return nil
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
		Sig:        sig,
		PubKey:     pubKey,
		SchnorrKey: pubKeySchnorr,
	}
	// txOut will represent the amount of tokens(reward) given to the person(toAddress) that executed CoinbaseTx

	// Value would be representative of the amount of coins in a transaction
	txOut := blockchain.TxOutput{
		// Value would be representative of the amount of coins in a transaction
		Value: t.Reward,
		// La Pubkey est nécessaire pour "déverrouiller" toutes les pièces dans une sortie. Cela indique que VOUS êtes celui qui l'a envoyé.
		// Vous êtes identifiable par votre PubKey
		// PubKey dans cette itération sera très simple, mais dans une application réelle, il s'agit d'un algorithme plus complexe
		PubKey: pubKey,
	} // You can see it follows {value, PubKey}

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	tx := blockchain.Transaction{
		Inputs:    []blockchain.TxInput{txIn},
		Outputs:   []blockchain.TxOutput{txOut},
		Timestamp: timestamp,
	}

	return &tx
}

/*
func (t *Transactions) privKeyToPublicKey(privKey string) (string, error){
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return "", err
	}

	// use TestNet3Params for interacting with bitcoin testnet
	// if we want to interact with main net should use MainNetParams
	addrPubKey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}
	return addrPubKey.EncodeAddress(), nil
}*/

// New instantiate a new transaction.
// from privkey
// to publickey
func (t *Transactions) New(from, to []byte, amount *big.Int) (*blockchain.Transaction, error) {
	var (
		inputs  []blockchain.TxInput
		outputs []blockchain.TxOutput
	)

	if !t.canPayTransactionFees(amount) {
		return nil, pkgError.ErrNotEnoughFunds
	}

	pubKey, err := wallet.GetPubKey(from)
	if err != nil {
		return nil, pkgError.ErrInternalError
	}

	// Find usable outputs.
	acc, validOutputs := t.FindSpendableOutputs(pubKey, amount)

	// Check if we have enough money to send the amount we request.
	if acc.Cmp(amount) < 0 {
		return nil, pkgError.ErrNotEnoughFunds
	}

	pubKeySchnorr, sig, err := t.getSchnorrKeys(pubKey, from)
	if err != nil {
		return nil, pkgError.ErrInternalError
	}

	// If we do, create inputs that indicate the outputs we are spending.
	for txID, outs := range validOutputs {
		decodeTxID, err := hex.DecodeString(txID)
		if err != nil {
			return nil, err
		}

		for _, out := range outs {
			input := blockchain.TxInput{
				ID:         decodeTxID,
				Out:        out,
				Sig:        sig,
				PubKey:     pubKey,
				SchnorrKey: pubKeySchnorr,
			}
			inputs = append(inputs, input)
		}
	}

	// we recover the value of the fees for the minor and we apply them.
	amountLessFees := t.setTransactionFees(amount)
	outputs = append(outputs, blockchain.TxOutput{Value: amountLessFees, PubKey: to})
	outputs = t.payTransactionFees(outputs)

	// If there is money left, make new exits from the difference.
	if acc.Cmp(amount) > 0 {
		newAmount := acc.Sub(acc, amount)
		outputs = append(outputs, blockchain.TxOutput{Value: newAmount, PubKey: pubKey})
	}

	// Initialize a new transaction with all the new inputs and outputs we have made.
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	tx := blockchain.Transaction{nil, inputs, outputs, timestamp}

	// Set a new id and send it back.
	tx.SetID()

	return &tx, nil
}

// FindUnspentTransactions allows to find all unspent transactions linked to a user.
func (t *Transactions) FindUnspentTransactions(pubKey []byte) []blockchain.Transaction {
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

// FindUTXO allow to find all unspent transactions,
// we can search for unspent outputs.
// Go through all the unspent transactions and see if we can unlock the outputs.
// Add them all to an array and return that array of TxOutputs.
func (t *Transactions) FindUTXO(pubKey []byte) []blockchain.TxOutput {
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

/*
type transaction struct {

}

func (t *Transactions) FindAllUserTransaction(pubKey string) []blockchain.TxOutput {
	listTransaction := []transaction{}

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
					listTransaction = append(listTransaction, transaction{
						out.
					})
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
*/

// FindUserBalance retrieve a user balance.
func (t *Transactions) FindUserBalance(pubKey []byte) *big.Int {
	balance := new(big.Int)
	UTXOs := t.FindUTXO(pubKey)

	for _, out := range UTXOs {
		balance = balance.Add(balance, out.Value)
	}
	return balance
}

// FindUserTokensSend allows to find the tokens sent by a user.
func (t *Transactions) FindUserTokensSend(pubKey []byte) *big.Int {
	tokensSend := new(big.Int)

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

// FindUserTokensReceived allows to find the tokens received by a user.
func (t *Transactions) FindUserTokensReceived(pubKey []byte) *big.Int {
	tokensReceived := new(big.Int)

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

// FindSpendableOutputs takes the address of an account and the amount we would like to spend.
// It returns a tuple that contains the amount we can spend and a map of the aggregate outputs that can get there.
func (t *Transactions) FindSpendableOutputs(pubKey []byte, amount *big.Int) (*big.Int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := t.FindUnspentTransactions(pubKey)
	accumulated := new(big.Int)

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
