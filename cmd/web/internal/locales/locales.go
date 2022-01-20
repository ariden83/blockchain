package locales

import (
	"github.com/BurntSushi/toml"
	"github.com/ariden83/blockchain/cmd/web/config"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"net/http"
)

type Locales struct {
	Bundle *i18n.Bundle
}

func New(cfg config.Locales) *Locales {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile(cfg.Path + "active.en.toml")
	bundle.MustLoadMessageFile(cfg.Path + "active.es.toml")
	bundle.MustLoadMessageFile(cfg.Path + "active.fr.toml")
	return &Locales{
		Bundle: bundle,
	}
}

func (l *Locales) Loc(r *http.Request) *i18n.Localizer {
	lang := r.FormValue("lang")
	accept := r.Header.Get("Accept-Language")
	return i18n.NewLocalizer(l.Bundle, lang, accept)
}

func (l *Locales) LocID(r *http.Request, ID string) string {
	lang := r.FormValue("lang")
	accept := r.Header.Get("Accept-Language")
	return i18n.NewLocalizer(l.Bundle, lang, accept).MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: ID}})
}
