package repo

import (
	"context"
	"dota2_fantasy/src/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

type LeagueParticipantRepo interface {
	Persist(pool *pgxpool.Pool, participant *model.LeagueParticipant) error
	GetParticipant(pool *pgxpool.Pool, authnID int, leagueObjectID string) (*model.LeagueParticipant, error)
}

func NewLeagueParticipantRepo() LeagueParticipantRepo {
	return leagueParticipantRepo{}
}

type leagueParticipantRepo struct{}

func (r leagueParticipantRepo) Persist(pool *pgxpool.Pool, participant *model.LeagueParticipant) error {
	err := pool.QueryRow(context.Background(), `INSERT INTO leage_participant (league_id, authn_id, league_admin, paid) VALUES ($1, $2, $3, $4)`,
		participant.LeagueID, participant.AuthnID, participant.LeagueAdmin, participant.Paid).Scan(&participant.ID)

	return err
}

func (r leagueParticipantRepo) GetParticipant(pool *pgxpool.Pool, authnID int, leagueObjectID string) (*model.LeagueParticipant, error) {
	lp := &model.LeagueParticipant{AuthnID: authnID}

	query := `SELECT lp.id, lp.league_id, lp.league_admin, lp.paid FROM league_participant as lp
				JOIN league ON league.id = lp.league_id WHERE league.object_id = $1 AND lp.authn_id = $2`

	err := pool.QueryRow(context.Background(), query, leagueObjectID, authnID).
		Scan(&lp.ID, lp.LeagueID, lp.LeagueAdmin, lp.Paid)

	return lp, err
}
