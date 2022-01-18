// https://www.blockchain-altcoin.com/privacy-policy
package explorer

import (
	"net/http"
)

func (e *Explorer) privacyPolicyPage(rw http.ResponseWriter, r *http.Request) {
	_, authorized := e.authorized(rw, r)
	data := FrontData{
		PageTitle:    e.metadata.Title + " - privacy policy",
		Authentified: authorized,
		Menus:        getMenus(),
	}

	templates.ExecuteTemplate(rw, "privacy-policy", data)
}
