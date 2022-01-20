package explorer

import (
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"net/http"
	"path/filepath"
	"text/template"
)

var templates *template.Template

var templateFunctions template.FuncMap

func (e *Explorer) initTemplates() {
	templateFunctions = template.FuncMap{
		"debug": debug,
	}
}

func debug(i interface{}) string {
	log := fmt.Sprintf("DEBUG: %v\n", i)
	fmt.Print(log)

	return log
}

func (e *Explorer) loadTemplates() {
	e.initTemplates()
	e.loadPages()
	e.loadPartials()
}

func (e *Explorer) loadPages() {
	dirs := []string{
		e.cfg.TemplatesDir + "pages/*.html",
		e.cfg.TemplatesDir + "pages/*/*.html",
	}

	files := []string{}
	for _, dir := range dirs {
		ff, err := filepath.Glob(dir)
		if err != nil {
			panic(err)
		}
		files = append(files, ff...)
	}

	t, err := template.New("templates").
		Funcs(templateFunctions).
		ParseFiles(files...)

	if err != nil {
		e.log.Error("fail to load templates", zap.Error(err))
	}
	templates = template.Must(t, err)
}

func (e *Explorer) loadPartials() {
	templates = template.Must(templates.ParseGlob(e.cfg.TemplatesDir + "partials/*.html"))
}

func (e *Explorer) ExecuteTemplate(rw http.ResponseWriter, r *http.Request, name string, data interface{}) {
	tmpl := templates.Lookup(name + lang(r))
	if tmpl == nil {
		e.log.Error(fmt.Sprintf("no template found to execute"), zap.String("template", name))
		return
	}
	tmpl.Execute(rw, data)
}

func lang(r *http.Request) string {
	l := r.FormValue("lang")
	accept := r.Header.Get("Accept-Language")

	t := parseTags(l, accept)
	if len(t) == 0 {
		return ""
	}
	m := language.NewMatcher(t)
	_, i, _ := m.Match(t...)
	d, _ := t[i].Base()
	if d.String() == "en" {
		return ""
	}
	return "-" + d.String()
}

func parseTags(langs ...string) []language.Tag {
	tags := []language.Tag{}
	for _, lang := range langs {
		t, _, err := language.ParseAcceptLanguage(lang)
		if err != nil {
			continue
		}
		tags = append(tags, t...)
	}
	return tags
}
