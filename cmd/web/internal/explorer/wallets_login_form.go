package explorer

import (
	"net/http"
)

type walletsLoginForm struct {
	PageTitle string
}

func (e *Explorer) walletsLoginForm(rw http.ResponseWriter, r *http.Request) {
	frontData := walletsLoginForm{"Wallets connexion"}
	templates.ExecuteTemplate(rw, "wallets_login_form", frontData)
}
