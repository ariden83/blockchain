package endpoint

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/handle"
	"github.com/ariden83/blockchain/internal/persistence"
	"github.com/davecgh/go-spew/spew"
	"io"
	"net/http"
	"sync"
)

var mutex = &sync.Mutex{}

type EndPoint struct {
	persistence *persistence.Persistence
}

func Init(p *persistence.Persistence) *EndPoint {
	return &EndPoint{
		persistence: p,
	}
}

func (e *EndPoint) Genesis() {
	lastHash, err := e.persistence.GetLastHash()
	handle.Handle(err)
	if lastHash == nil {
		fmt.Println("No existing blockchain found")

		genesis := blockchain.Genesis()
		fmt.Println("Genesis proved")

		lastHash = genesis.Hash

		serializeBLock, err := Serialize(genesis)
		handle.Handle(err)

		err = e.persistence.Update(lastHash, serializeBLock)
		handle.Handle(err)
	} else {

		val, err := e.persistence.GetCurrentHashSerialize(lastHash)
		handle.Handle(err)
		block, err := Deserialize(val)
		handle.Handle(err)

		e.persistence.SetLastHash(lastHash)

		mutex.Lock()
		blockchain.Blockchain = append(blockchain.Blockchain, *block)
		mutex.Unlock()

		spew.Dump(blockchain.Blockchain)
	}

	return
}

// Message takes incoming JSON payload for writing heart rate
type Message struct {
	BPM int `json:"bpm"`
}

func (e *EndPoint) GenerateBlock(m Message) blockchain.Block {
	//ensure atomicity when creating new block
	lastHash, err := e.persistence.GetLastHash()
	handle.Handle(err)

	if lastHash == nil {
		handle.Handle(fmt.Errorf("no hash found"))
	}

	serializeBloc, err := e.persistence.GetCurrentHashSerialize(lastHash)
	handle.Handle(err)
	block, err := Deserialize(serializeBloc)
	handle.Handle(err)

	mutex.Lock()
	newBlock := blockchain.GenerateBlock(lastHash, block.Index, m.BPM)
	mutex.Unlock()

	if blockchain.IsBlockValid(newBlock, blockchain.Blockchain[len(blockchain.Blockchain)-1]) {

		mutex.Lock()
		blockchain.Blockchain = append(blockchain.Blockchain, newBlock)
		mutex.Unlock()

		ser, err := Serialize(&newBlock)
		handle.Handle(err)

		err = e.persistence.Update(newBlock.Hash, ser)
		handle.Handle(err)
		spew.Dump(blockchain.Blockchain)
	} else {
		handle.Handle(fmt.Errorf("new block is invalid"))
	}

	return newBlock
}

func (e *EndPoint) PrintBlockChain(w http.ResponseWriter) {
	iterator := e.Iterator()

	for {
		block := iterator.Next()

		io.WriteString(w, fmt.Sprintf("Previous hash: %x\n", block.PrevHash))
		io.WriteString(w, fmt.Sprintf("data: %+v\n", block))
		io.WriteString(w, fmt.Sprintf("hash: %x\n", block.Hash))
		/*pow := blockchain.NewProofOfWork(block)
		io.WriteString(w, fmt.Sprintf("Pow: %s\n", strconv.FormatBool(pow.Validate())))*/
		io.WriteString(w, "")
		// This works because the Genesis block has no PrevHash to point to.
		if len(block.PrevHash) == 0 {
			break
		}
	}

	return
}

type BlockChainIterator struct {
	CurrentHash []byte
	persistence *persistence.Persistence
	// Database    *badger.DB
}

// Iterator takes our BlockChain struct and returns it as a BlockCHainIterator struct
func (e *EndPoint) Iterator() *BlockChainIterator {
	iterator := BlockChainIterator{
		CurrentHash: e.persistence.LastHash,
		persistence: e.persistence,
	}

	return &iterator
}

func (b *BlockChainIterator) Next() *blockchain.Block {
	val, err := b.persistence.GetCurrentHashSerialize(b.CurrentHash)
	handle.Handle(err)
	block, err := Deserialize(val)
	handle.Handle(err)
	b.CurrentHash = block.PrevHash
	return block
}

func Deserialize(data []byte) (*blockchain.Block, error) {
	var block blockchain.Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	return &block, err
}

func Serialize(b *blockchain.Block) ([]byte, error) {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)
	return res.Bytes(), err
}

func (e *EndPoint) Close() {
	e.persistence.Close()
}
