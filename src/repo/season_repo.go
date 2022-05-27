package repo

import (
	"context"
	"dota2_fantasy/src/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

type SeasonRepo interface {
	Persist(pool *pgxpool.Pool, roster *model.Season) error
}

func NewSeasonRepo() seasonRepo {
	return seasonRepo{}
}

type seasonRepo struct{}

func (r seasonRepo) Persist(pool *pgxpool.Pool, roster *model.Season) error {
	err := pool.QueryRow(context.Background(), `INSERT INTO season (converence_id, name, start_date, end_date) VALUES ($1, $2, $3, $4)`,
		roster.ConferenceID, roster.Name, roster.StartDate, roster.EndDate).Scan(&roster.ID)

	return err
}
