package repo

import (
	"context"
	"dota2_fantasy/src/model"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/thanhpk/randstr"
)

type AuthnRepo interface {
	Persist(pool *pgxpool.Pool, authn *model.Authn) error
	GetAuthnByToken(pool *pgxpool.Pool, token string) (*model.Authn, error)
	GetAuthnByEmail(pool *pgxpool.Pool, email string) (*model.Authn, error)
	GenerateNewSessionToken(pool *pgxpool.Pool, authnID int) (string, error)
	ClearSessionToken(pool *pgxpool.Pool, authnID int) error
	UpdateLastActionTime(pool *pgxpool.Pool, authnID int) error
}

func NewAuthnRepo() AuthnRepo {
	return authnRepo{}
}

type authnRepo struct{}

func (c authnRepo) Persist(pool *pgxpool.Pool, authn *model.Authn) error {
	token := randstr.String(128)

	authn.SessionToken = token

	err := pool.QueryRow(context.Background(), `INSERT INTO authn (super_admin, email, display_name, last_action, session_token) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		authn.SuperAdmin, authn.Email, authn.DisplayName, time.Now(), token).Scan(&authn.ID)

	return err
}

func (c authnRepo) GetAuthnByToken(pool *pgxpool.Pool, token string) (*model.Authn, error) {

	authn := &model.Authn{}
	err := pool.QueryRow(context.Background(), `SELECT id, super_admin, email, display_name, last_action, session_token FROM authn WHERE session_token = $1`, token).
		Scan(&authn.ID, &authn.SuperAdmin, &authn.Email, &authn.DisplayName, &authn.LastAction, &authn.SessionToken)

	return authn, err
}

func (c authnRepo) GetAuthnByEmail(pool *pgxpool.Pool, email string) (*model.Authn, error) {

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
