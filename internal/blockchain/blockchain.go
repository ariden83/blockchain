package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"
)

const (
	difficulty int = 1
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

//TxOutput represents a transaction in the blockchain
//For Example, I sent you 5 coins. Value would == 5, and it would have my unique PubKey
type TxOutput struct {
	Value int
	//Value would be representative of the amount of coins in a transaction
	PubKey string
	//The Pubkey is needed to "unlock" any coins within an Output. This indicated that YOU are the one that sent it.
	//You are indentifiable by your PubKey
	//PubKey in this iteration will be very straightforward, however in an actual application this is a more complex algorithm
}

//Important to note that each output is Indivisible.
//You cannot "make change" with any output.
//If the Value is 10, in order to give someone 5, we need to make two five coin outputs.

//TxInput is represntative of a reference to a previous TxOutput
type TxInput struct {
	ID []byte
	//ID will find the Transaction that a specific output is inside of
	Out int
	//Out will be the index of the specific output we found within a transaction.
	//For example if a transaction has 4 outputs, we can use this "out" field to specify which output we are looking for
	Sig string
	//This would be a script that adds data to an outputs' PubKey
	//however for this tutorial the Sig will be indentical to the PubKey.
}

type Block struct {
	Index        *big.Int
	Timestamp    string
	Transactions []*Transaction
	Hash         []byte
	PrevHash     []byte
	Difficulty   int
	Nonce        string
}

func (tx *Transaction) IsCoinBase() bool {
	//This checks a transaction and will only return true if it is a newly minted "coin"
	// Aka a CoinBase transaction
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}
func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
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

var mutex = &sync.Mutex{}
var BlockChain []Block

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

// SHA256 hasing
func calculateHash(block Block) []byte {
	record := block.Index.String() + block.Timestamp + string(block.PrevHash) + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return []byte(hex.EncodeToString(hashed))
}

// make sure block is valid by checking index, and comparing the hash of the previous block
func IsBlockValid(newBlock, oldBlock Block) bool {
	newIndexWaiting := oldBlock.Index.Add(oldBlock.Index, big.NewInt(1))

	if newIndexWaiting.Cmp(newBlock.Index) != 0 {
		return false
	}

	res := bytes.Compare(oldBlock.Hash, newBlock.PrevHash)
	if res != 0 {
		return false
	}

	res = bytes.Compare(calculateHash(newBlock), newBlock.Hash)
	if res != 0 {
		return false
	}

	return true
}

// create a new block using previous block's hash
func AddBlock(lastHash []byte, index *big.Int, coinBase *Transaction) Block {
	t := time.Now()
	newIndex := index.Add(index, big.NewInt(1))

	var newBlock Block = Block{
		Index:        newIndex,
		Timestamp:    t.String(),
		PrevHash:     lastHash,
		Difficulty:   difficulty,
		Transactions: []*Transaction{coinBase},
	}

	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		if !isHashValid(calculateHash(newBlock), newBlock.Difficulty) {
			fmt.Println(calculateHash(newBlock), " do more work!")
			time.Sleep(time.Second)
			continue
		} else {
			fmt.Println(calculateHash(newBlock), " work done!")
			newBlock.Hash = calculateHash(newBlock)
			break
		}

	}
	return newBlock
}

func isHashValid(hash []byte, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(string(hash), prefix)
}
