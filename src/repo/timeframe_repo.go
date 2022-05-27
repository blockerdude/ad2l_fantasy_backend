package repo

import (
	"context"
	"dota2_fantasy/src/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TimeframeRepo interface {
	Persist(pool *pgxpool.Pool, roster *model.Timeframe) error
}

func NewTimeframeRepo() timeframeRepo {
	return timeframeRepo{}
}

type timeframeRepo struct{}

func (r timeframeRepo) Persist(pool *pgxpool.Pool, timeframe *model.Timeframe) error {
	err := pool.QueryRow(context.Background(), `INSERT INTO timeframe (season_id, open, close_date, name) VALUES ($1, $2, $3, $4)`,
		timeframe.SeasonID, timeframe.Open, timeframe.CloseDate, timeframe.Name).Scan(&timeframe.ID)

	return err
}
