package http

import (
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/internal/iterator"
)

/*
func (e *EndPoint) handlePrintBlockChain(w http.ResponseWriter, _ *http.Request) {
	iterator := e.Iterator()

	for {
		block, err := iterator.Next()
		if err != nil {
			e.log.Error("fail to iterate next block", zap.Error(err))
			return
		}

		if _, err := io.WriteString(w, fmt.Sprintf("Previous hash: %x\ndata: %+v\nhash: %x\n",
			block.PrevHash,
			block,
			block.Hash)); err != nil {
			e.log.Error("fail to write string", zap.Error(err))
			return
		}

		// This works because the Genesis block has no PrevHash to point to.
		if len(block.PrevHash) == 0 {
			break
		}
	}
	e.event.Push(event.Message{Type: event.BlockChain})
}
*/

// Iterator takes our BlockChain struct and returns it as a BlockCHainIterator struct
/*func (e *EndPoint) Iterator() *iterator.BlockChainIterator {
	iterator := iterator.BlockChainIterator{
		CurrentHash: e.persistence.LastHash(),
		Persistence: e.persistence,
	}

	return &iterator
}*/

func (e *EndPoint) handleGetBlockChain(w http.ResponseWriter, _ *http.Request) {
	iterator := iterator.New(e.persistence)

	for {
		block, err := iterator.Next()
		if err != nil {
			e.log.Error("fail to iterate next block", zap.Error(err))
			return
		}

		if _, err = io.WriteString(w, fmt.Sprintf("Previous hash: %x\ndata: %+v\nhash: %x\n",
			block.PrevHash,
			block,
			block.Hash)); err != nil {
			e.log.Error("fail to WriteString", zap.Error(err))
			return
		}
		/*pow := blockchain.NewProofOfWork(block)
		io.WriteString(w, fmt.Sprintf("Pow: %s\n", strconv.FormatBool(pow.Validate())))*/
		// This works because the Genesis block has no PrevHash to point to.
		if len(block.PrevHash) == 0 {
			break
		}
	}
	// e.event.Push(event.Message{Type: event.BlockChain})
}
