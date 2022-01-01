// https://www.blockchain-altcoin.com/terms-of-service
package explorer

import (
	"net/http"
)

func (e *Explorer) termsOfServicePage(rw http.ResponseWriter, r *http.Request) {
	data := homeData{e.metadata.Title + " - terms of service"}

	templates.ExecuteTemplate(rw, "terms-of-service", data)
}
