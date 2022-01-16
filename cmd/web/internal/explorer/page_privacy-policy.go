// https://www.blockchain-altcoin.com/privacy-policy
package explorer

import (
	"net/http"
)

func (e *Explorer) privacyPolicyPage(rw http.ResponseWriter, r *http.Request) {
	_, authorized := e.authorized(rw, r)
	data := frontData{e.metadata.Title + " - privacy policy", authorized}

	templates.ExecuteTemplate(rw, "privacy-policy", data)
}
