package explorer

import (
	"fmt"
	"net/http"
)

type walletsIndexData struct {
	PageTitle string
	Addresses []string
}

func walletsIndex(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("***************************************** walletsIndex")
	addresses := []string{}
	data := walletsIndexData{"Wallets", addresses}

	templates.ExecuteTemplate(rw, "wallets_index", data)
}
