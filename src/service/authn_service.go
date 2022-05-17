package service

import (
	"context"
	fantasycontext "dota2_fantasy/src/fantasyContext"
	"dota2_fantasy/src/model"
	"dota2_fantasy/src/repo"
	"dota2_fantasy/src/util"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type AuthnService interface {
	HandleOIDCLogin(ctx context.Context, oidcCode string) (*model.Authn, error)
}

func NewAuthnService(config util.Config, authnRepo repo.AuthnRepo) AuthnService {
	return authnService{
		config:    config,
		authnRepo: authnRepo,
	}
}

type authnService struct {
	config    util.Config
	authnRepo repo.AuthnRepo
}

type fetchedEmail struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

func (a authnService) HandleOIDCLogin(ctx context.Context, oidcCode string) (*model.Authn, error) {

	accessToken, err := a.getGoogleAccessToken(oidcCode)
	if err != nil {
		return nil, err
	}

	emailInfo, err := a.getGoogleEmailFromAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	return a.authnRepo.GetUserByEmail(fantasycontext.GetDBPool(ctx), emailInfo.Email)
}

func (a authnService) getGoogleAccessToken(oidcCode string) (string, error) {
	baseURL := "https://oauth2.googleapis.com/token"
	grantType := "authorization_code"

	postTokenCall := fmt.Sprintf("%s?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s&grant_type=%s",
		baseURL, a.config.Secrets.GoogleClientID, a.config.Secrets.GoogleClientSecret, a.config.OIDC.ServerRedirectURL, oidcCode, grantType)

	tokenResponse, err := http.Post(postTokenCall, "application/json", nil)
	if err != nil {
		return "", err
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
		return "", decodeErr
	}

	return tokenBody.AccessToken, nil
}

func (a authnService) getGoogleEmailFromAccessToken(accessToken string) (*fetchedEmail, error) {
	// Get emailInfo
	emailCall := "https://www.googleapis.com/oauth2/v2/userinfo"

	emailCallReq, err := http.NewRequest("GET", emailCall, nil)
	if err != nil {
		return nil, err
	}

	emailCallReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	client := http.Client{}
	emailResponse, err := client.Do(emailCallReq)
	if err != nil {
		return nil, err
	}

	emailRes := &fetchedEmail{}

	decodeErr := json.NewDecoder(emailResponse.Body).Decode(emailRes)
	if decodeErr != nil {
		return nil, decodeErr
	}

	if !emailRes.VerifiedEmail {
		return nil, errors.New("Email not verified")
	}

	return emailRes, nil
}
