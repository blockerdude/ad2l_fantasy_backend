package repo

import (
	"context"
	"dota2_fantasy/src/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PlayerScoreRepo interface {
	Persist(pool *pgxpool.Pool, score *model.PlayerScore) error
}

func NewPlayerScoreRepo() PlayerScoreRepo {
	return playerScoreRepo{}
}

type playerScoreRepo struct{}

func (r playerScoreRepo) Persist(pool *pgxpool.Pool, score *model.PlayerScore) error {

	// TODO: will need to test that nils are inserted as nulls appropriately for total and substitute_id

	_, err := pool.Exec(context.Background(), `INSERT INTO player_score
		(timeframe_id, player_id, match_id, substitute_id, total, kills, deaths, last_hits, denies, teamfight_participation,
		gpm, tower_kills, rosh_kills, obs_placed, camps_stacked, runes_taken, first_blood, stun_seconds)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`,
		score.TimeframeID, score.PlayerID, score.MatchID, score.SubstituteID, score.Total, score.Kills,
		score.Deaths, score.LastHits, score.Denies, score.TeamfightParticipation, score.GoldPerMinute, score.TowerKills,
		score.RoshKills, score.ObserversPlaced, score.CampsStacked, score.RunesTaken, score.FirstBlood, score.StunSeconds)

	return err
}
