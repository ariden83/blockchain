package http

import (
	"fmt"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/iterator"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func (e *EndPoint) handlePrintBlockChain(w http.ResponseWriter, _ *http.Request) {
	iterator := e.Iterator()

	for {
		block, err := iterator.Next()
		if err != nil {
			e.log.Fatal("fail to iterate next block", zap.Error(err))
		}

		if _, err := io.WriteString(w, fmt.Sprintf("Previous hash: %x\ndata: %+v\nhash: %x\n",
			block.PrevHash,
			block,
			block.Hash)); err != nil {
			e.log.Error("fail to write string", zap.Error(err))
		}

		/*pow := blockchain.NewProofOfWork(block)
		io.WriteString(w, fmt.Sprintf("Pow: %s\n", strconv.FormatBool(pow.Validate())))*/
		// This works because the Genesis block has no PrevHash to point to.
		if len(block.PrevHash) == 0 {
			break
		}
	}
	e.event.Push(event.Message{Type: event.BlockChain})
}

// Iterator takes our BlockChain struct and returns it as a BlockCHainIterator struct
func (e *EndPoint) Iterator() *iterator.BlockChainIterator {
	iterator := iterator.BlockChainIterator{
		CurrentHash: e.persistence.LastHash(),
		Persistence: e.persistence,
	}

	return &iterator
}
