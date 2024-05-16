package explorer

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/mailjet/mailjet-apiv3-go/v3"
	"go.uber.org/zap"
)

type contactDetails struct {
	Name    string
	Email   string
	Subject string
	Message string
}

type contactShowData struct {
	*FrontData
	Success bool
	Error   error
	Form    contactDetails
}

func (e *Explorer) contactPage(rw http.ResponseWriter, r *http.Request) {
	_, authorized := e.authorized(rw, r)

	data := contactShowData{
		Success: false,
		Form: contactDetails{
			Name:    r.FormValue("name"),
			Email:   r.FormValue("email"),
			Subject: r.FormValue("subject"),
			Message: r.FormValue("message"),
		},
		FrontData: &FrontData{
			PageTitle:    e.metadata.Title + "- contact-us",
			Authentified: authorized,
			Menus:        getMenus(),
		},
	}

	if r.Method != http.MethodPost {
		templates.ExecuteTemplate(rw, "contact", data)
		return
	}

	mj := mailjet.NewMailjetClient(e.cfg.Mails.PublicKey, e.cfg.Mails.SecretKey)

	if e.cfg.Mails.ProxyURL != "" {
		client := e.setupProxy(e.cfg.Mails.ProxyURL)
		mj.SetClient(client)
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: r.FormValue("email"),
				Name:  r.FormValue("name"),
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: "adrienparrochia@gmail.com",
					Name:  "ariden",
				},
			},
			Subject:  r.FormValue("subject"),
			TextPart: r.FormValue("message"),
			HTMLPart: r.FormValue("message"),
		},
	}
	messages := &mailjet.MessagesV31{Info: messagesInfo}

	if _, err := mj.SendMailV31(messages); err != nil {
		e.log.Error("fail to call mailjet", zap.Error(err))
		data.Error = errors.New("internal error")
	} else {
		data.Success = true
	}

	templates.ExecuteTemplate(rw, "contact", data)
}

func (e *Explorer) setupProxy(proxyURLStr string) *http.Client {
	proxyURL, err := url.Parse(proxyURLStr)
	if err != nil {
		e.log.Error("fail to call mailjet", zap.Error(err))
	}
	tr := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	client := &http.Client{}
	client.Transport = tr

	return client
}
