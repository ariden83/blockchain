package explorer

import (
	"net/http"
)

func (e *Explorer) sitemapPage(rw http.ResponseWriter, r *http.Request) {
	_, authorized := e.authorized(rw, r)
	data := FrontData{
		PageTitle:    "Sitemap | " + e.metadata.Title,
		Authentified: authorized,
		Menus:        getMenus(),
	}

	templates.ExecuteTemplate(rw, "sitemap", data)
}
