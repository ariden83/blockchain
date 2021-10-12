package endpoint

import (
	"encoding/json"
	"fmt"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/handle"
	"github.com/ariden83/blockchain/internal/utils"
	"github.com/davecgh/go-spew/spew"
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
	w.Header().Set("Content-Type", "application/json")
	var m WriteBlockInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
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

	respondWithJSON(w, r, http.StatusCreated, newBlock)

}

func (e *EndPoint) getLastBlock() ([]byte, *big.Int) {
	lastHash, err := e.persistence.GetLastHash()
	handle.Handle(err)

	if lastHash == nil {
		handle.Handle(fmt.Errorf("no hash found"))
	}

	serializeBloc, err := e.persistence.GetCurrentHashSerialize(lastHash)
	handle.Handle(err)
	block, err := utils.Deserialize(serializeBloc)
	handle.Handle(err)

	return lastHash, block.Index
}

func (e *EndPoint) WriteBlock(p WriteBlockInput) blockchain.Block {
	lastHash, index := e.getLastBlock()

	//mutex.Lock()
	cbtx := e.transaction.CoinBaseTx(p.PubKey, "")
	newBlock := blockchain.AddBlock(lastHash, index, cbtx)
	//mutex.Unlock()

	if blockchain.IsBlockValid(newBlock, blockchain.BlockChain[len(blockchain.BlockChain)-1]) {

		mutex.Lock()
		blockchain.BlockChain = append(blockchain.BlockChain, newBlock)
		mutex.Unlock()

		ser, err := utils.Serialize(&newBlock)
		handle.Handle(err)

		err = e.persistence.Update(newBlock.Hash, ser)
		handle.Handle(err)
		spew.Dump(blockchain.BlockChain)
	} else {
		handle.Handle(fmt.Errorf("new block is invalid"))
	}

	return newBlock
}
