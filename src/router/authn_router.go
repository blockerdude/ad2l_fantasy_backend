package router

import (
	"dota2_fantasy/src/util"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type AuthnRouter struct {
	secrets util.Secrets
}

func NewAuthnRouter(secrets util.Secrets) *AuthnRouter {
	return &AuthnRouter{secrets: secrets}
}

func (ar AuthnRouter) getRedirectURL(w http.ResponseWriter, r *http.Request) {

	// https: //accounts.google.com/o/oauth2/v2/auth?&state=myCustomValueForState&prompt=select_account
	baseURL := "https://accounts.google.com/o/oauth2/v2/auth"
	redirectURL := "http://localhost:8080/api/handleOIDC"
	responseType := "code"
	scope := "https://www.googleapis.com/auth/userinfo.email"
	accessType := "online"
	state := "todo-state-value"
	prompt := "select_account"
	redirect := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=%s&scope=%s&access_type=%s&state=%s&prompt=%s",
		baseURL, ar.secrets.GoogleClientID, redirectURL, responseType, scope, accessType, state, prompt)
	// w.WriteHeader(http.StatusSeeOther)
	w.Write([]byte(redirect))
}

func (ar AuthnRouter) handleOIDCResponse(w http.ResponseWriter, r *http.Request) {

	a := r.URL.Query()
	for b, c := range a {
		fmt.Printf("%s:%v\n", b, c)
	}

	// w.WriteHeader(http.StatusSeeOther)
	w.Write([]byte("GOT HERE"))
}

func (ar AuthnRouter) SetupRoutes(baseRouter *mux.Router) {
	subRouter := baseRouter.PathPrefix("/api").Subrouter()

	subRouter.HandleFunc("/startOIDC", ar.getRedirectURL)
	subRouter.HandleFunc("/handleOIDC", ar.handleOIDCResponse)

	// baseRouter.
}
