// https://www.blockchain-altcoin.com/privacy-policy
package explorer

import (
	"net/http"
)

func (e *Explorer) privacyPolicyPage(rw http.ResponseWriter, r *http.Request) {
	data := homeData{e.metadata.Title + " - privacy policy"}

	templates.ExecuteTemplate(rw, "privacy-policy", data)
}
