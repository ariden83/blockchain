package http

import (
	"fmt"
	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/internal/iterator"
	"go.uber.org/zap"
	"io"
	"net/http"
)

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
			block.Hash,
		)); err != nil {
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
	e.event.Push(event.Message{Type: event.BlockChain})
}
