package router

import (
	"dota2_fantasy/src/util"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

	// TODO: need to validate response.State. Likely need to store something in a cookie

	// Get Token
	baseURL := "https://oauth2.googleapis.com/token"
	redirectURI := "http://localhost:8080/api/handleOIDC"
	grantType := "authorization_code"

	postTokenCall := fmt.Sprintf("%s?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s&grant_type=%s",
		baseURL, ar.secrets.GoogleClientID, ar.secrets.GoogleClientSecret, redirectURI, response.code, grantType)

	fmt.Println(postTokenCall)

	tokenResponse, err := http.Post(postTokenCall, "application/json", nil)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	type googleTokenResponse struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
		IDToken     string `json:"id_token"`
	}

	tokenBody := &googleTokenResponse{}
	decodeErr := json.NewDecoder(tokenResponse.Body).Decode(tokenBody)
	if decodeErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	// Get emailInfo
	emailCall := "https://www.googleapis.com/oauth2/v2/userinfo"

	emailCallReq, err := http.NewRequest("GET", emailCall, nil)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	emailCallReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenBody.AccessToken))
	client := http.Client{}
	emailResponse, err := client.Do(emailCallReq)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	type fetchedEmail struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Picture       string `json:"picture"`
	}

	emailRes := &fetchedEmail{}

	decodeErr = json.NewDecoder(emailResponse.Body).Decode(emailRes)
	if decodeErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	if !emailRes.VerifiedEmail {
		w.WriteHeader(http.StatusUnauthorized)
	}
	// TODO: value here should be a session token
	cookie := http.Cookie{Name: "OIDC-Cookie", Value: "my-test-value", HttpOnly: true, Secure: true, Path: "/"}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "http://localhost:3000", http.StatusSeeOther)
}

func (ar AuthnRouter) SetupRoutes(baseRouter *mux.Router) {
	subRouter := baseRouter.PathPrefix("/api").Subrouter()

	subRouter.HandleFunc("/startOIDC", ar.getRedirectURL)
	subRouter.HandleFunc("/handleOIDC", ar.handleOIDCResponse)

	// baseRouter.
}
