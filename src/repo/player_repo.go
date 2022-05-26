package repo

import (
	"context"
	"dota2_fantasy/src/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PlayerRepo interface {
	Persist(pool *pgxpool.Pool, player *model.Player) error
}

func NewPlayerRepo() playerRepo {
	return playerRepo{}
}

type playerRepo struct{}

func (r playerRepo) Persist(pool *pgxpool.Pool, player *model.Player) error {
	err := pool.QueryRow(context.Background(), `INSERT INTO player (steam_id, name) VALUES ($1, $2) RETURNING id`,
		player.SteamID, player.Name).Scan(&player.ID)

	return err
}
