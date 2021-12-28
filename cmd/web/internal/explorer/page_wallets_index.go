package explorer

import (
	"net/http"
)

type walletsIndexData struct {
	PageTitle string
	Addresses []string
}

func walletsIndexPage(rw http.ResponseWriter, r *http.Request) {
	addresses := []string{}
	data := walletsIndexData{"Wallets", addresses}

	templates.ExecuteTemplate(rw, "wallets_index", data)
}
