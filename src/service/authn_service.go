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

	"github.com/jackc/pgx/v4"
)

type AuthnService interface {
	// HandleOIDCLogin returns the session token
	HandleOIDCLogin(ctx context.Context, oidcCode string) (string, error)

	HandleOIDCSignup(ctx context.Context, oidcCode string) (string, error)

	GetAuthnByToken(ctx context.Context, token string) (*model.Authn, error)

	LogoutUser(ctx context.Context) error

	UpdateLastActionTime(ctx context.Context) error
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

func (a authnService) HandleOIDCLogin(ctx context.Context, oidcCode string) (string, error) {

	accessToken, err := a.getGoogleAccessToken(oidcCode, a.config.OIDC.LoginRedirectURL)
	if err != nil {
		return "", err
	}

	emailInfo, err := a.getGoogleEmailFromAccessToken(accessToken)
	if err != nil {
		return "", err
	}

	pool := fantasycontext.GetDBPool(ctx)
	authn, err := a.authnRepo.GetAuthnByEmail(pool, emailInfo.Email)
	if err != nil {
		return "", err
	}

	return a.authnRepo.GenerateNewSessionToken(pool, authn.ID)
}

func (a authnService) HandleOIDCSignup(ctx context.Context, oidcCode string) (string, error) {
	accessToken, err := a.getGoogleAccessToken(oidcCode, a.config.OIDC.SignupRedirectURL)
	if err != nil {
		return "", err
	}

	emailInfo, err := a.getGoogleEmailFromAccessToken(accessToken)
	if err != nil {
		return "", err
	}

	pool := fantasycontext.GetDBPool(ctx)
	_, err = a.authnRepo.GetAuthnByEmail(pool, emailInfo.Email)
	if err != pgx.ErrNoRows {
		return "", err
	}

	newAuthn := &model.Authn{
		SuperAdmin:  false,
		Email:       emailInfo.Email,
		DisplayName: "",
	}

	if err := a.authnRepo.Persist(pool, newAuthn); err != nil {
		return "", err
	}

	return newAuthn.SessionToken, nil

}

func (a authnService) getGoogleAccessToken(oidcCode, redirectURL string) (string, error) {
	baseURL := "https://oauth2.googleapis.com/token"
	grantType := "authorization_code"

	postTokenCall := fmt.Sprintf("%s?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s&grant_type=%s",
		baseURL, a.config.Secrets.GoogleClientID, a.config.Secrets.GoogleClientSecret, redirectURL, oidcCode, grantType)

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

func (a authnService) GetAuthnByToken(ctx context.Context, token string) (*model.Authn, error) {
	return a.authnRepo.GetAuthnByToken(fantasycontext.GetDBPool(ctx), token)
}

func (a authnService) LogoutUser(ctx context.Context) error {
	return a.authnRepo.ClearSessionToken(fantasycontext.GetDBPool(ctx), fantasycontext.GetAuthn(ctx).ID)
}

func (a authnService) UpdateLastActionTime(ctx context.Context) error {
	return a.authnRepo.UpdateLastActionTime(fantasycontext.GetDBPool(ctx), fantasycontext.GetAuthn(ctx).ID)
}
