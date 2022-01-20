package explorer

import (
	"net/http"
)

func (e *Explorer) homePage(rw http.ResponseWriter, r *http.Request) {

	data := e.frontData(rw, r).Title(e.locales.LocID(r, "HomePageTitle"))

	e.ExecuteTemplate(rw, r, "home", data)
}
