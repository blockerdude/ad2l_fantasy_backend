package repo

import (
	"context"
	"dota2_fantasy/src/model"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/thanhpk/randstr"
)

type AuthnRepo interface {
	GetUserByToken(conn *pgxpool.Pool, token string) (*model.Authn, error)
	GetUserByEmail(conn *pgxpool.Pool, email string) (*model.Authn, error)
	GenerateNewSessionToken(pool *pgxpool.Pool, authnID int) (string, error)
	ClearSessionToken(pool *pgxpool.Pool, authnID int) error
	UpdateLastActionTime(pool *pgxpool.Pool, authnID int) error
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

func (c authnRepo) GenerateNewSessionToken(pool *pgxpool.Pool, authnID int) (string, error) {

	token := randstr.String(128)

	_, err := pool.Exec(context.Background(), `UPDATE authn SET session_token = $2, last_action = $3 WHERE id = $1`, authnID, token, time.Now())

	return token, err
}

func (c authnRepo) ClearSessionToken(pool *pgxpool.Pool, authnID int) error {
	_, err := pool.Exec(context.Background(), `UPDATE authn SET session_token = $2 WHERE id = $1`, authnID, "")

	return err
}

func (c authnRepo) UpdateLastActionTime(pool *pgxpool.Pool, authnID int) error {
	_, err := pool.Exec(context.Background(), `UPDATE authn SET last_action = $2 WHERE id = $1`, authnID, time.Now())
	return err
}
