package explorer

import (
	"net/http"

	"github.com/ariden83/blockchain/internal/blockchain"
)

type blocksIndexData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func blocksIndex(rw http.ResponseWriter, r *http.Request) {
	blocks := []*blockchain.Block{}
	data := blocksIndexData{"Blocks", blocks}

	templates.ExecuteTemplate(rw, "blocks_index", data)
}
