package repo

import (
	"context"
	"dota2_fantasy/src/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

type AuthnRepo interface {
	GetUserByToken(conn *pgxpool.Pool, token string) (*model.Authn, error)
	GetUserByEmail(conn *pgxpool.Pool, email string) (*model.Authn, error)
}

func NewAuthnRepo() AuthnRepo {
	return authnRepo{}
}

type authnRepo struct{}

func (c authnRepo) GetUserByToken(pool *pgxpool.Pool, token string) (*model.Authn, error) {

	authn := &model.Authn{}
	err := pool.QueryRow(context.Background(), `SELECT id, super_admin, email, display_name, last_action, session_token FROM authn WHERE session_token = $1`, token).
		Scan(&authn.ID, &authn.SuperAdmin, &authn.Email, &authn.DisplayName, &authn.LastAction, &authn.SessionToken)

	return authn, err
}

func (c authnRepo) GetUserByEmail(pool *pgxpool.Pool, email string) (*model.Authn, error) {

	authn := &model.Authn{}
	err := pool.QueryRow(context.Background(), `SELECT id, super_admin, email, display_name, last_action, session_token FROM authn WHERE email = $1`, email).
		Scan(&authn.ID, &authn.SuperAdmin, &authn.Email, &authn.DisplayName, &authn.LastAction, &authn.SessionToken)

	return authn, err
}
