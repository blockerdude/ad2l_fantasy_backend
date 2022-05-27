package repo

import (
	"context"
	"dota2_fantasy/src/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

type LeagueParticipantsRepo interface {
	Persist(pool *pgxpool.Pool, participant *model.LeagueParticipant) error
}

func NewLeagueParticipantsRepo() LeagueParticipantsRepo {
	return leagueParticipantsRepo{}
}

type leagueParticipantsRepo struct{}

func (r leagueParticipantsRepo) Persist(pool *pgxpool.Pool, participant *model.LeagueParticipant) error {
	err := pool.QueryRow(context.Background(), `INSERT INTO leage_participant (league_id, authn_id, league_admin, paid) VALUES ($1, $2, $3, $4)`,
		participant.LeagueID, participant.AuthnID, participant.LeagueAdmin, participant.Paid).Scan(&participant.ID)

	return err
}
