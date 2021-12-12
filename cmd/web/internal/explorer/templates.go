package explorer

import (
	"text/template"
)

var templates *template.Template

var templateFunctions template.FuncMap

func (e *Explorer) initTemplates() {
	templateFunctions = template.FuncMap{
		"debug":       debug,
		"increment":   increment,
		"add":         add,
		"unixToHuman": unixToHuman,
		"homeURL":     e.homeURL,
		"blockURL":    e.blockURL,
		"txURL":       e.txURL,
		"walletURL":   e.walletURL,
	}
}

func (e *Explorer) loadTemplates() {
	e.initTemplates()
	e.loadPages()
	e.loadPartials()
}

func (e *Explorer) loadPages() {
	templates = template.Must(
		template.
			New("templates").
			Funcs(templateFunctions).
			// Delims("[[", "]]").
			ParseGlob(e.cfg.TemplatesDir + "pages/*.gohtml"))
}

func (e *Explorer) loadPartials() {
	templates = template.Must(
		templates.
			ParseGlob(e.cfg.TemplatesDir + "partials/*.gohtml"))
}
