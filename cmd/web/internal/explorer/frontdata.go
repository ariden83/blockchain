package explorer

import (
	"net/http"
)

type FrontData struct {
	PageTitle         string
	Authentified      bool
	Menus             []Menus
	Javascripts       []string
	ModuleJavascripts []string
	CSS               []string
	siteTitle         string
}

type Menus struct {
	Identifier  string
	Name        string
	Title       string
	URL         string
	HasChildren bool
	Pre         string
}

func (e *Explorer) frontData(rw http.ResponseWriter, r *http.Request) *FrontData {
	_, authorized := e.authorized(rw, r)
	f := &FrontData{
		Authentified: authorized,
		Menus:        getMenus(),
		siteTitle:    e.metadata.Title,
	}
	return f
}

func (f *FrontData) JS(js []string) *FrontData {
	f.Javascripts = js
	return f
}

func (f *FrontData) ModuleJS(js []string) *FrontData {
	f.ModuleJavascripts = js
	return f
}

func (f *FrontData) Css(css []string) *FrontData {
	f.CSS = css
	return f
}

func (f *FrontData) Title(title string) *FrontData {
	f.PageTitle = f.siteTitle + " - " + title
	return f
}
