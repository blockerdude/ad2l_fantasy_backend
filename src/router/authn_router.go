package router

import (
	fantasycontext "dota2_fantasy/src/fantasyContext"
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
	config     util.Config
	middleware Middleware
	authnSvc   service.AuthnService
}

func NewAuthnRouter(conf util.Config, middleware Middleware, authnSvc service.AuthnService) *AuthnRouter {
	return &AuthnRouter{config: conf, middleware: middleware, authnSvc: authnSvc}
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

func (ar AuthnRouter) handleOIDCResponse(w http.ResponseWriter, req *http.Request) {

	response := googleOIDCResponse{}
	queryRes := req.URL.Query()
	response.code = queryRes.Get("code")
	response.scope = queryRes.Get("scope")
	response.prompt = queryRes.Get("prompt")
	response.state = queryRes.Get("state")
	var err error
	response.authuser, err = strconv.Atoi(queryRes.Get("authuser"))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	oidcCookie, err := req.Cookie(OIDC_COOKIE)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if oidcCookie.Value != response.state {
		// State not matching indicates a malicious attempt
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionToken, err := ar.authnSvc.HandleOIDCLogin(req.Context(), response.code)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	deleteOIDCCookie := http.Cookie{Name: OIDC_COOKIE, Value: "", HttpOnly: true, Secure: true, MaxAge: -1, Path: "/"}
	SessionCookie := http.Cookie{Name: SESSION_COOKIE, Value: sessionToken, HttpOnly: true, Secure: true, Path: "/"}
	http.SetCookie(w, &SessionCookie)
	http.SetCookie(w, &deleteOIDCCookie)
	http.Redirect(w, req, ar.config.OIDC.UIBaseURL, http.StatusSeeOther)
}

func (ar AuthnRouter) getLoggedInAuthn(w http.ResponseWriter, req *http.Request) {

	authn := fantasycontext.GetAuthn(req.Context())

	res, err := json.Marshal(authn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(res)
}

func (ar AuthnRouter) logoutHandler(w http.ResponseWriter, req *http.Request) {
	ar.authnSvc.LogoutUser(req.Context())

	SessionCookie := http.Cookie{Name: SESSION_COOKIE, Value: "", HttpOnly: true, Secure: true, MaxAge: -1, Path: "/"}
	http.SetCookie(w, &SessionCookie)

	w.WriteHeader(http.StatusNoContent)
}

func (ar AuthnRouter) SetupRoutes(baseRouter *mux.Router) {
	subRouter := baseRouter.PathPrefix("/api").Subrouter()

	subRouter.HandleFunc("/startOIDC", ar.middleware.WithBaseMiddleware(ar.getRedirectURL))

	subRouter.HandleFunc("/handleOIDC", ar.middleware.WithBaseMiddleware(ar.handleOIDCResponse))

	subRouter.HandleFunc("/authn", ar.middleware.WithBaseMiddleware(ar.getLoggedInAuthn, ar.middleware.RequireLogin)).Methods(http.MethodGet)

	subRouter.Handle("/logout", ar.middleware.WithBaseMiddleware(ar.logoutHandler, ar.middleware.RequireLogin)).Methods(http.MethodGet)
	// baseRouter.
}
