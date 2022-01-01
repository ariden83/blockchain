package explorer

import (
	"github.com/ariden83/blockchain/internal/middleware"
	"net/http"
)

func (e *Explorer) loadRoutes() {
	e.router.HandleFunc("/", homePage).Methods(http.MethodGet)
	e.router.HandleFunc("/login", e.loginPage).Methods(http.MethodGet)
	e.router.HandleFunc("/create", e.walletsCreatePage).Methods(http.MethodGet)
	e.router.HandleFunc("/404", notFoundPage).Methods(http.MethodGet)
	e.router.HandleFunc("/privacy-policy", e.privacyPolicyPage).Methods(http.MethodGet)
	e.router.HandleFunc("/terms-of-service", e.termsOfServicePage).Methods(http.MethodGet)
	/*

		e.router.HandleFunc("/blocks", blocksIndexPage).Methods(http.MethodGet)
		e.router.HandleFunc("/blocks/{hash:[0-9a-f]+}", blocksShowPage).Methods(http.MethodGet)
		e.router.HandleFunc("/blocks", blocksCreatePage).Methods("POST")
		e.router.HandleFunc("/blocks/mine", blocksMinePage).Methods(http.MethodGet)

		e.router.HandleFunc("/transactions/{id:[0-9a-f]+}", txsShowPage).Methods(http.MethodGet)

		e.router.HandleFunc("/wallets", walletsIndexPage).Methods(http.MethodGet)
		e.router.HandleFunc("/wallets/server", walletsServerPage).Methods(http.MethodGet)
		e.router.HandleFunc("/wallets/{address:[0-9a-f]+}", walletsShowPage).Methods(http.MethodGet)

	*/

}

func (e *Explorer) loadAPIRoutes() {
	s := e.router.PathPrefix("/api").Subrouter()
	s.Use(middleware.DefaultHeader)

	// s.HandleFunc("/auth/google/login", e.oauthGoogleLogin)
	// s.HandleFunc("/auth/google/callback", e.oauthGoogleCallback)

	s.HandleFunc("/oauth2", e.authorize)

	s.HandleFunc("/auth", e.oauthHandler).Methods(http.MethodPost)
	s.HandleFunc("/login", e.loginHandler).Methods(http.MethodPost)

	s.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		err := e.authServer.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	s.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		err := e.authServer.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

}
