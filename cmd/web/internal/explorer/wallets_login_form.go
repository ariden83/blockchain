package explorer

import (
	"fmt"
	"net/http"
)

type walletsLoginForm struct {
	PageTitle string
}

func (e *Explorer) walletsLoginForm(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("***************************************** walletsLoginForm")
	frontData := walletsLoginForm{"Wallets connexion"}
	fmt.Println("***************************************** 32+")
	templates.ExecuteTemplate(rw, "wallets_login_form", frontData)
}
