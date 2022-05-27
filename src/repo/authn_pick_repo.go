package repo

import (
	"context"
	"dota2_fantasy/src/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

type AuthnPickRepo interface {
	Persist(pool *pgxpool.Pool, pick *model.AuthnPick) error
}

func NewAuthnPickRepo() AuthnPickRepo {
	return authnPickRepo{}
}

type authnPickRepo struct{}

func (r authnPickRepo) Persist(pool *pgxpool.Pool, pick *model.AuthnPick) error {
	_, err := pool.Exec(context.Background(), `INSERT INTO authn_pick(league_participant_id, timeframe_id, player_ids) VALUES ($1, $2, $3)`,
		pick.LeagueParticipantID, pick.TimeframeID, pick.PlayerIDs)

	return err
}
