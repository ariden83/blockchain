package explorer

import (
	"fmt"
	"net/http"

	"errors"
	"github.com/ariden83/blockchain/cmd/web/internal/utils"
	"github.com/ariden83/blockchain/internal/transactions"
)

type txsShowData struct {
	PageTitle string
	Tx        *transactions.Transactions
}

var (
	ErrTxNotFound = errors.New("transaction not found")
)

func txsShow(rw http.ResponseWriter, r *http.Request) {
	id := utils.GetRoute(r, "id")
	fmt.Println(fmt.Sprintf("%s", id))
	tx := &transactions.Transactions{}
	/* tx, err := blockchain.FindTx(id)
	if err == ErrTxNotFound {
		http.Redirect(rw, r, "/404", http.StatusFound)
	}
	*/

	data := txsShowData{"Show Transaction", tx}
	templates.ExecuteTemplate(rw, "txs_show", data)
}
