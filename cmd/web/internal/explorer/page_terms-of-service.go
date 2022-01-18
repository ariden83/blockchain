// https://www.blockchain-altcoin.com/terms-of-service
package explorer

import (
	"net/http"
)

func (e *Explorer) termsOfServicePage(rw http.ResponseWriter, r *http.Request) {
	_, authorized := e.authorized(rw, r)
	data := FrontData{
		PageTitle:    e.metadata.Title + " - terms of service",
		Authentified: authorized,
		Menus:        getMenus(),
	}

	templates.ExecuteTemplate(rw, "terms-of-service", data)
}
