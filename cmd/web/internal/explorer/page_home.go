package explorer

import (
	"net/http"
)

type frontData struct {
	PageTitle    string
	Authentified bool
}

func (e *Explorer) homePage(rw http.ResponseWriter, r *http.Request) {
	_, authorized := e.authorized(rw, r)
	data := frontData{PageTitle: "Home", Authentified: authorized}
	templates.ExecuteTemplate(rw, "home", data)
}
