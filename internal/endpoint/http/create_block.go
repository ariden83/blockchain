package http

import (
	"fmt"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/p2p/address"
	"github.com/ariden83/blockchain/internal/p2p/validation"
	"github.com/ariden83/blockchain/internal/utils"
	"go.uber.org/zap"
	"io"
	"math/big"
	"net/http"
)

// Message takes incoming JSON payload for writing heart rate
type CreateBlockInput struct {
	Address string `json:"address"`
	PubKey  string `json:"key"`
}

// handleCreateBlock takes JSON payload as an input for heart rate (BPM)
func (e *EndPoint) handleCreateBlock(rw http.ResponseWriter, r *http.Request) {
	req := &CreateBlockInput{}

	log := e.log.With(zap.String("input", "createBlock"))
	if err := e.decodeBody(rw, log, r.Body, req); err != nil {
		return
	}

	if req.Address == "" {
		io.WriteString(rw, "No address set")
		return
	}

	if req.PubKey == "" {
		io.WriteString(rw, "No pub key set")
		return
	}

	newBlock := e.WriteBlock(*req)
	e.JSONRes(rw, http.StatusCreated, newBlock)

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

func (e *EndPoint) WriteBlock(p CreateBlockInput) blockchain.Block {
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
