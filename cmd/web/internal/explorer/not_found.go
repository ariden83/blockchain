package explorer

import "net/http"

type notFoundData struct {
	PageTitle string
}

func notFoundPage(rw http.ResponseWriter, r *http.Request) {
	data := notFoundData{"404 Not Found"}

	templates.ExecuteTemplate(rw, "not_found", data)
}
