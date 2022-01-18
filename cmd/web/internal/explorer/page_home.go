package explorer

import (
	"net/http"
)

type FrontData struct {
	PageTitle    string
	Authentified bool
	Menus        []Menus
	Javascripts  []string
}

type Menus struct {
	Identifier  string
	Name        string
	Title       string
	URL         string
	HasChildren bool
	Pre         string
}

func (e *Explorer) homePage(rw http.ResponseWriter, r *http.Request) {
	_, authorized := e.authorized(rw, r)
	data := FrontData{
		PageTitle:    "Home",
		Authentified: authorized,
		Menus:        getMenus(),
	}
	templates.ExecuteTemplate(rw, "home", data)
}
