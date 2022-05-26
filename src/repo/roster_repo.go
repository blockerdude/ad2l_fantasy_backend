package repo

import (
	"context"
	"dota2_fantasy/src/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

type RosterRepo interface {
	PersistRoster(pool *pgxpool.Pool, roster *model.Roster) error
	PersistRosterPlayer(pool *pgxpool.Pool, rosterPlayer *model.RosterPlayer) error
}

func NewRosterRepo() rosterRepo {
	return rosterRepo{}
}

type rosterRepo struct{}

func (r playerRepo) PersistRoster(pool *pgxpool.Pool, roster *model.Roster) error {
	err := pool.QueryRow(context.Background(), `INSERT INTO roster (season_id, team_id) VALUES ($1, $2)`,
		roster.SeasonID, roster.TeamID).Scan(&roster.ID)

	return err
}

func (r playerRepo) PersistRosterPlayer(pool *pgxpool.Pool, rosterPlayer *model.RosterPlayer) error {
	_, err := pool.Exec(context.Background(), `INSERT INTO roster_player (player_id, roster_id, position) VALUES ($1, $2, $3)`,
		rosterPlayer.PlayerID, rosterPlayer.RosterID, rosterPlayer.Position)

	return err
}
