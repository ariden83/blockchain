package http

import (
	"encoding/json"
	"fmt"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/p2p/address"
	"github.com/ariden83/blockchain/internal/p2p/validation"
	"github.com/ariden83/blockchain/internal/utils"
	"io"
	"math/big"
	"net/http"
)

// Message takes incoming JSON payload for writing heart rate
type WriteBlockInput struct {
	Address string `json:"address"`
	PubKey  string `json:"key"`
}

// takes JSON payload as an input for heart rate (BPM)
func (e *EndPoint) handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m WriteBlockInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		e.respondWithJSON(w, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	if m.Address == "" {
		io.WriteString(w, "No address set")
		return
	}

	if m.PubKey == "" {
		io.WriteString(w, "No pub key set")
		return
	}

	newBlock := e.WriteBlock(m)
	e.respondWithJSON(w, http.StatusCreated, newBlock)

}

func (e *EndPoint) getLastBlock() ([]byte, *big.Int) {
	lastHash, err := e.persistence.GetLastHash()
	e.Handle(err)

	if lastHash == nil {
		e.Handle(fmt.Errorf("no hash found"))
	}

	serializeBloc, err := e.persistence.GetCurrentHashSerialize(lastHash)
	e.Handle(err)
	block, err := utils.Deserialize(serializeBloc)
	e.Handle(err)

	return lastHash, block.Index
}

func (e *EndPoint) WriteBlock(p WriteBlockInput) blockchain.Block {
	lastHash, index := e.getLastBlock()

	//mutex.Lock()
	cbtx := e.transaction.CoinBaseTx(p.PubKey, "")
	cbtx.SetID()

	newBlock := blockchain.AddBlock(lastHash, index, cbtx)
	//mutex.Unlock()

	if blockchain.IsBlockValid(newBlock, blockchain.BlockChain[len(blockchain.BlockChain)-1]) {

		mutex.Lock()
		e.event.PushBlock(validation.New(newBlock, address.GetCurrentAddress()))
		mutex.Unlock()

		/*mutex.Lock()
		blockchain.BlockChain = append(blockchain.BlockChain, newBlock)
		mutex.Unlock()

		ser, err := utils.Serialize(&newBlock)
		e.Handle(err)

		e.event.Push(event.BlockChain, "")

		err = e.persistence.Update(newBlock.Hash, ser)
		e.Handle(err)
		spew.Dump(blockchain.BlockChain)*/
	} else {
		e.Handle(fmt.Errorf("new block created is invalid"))
	}

	return newBlock
}
