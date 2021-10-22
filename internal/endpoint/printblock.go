package endpoint

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
