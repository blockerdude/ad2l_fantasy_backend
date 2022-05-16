package router

import (
	"dota2_fantasy/src/service"
	"dota2_fantasy/src/util"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thanhpk/randstr"
)

const (
	OIDC_COOKIE    = "oidc"
	SESSION_COOKIE = "fantasy-session"
)

type AuthnRouter struct {
	config   util.Config
	authnSvc service.AuthnService
}

func NewAuthnRouter(conf util.Config, authnSvc service.AuthnService) *AuthnRouter {
	return &AuthnRouter{config: conf, authnSvc: authnSvc}
}

func (ar AuthnRouter) getRedirectURL(w http.ResponseWriter, r *http.Request) {

	baseURL := "https://accounts.google.com/o/oauth2/v2/auth"
	redirectURL := ar.config.OIDC.ServerRedirectURL
	responseType := "code"
	scope := "https://www.googleapis.com/auth/userinfo.email"
	accessType := "online"
	state := randstr.String(32)
	prompt := "select_account"
	redirect := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=%s&scope=%s&access_type=%s&state=%s&prompt=%s",
		baseURL, ar.config.Secrets.GoogleClientID, redirectURL, responseType, scope, accessType, state, prompt)

	cookie := http.Cookie{Name: OIDC_COOKIE, Value: state, HttpOnly: true, Secure: true, Path: "/"}

	http.SetCookie(w, &cookie)
	w.Write([]byte(redirect))
}

type googleOIDCResponse struct {
	code     string
	scope    string
	authuser int
	prompt   string
	state    string
}

func (ar AuthnRouter) handleOIDCResponse(w http.ResponseWriter, r *http.Request) {

	response := googleOIDCResponse{}
	queryRes := r.URL.Query()
	response.code = queryRes.Get("code")
	response.scope = queryRes.Get("scope")
	response.prompt = queryRes.Get("prompt")
	response.state = queryRes.Get("state")
	var err error
	response.authuser, err = strconv.Atoi(queryRes.Get("authuser"))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	oidcCookie, err := r.Cookie(OIDC_COOKIE)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	if oidcCookie.Value != response.state {
		// State not matching indicates a malicious attempt
		w.WriteHeader(http.StatusUnauthorized)
	}

	authn, err := ar.authnSvc.HandleOIDCLogin(response.code)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	body, err := json.Marshal(authn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "%v", string(body))

	deleteOIDCCookie := http.Cookie{Name: OIDC_COOKIE, Value: "", HttpOnly: true, Secure: true, MaxAge: -1, Path: "/"}
	SessionCookie := http.Cookie{Name: SESSION_COOKIE, Value: "my-test-value", HttpOnly: true, Secure: true, Path: "/"}
	http.SetCookie(w, &SessionCookie)
	http.SetCookie(w, &deleteOIDCCookie)
	http.Redirect(w, r, ar.config.OIDC.UIBaseURL, http.StatusSeeOther)
}

func (ar AuthnRouter) SetupRoutes(baseRouter *mux.Router) {
	subRouter := baseRouter.PathPrefix("/api").Subrouter()

	subRouter.HandleFunc("/startOIDC", ar.getRedirectURL)

	subRouter.HandleFunc("/handleOIDC", ar.handleOIDCResponse)

	// baseRouter.
}
