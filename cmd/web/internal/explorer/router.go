package explorer

import (
	"net/http"
)

func (e *Explorer) loadRoutes() {
	e.router.HandleFunc("/", e.homePage).Methods(http.MethodGet)
	e.router.HandleFunc("/contact", e.contactPage).Methods(http.MethodGet, http.MethodPost)
	e.router.HandleFunc("/about", e.aboutPage).Methods(http.MethodGet)
	e.router.HandleFunc("/404", notFoundPage).Methods(http.MethodGet)
	e.router.HandleFunc("/privacy-policy", e.privacyPolicyPage).Methods(http.MethodGet)
	e.router.HandleFunc("/terms-of-service", e.termsOfServicePage).Methods(http.MethodGet)
}
func (e *Explorer) loadNonConnectedRoutes() {
	s := e.router.PathPrefix("/").Subrouter().StrictSlash(true)
	s.HandleFunc("/login", e.loginPage).Methods(http.MethodGet)
	s.HandleFunc("/inscription", e.inscriptionPage).Methods(http.MethodGet)
	s.HandleFunc("/authorize", e.authorizePage).Methods(http.MethodGet)
	s.Use(e.hasValidToken)
}
func (e *Explorer) loadConnectedRoutes() {
	s := e.router.PathPrefix("/").Subrouter().StrictSlash(true)
	s.HandleFunc("/wallet", e.walletPage).Methods(http.MethodGet)
	s.HandleFunc("/logout", e.logoutPage).Methods(http.MethodGet)
	s.Use(e.validateToken)
}

func (e *Explorer) loadAPIRoutes() {
	s := e.router.PathPrefix("/api").Subrouter().StrictSlash(true)
	s.Use(jsonHeader)
}
func (e *Explorer) loadAPINonConnectedRoutes() {
	s := e.router.PathPrefix("/api").Subrouter().StrictSlash(true)
	s.HandleFunc("/login", e.loginAPI).Methods(http.MethodPost)
	s.HandleFunc("/inscription", e.inscriptionAPI).Methods(http.MethodPost)
	s.Use(jsonHeader)
}
func (e *Explorer) loadAPIConnectedRoutes() {
	s := e.router.PathPrefix("/api").Subrouter().StrictSlash(true)
	s.HandleFunc("/token", e.tokenAPI).Methods(http.MethodPost)
	s.Use(e.validateToken)
	s.Use(jsonHeader)
}
