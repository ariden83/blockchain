package explorer

import (
	"fmt"
	"net/http"
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
}

func (e *Explorer) contactPage(rw http.ResponseWriter, r *http.Request) {
	_, authorized := e.authorized(rw, r)

	data := contactShowData{
		Success: false,
	}
	data.FrontData = &FrontData{
		PageTitle:    e.metadata.Title + "- contact-us",
		Authentified: authorized,
		Menus:        getMenus(),
	}

	if r.Method != http.MethodPost {
		templates.ExecuteTemplate(rw, "contact", data)
		return
	}

	// get form data
	formData := contactDetails{
		Name:    r.FormValue("name"),
		Email:   r.FormValue("email"),
		Subject: r.FormValue("subject"),
		Message: r.FormValue("message"),
	}

	// do something with the submitted form data
	fmt.Printf("%+v\n", formData)
	data.Success = true

	templates.ExecuteTemplate(rw, "contact", data)
}
