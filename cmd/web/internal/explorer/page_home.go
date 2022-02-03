package explorer

import (
	"net/http"
)

func (e *Explorer) homePage(rw http.ResponseWriter, r *http.Request) {

	data := e.frontData(rw, r).JS([]string{
		"/static/home/home.js?v0.0.12",
	}).Title(e.locales.LocID(r, "HomePageTitle"))

	e.ExecuteTemplate(rw, r, "home", data)
}
