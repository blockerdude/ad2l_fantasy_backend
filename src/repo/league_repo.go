package repo

import (
	"context"
	"dota2_fantasy/src/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

type LeagueRepo interface {
	Persist(pool *pgxpool.Pool, league *model.League) error
}

func NewLeagueRepo() leagueRepo {
	return leagueRepo{}
}

type leagueRepo struct{}

func (r leagueRepo) Persist(pool *pgxpool.Pool, league *model.League) error {
	err := pool.QueryRow(context.Background(), `INSERT INTO league (object_id, season_id, name, description) VALUES ($1, $2, $3, $4)`,
		league.ObjectID, league.SeasonID, league.Name, league.Description).Scan(&league.ID)

	return err
}
