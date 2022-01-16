package explorer

import (
	"net/http"
)

func (e *Explorer) loadRoutes() {
	e.router.HandleFunc("/", e.homePage).Methods(http.MethodGet)
	e.router.HandleFunc("/login", e.loginPage).Methods(http.MethodGet)
	e.router.HandleFunc("/inscription", e.inscriptionPage).Methods(http.MethodGet)
	e.router.HandleFunc("/404", notFoundPage).Methods(http.MethodGet)
	e.router.HandleFunc("/privacy-policy", e.privacyPolicyPage).Methods(http.MethodGet)
	e.router.HandleFunc("/terms-of-service", e.termsOfServicePage).Methods(http.MethodGet)
}

func (e *Explorer) loadConnectedRoutes() {
	s := e.router.PathPrefix("/").Subrouter().StrictSlash(true)
	s.HandleFunc("/wallet", e.walletPage).Methods(http.MethodGet)

	s.Use(e.validateToken)
}

func (e *Explorer) loadAPIRoutes() {
	s := e.router.PathPrefix("/api").Subrouter().StrictSlash(true)

	// s.HandleFunc("/auth/google/login", e.oauthGoogleLogin)
	// s.HandleFunc("/auth/google/callback", e.oauthGoogleCallback)

	s.HandleFunc("/login", e.loginAPI).Methods(http.MethodPost)
	s.HandleFunc("/authorize", e.authorizeAPI).Methods(http.MethodGet)

	s.Use(jsonHeader)
}

func (e *Explorer) loadAPIConnectedRoutes() {
	s := e.router.PathPrefix("/api").Subrouter().StrictSlash(true)

	s.HandleFunc("/token", e.tokenAPI).Methods(http.MethodGet)
	s.Use(e.validateToken)
	s.Use(jsonHeader)
}
