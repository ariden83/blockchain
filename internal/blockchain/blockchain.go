package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/ariden83/blockchain/internal/blockchain/difficulty"
)

type Validation struct {
	Total   int
	Refused int
	Ok      int
}

type Block struct {
	Index        *big.Int
	Validation   Validation
	Timestamp    int64
	Transactions []*Transaction
	Hash         []byte
	PrevHash     []byte
	// Nœud racine du hachage de réception
	ReceiptHash []byte
	// Nœud racine du hachage de transaction
	TransactionHashRoot []byte
	Difficulty          int
	// Code à usage unique choisi au hasard pour transmettre le mot de passe en toute sécurité et empêcher les attaques par rejeu
	Nonce string
	// Il détermine combien de transactions peuvent être stockées dans un bloc en fonction de la somme de gaz
	// Par exemple, lorsque la limite de gaz du bloc est de 100 et que nous avons des transactions dont les limites
	// de gaz sont de 50, 50 et 10. Block ne peut stocker que les deux premières transactions (les 50)
	// mais pas la dernière (10).
	GasLimit  int
	GasUsed   int
	ExtraData []interface{}
}

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
	// le nonce dans la transaction est un nonce de compte qui représente un ordre de transaction qu'un compte crée.
	Timestamp int64
}

// Important to note that each output is Indivisible.
// Vous ne pouvez pas "faire de changement" avec n'importe quelle sortie.
// Si la valeur est 10, afin de donner 5 à quelqu'un, nous devons faire deux sorties de cinq pièces.
// TxInput is represntative of a reference to a previous TxOutput
type TxInput struct {
	// ID will find the Transaction that a specific output is inside of
	ID []byte
	// Out will be the index of the specific output we found within a transaction.
	// For example if a transaction has 4 outputs, we can use this "out" field to specify which output we are looking for
	Out int
	// This would be a script that adds data to an outputs' PubKey
	// however for this tutorial the Sig will be indentical to the PubKey.
	Sig        []byte
	PubKey     []byte
	SchnorrKey []byte
}

//TxOutput represents a transaction in the blockchain
//For Example, I sent you 5 coins. Value would == 5, and it would have my unique PubKey
type TxOutput struct {
	// Value would be representative of the amount of coins in a transaction
	Value *big.Int
	// La Pubkey est nécessaire pour "déverrouiller" toutes les pièces dans une sortie. Cela indique que VOUS êtes celui qui l'a envoyé.
	// Vous êtes identifiable par votre PubKey
	// PubKey dans cette itération sera très simple, mais dans une application réelle, il s'agit d'un algorithme plus complexe
	PubKey []byte
}

func (tx *Transaction) IsCoinBase() bool {
	//This checks a transaction and will only return true if it is a newly minted "coin"
	// Aka a CoinBase transaction
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func (in *TxInput) CanUnlock(data []byte) bool {
	res := bytes.Compare(in.Sig, data)
	return res == 0
}

func (out *TxOutput) CanBeUnlocked(data []byte) bool {
	res := bytes.Compare(out.PubKey, data)
	return res == 0
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

var (
	mutex      = &sync.Mutex{}
	BlockChain []Block
)

/*
type BlockchainConstrucktor struct {}

func Init() *BlockchainConstrucktor{
	return &BlockchainConstrucktor{}
}*/

func Genesis(coinBase *Transaction) *Block {
	genesisBlock := Block{}
	genesisBlock = AddBlock([]byte{}, big.NewInt(1), coinBase)

	spew.Dump(genesisBlock)

	mutex.Lock()
	BlockChain = append(BlockChain, genesisBlock)
	mutex.Unlock()

	return &genesisBlock
}

// calculateHash SHA256 hashing
func calculateHash(block Block) []byte {
	record := block.Index.String() + strconv.FormatInt(block.Timestamp, 16) + string(block.PrevHash) + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return []byte(hex.EncodeToString(hashed))
}

// IsBlockValid make sure block is valid by checking index, and comparing the hash of the previous block
func IsBlockValid(newBlock, oldBlock Block) bool {
	newIndexWaiting := big.NewInt(0)
	newIndexWaiting = newIndexWaiting.Add(oldBlock.Index, big.NewInt(1))

	if newIndexWaiting.Cmp(newBlock.Index) != 0 {
		fmt.Println(fmt.Sprintf("block is invalid, with index have %d want %d", newBlock.Index, newIndexWaiting))
		return false
	}

	res := bytes.Compare(oldBlock.Hash, newBlock.PrevHash)
	if res != 0 {
		fmt.Println(fmt.Sprintf("block is invalid, prev hash is %s want %s", newBlock.PrevHash, oldBlock.Hash))
		return false
	}

	res = bytes.Compare(calculateHash(newBlock), newBlock.Hash)
	if res != 0 {
		fmt.Println(fmt.Sprintf("block is invalid %d with compare calculateHash", res))
		return false
	}

	return true
}

func NextID(index *big.Int) *big.Int {
	newIndex := big.NewInt(0)
	return newIndex.Add(index, big.NewInt(1))
}

// AddBlock create a new block using previous block's hash
func AddBlock(lastHash []byte, index *big.Int, coinBase *Transaction) Block {
	var (
		t        = time.Now().UnixNano() / int64(time.Millisecond)
		newBlock Block
		i        uint
	)

	newBlock = Block{
		Index:        NextID(index),
		Timestamp:    t,
		PrevHash:     lastHash,
		Difficulty:   difficulty.Current.Int(),
		Transactions: []*Transaction{coinBase},
	}

	for i = 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		if !isHashValid(calculateHash(newBlock), newBlock.Difficulty) {
			fmt.Println(calculateHash(newBlock), " do more work!")
			// time.Sleep(10 * time.Millisecond)
			continue
		} else {
			fmt.Println(calculateHash(newBlock), " work done!")
			newBlock.Hash = calculateHash(newBlock)
			break
		}
	}
	difficulty.Current.Update(i)
	newBlock.Difficulty = difficulty.Current.Int()

	return newBlock
}

func GetLastBlock() Block {
	return BlockChain[len(BlockChain)-1]
}

func isHashValid(hash []byte, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(string(hash), prefix)
}

// @todo parse all blockChain
func IsValid(_ []Block) bool {

	return true
}
