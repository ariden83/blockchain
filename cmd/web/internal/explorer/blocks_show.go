package explorer

import (
	"fmt"
	"net/http"

	"errors"
	"github.com/ariden83/blockchain/cmd/web/internal/utils"
	"github.com/ariden83/blockchain/internal/blockchain"
)

type blocksShowData struct {
	PageTitle string
	Block     *blockchain.Block
}

var (
	ErrBlockNotFound = errors.New("block not found")
)

func blocksShowPage(rw http.ResponseWriter, r *http.Request) {
	hash := utils.GetRoute(r, "hash")

	fmt.Println(fmt.Sprintf("%s", hash))
	block := &blockchain.Block{}
	/* block, err := blockchain.FindBlock(hash)
	if err == ErrBlockNotFound {
		http.Redirect(rw, r, "/404", http.StatusFound)
	} */

	data := blocksShowData{"Show Block", block}
	templates.ExecuteTemplate(rw, "blocks_show", data)
}
