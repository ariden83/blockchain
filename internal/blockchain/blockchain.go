package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	difficulty int = 1
)

type Block struct {
	Index      *big.Int
	Timestamp  string
	BPM        int
	Hash       []byte
	PrevHash   []byte
	Difficulty int
	Nonce      string
}

var mutex = &sync.Mutex{}
var Blockchain []Block

/*
type BlockchainConstrucktor struct {}

func Init() *BlockchainConstrucktor{
	return &BlockchainConstrucktor{}
}*/

func Genesis() *Block {
	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{big.NewInt(1), t.String(), 0, calculateHash(genesisBlock), []byte{}, difficulty, ""}
	spew.Dump(genesisBlock)

	mutex.Lock()
	Blockchain = append(Blockchain, genesisBlock)
	mutex.Unlock()

	return &genesisBlock
}

// SHA256 hasing
func calculateHash(block Block) []byte {
	record := block.Index.String() + block.Timestamp + strconv.Itoa(block.BPM) + string(block.PrevHash) + block.Nonce
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
func GenerateBlock(lastHash []byte, index *big.Int, BPM int) Block {
	var newBlock Block

	t := time.Now()
	newIndex := index.Add(index, big.NewInt(1))

	newBlock.Index = newIndex
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = lastHash
	newBlock.Difficulty = difficulty

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
