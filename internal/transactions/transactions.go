package transactions

import (
	"errors"
	"math/big"
	"sync"
	"time"

	"encoding/hex"
	"github.com/gcash/bchd/bchec"
	"go.uber.org/zap"

	"github.com/ariden83/blockchain/config"
	"github.com/ariden83/blockchain/internal/blockchain"
	signschnorr "github.com/ariden83/blockchain/internal/blockchain/signSchnorr"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/iterator"
	"github.com/ariden83/blockchain/internal/persistence"
	"github.com/ariden83/blockchain/internal/wallet"
	pkgError "github.com/ariden83/blockchain/pkg/errors"
)

var ErrNotEnoughFunds = errors.New("Not enough funds")

type Transactions struct {
	Reward          *big.Int
	serverPublicKey string
	persistence     persistence.IPersistence
	event           *event.Event
	log             *zap.Logger
}

type ITransaction interface {
	New([]byte, []byte, []byte, *big.Int) (*blockchain.Transaction, error)
	CoinBaseTx([]byte) *blockchain.Transaction
	FindUserBalance([]byte) *big.Int
	FindUserTokensSend([]byte) *big.Int
	FindUserTokensReceived([]byte) *big.Int
	WriteBlock([]byte) (*blockchain.Block, error)
	GetLastBlock() ([]byte, *big.Int, error)
	SendBlock(input SendBlockInput) error
}

var mutex = &sync.Mutex{}

func New(options ...func(*Transactions)) *Transactions {
	t := &Transactions{}

	for _, o := range options {
		o(t)
	}

	return t
}

func WithConfig(cfg config.Transactions) func(*Transactions) {
	return func(e *Transactions) {
		e.Reward = cfg.Reward
	}
}

func WithPersistence(p persistence.IPersistence) func(*Transactions) {
	return func(e *Transactions) {
		e.persistence = p
	}
}

func WithEvents(evt *event.Event) func(*Transactions) {
	return func(e *Transactions) {
		e.event = evt
	}
}

func WithLogs(logs *zap.Logger) func(*Transactions) {
	return func(e *Transactions) {
		e.log = logs.With(zap.String("service", "transactions"))
	}
}

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

// new transaction
// from privkey
// to publickey
func (t *Transactions) New(pubKey, privKey, to []byte, amount *big.Int) (*blockchain.Transaction, error) {
	var (
		inputs  []blockchain.TxInput
		outputs []blockchain.TxOutput
	)

	if !t.canPayTransactionFees(amount) {
		return nil, pkgError.ErrNotEnoughFunds
	}

	// Trouver des sorties utilisables
	acc, validOutputs := t.FindSpendableOutputs(pubKey, amount)

	// Vérifiez si nous avons assez d'argent pour envoyer le montant que nous demandons
	if acc.Cmp(amount) < 0 {
		return nil, pkgError.ErrNotEnoughFunds
	}

	pubKeySchnorr, sig, err := t.getSchnorrKeys(pubKey, privKey)
	if err != nil {
		return nil, pkgError.ErrInternalError
	}

	// Si nous le faisons, créez des inputs qui indiquent les outputs que nous dépensons
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

	// on récupère la valeur des frais pour le mineur et on les applique
	amountLessFees := t.setTransactionFees(amount)
	outputs = append(outputs, blockchain.TxOutput{Value: amountLessFees, PubKey: to})
	outputs = t.payTransactionFees(outputs)

	// S'il reste de l'argent, faites de nouvelles sorties à partir de la différence.
	if acc.Cmp(amount) > 0 {
		newAmount := acc.Sub(acc, amount)
		outputs = append(outputs, blockchain.TxOutput{Value: newAmount, PubKey: pubKey})
	}

	// Initialiser une nouvelle transaction avec toutes les nouvelles entrées et sorties que nous avons effectuées
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	tx := blockchain.Transaction{nil, inputs, outputs, timestamp}

	// Définissez un nouvel identifiant et renvoyez-le.
	tx.SetID()

	return &tx, nil
}

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

// Maintenant que nous avons un moyen de trouver toutes les transactions non dépensées,
// nous pouvons rechercher les sorties non dépensées.
// Parcourez toutes les transactions non dépensées et voyez si nous pouvons déverrouiller les sorties.
// Ajoutez-les tous à un tableau et retournez ce tableau de TxOutputs
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

func (t *Transactions) FindUserBalance(pubKey []byte) *big.Int {
	var balance *big.Int = new(big.Int)
	UTXOs := t.FindUTXO(pubKey)

	for _, out := range UTXOs {
		balance = balance.Add(balance, out.Value)
	}
	return balance
}

func (t *Transactions) FindUserTokensSend(pubKey []byte) *big.Int {
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

func (t *Transactions) FindUserTokensReceived(pubKey []byte) *big.Int {
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
func (t *Transactions) FindSpendableOutputs(pubKey []byte, amount *big.Int) (*big.Int, map[string][]int) {
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
