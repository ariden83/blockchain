package explorer

import (
	"net/http"
)

func (e *Explorer) aboutPage(rw http.ResponseWriter, r *http.Request) {
	_, authorized := e.authorized(rw, r)
	data := FrontData{
		PageTitle:    e.metadata.Title + "- About us",
		Authentified: authorized,
		Menus:        getMenus(),
	}
	templates.ExecuteTemplate(rw, "about", data)
}
