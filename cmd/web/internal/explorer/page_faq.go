package explorer

import (
	"net/http"
)

func (e *Explorer) faqPage(rw http.ResponseWriter, r *http.Request) {
	data := e.frontData(rw, r).Title("Frequently Asked Questions")
	templates.ExecuteTemplate(rw, "faq", data)
}
