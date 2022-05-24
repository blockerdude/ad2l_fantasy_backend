package repo

import (
	"context"
	"dota2_fantasy/src/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TeamRepo interface {
	GetAllTeams(pool *pgxpool.Pool) ([]model.Team, error)
	FindByName(pool *pgxpool.Pool, name string) (*model.Team, error)
	Persist(pool *pgxpool.Pool, team *model.Team) error
}

func NewTeamRepo() TeamRepo {
	return teamRepo{}
}

type teamRepo struct{}

func (r teamRepo) GetAllTeams(pool *pgxpool.Pool) ([]model.Team, error) {

	rows, err := pool.Query(context.Background(), `SELECT id, conference_id, name FROM team`)
	if err != nil {
		return nil, err
	}

	teams := make([]model.Team, 0)
	for rows.Next() {

		team := model.Team{}
		err = rows.Scan(&team.ID, &team.ConferenceID, &team.Name)
		if err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	return teams, nil
}

func (r teamRepo) FindByName(pool *pgxpool.Pool, name string) (*model.Team, error) {
	team := &model.Team{Name: name}
	err := pool.QueryRow(context.Background(), `SELECT id, conference_id FROM team WHERE name = $1`, name).
		Scan(&team.ID, &team.ConferenceID)

	return team, err
}

func (r teamRepo) Persist(pool *pgxpool.Pool, team *model.Team) error {
	err := pool.QueryRow(context.Background(), `INSERT INTO team (conference_id, name) VALUES ($1, $2) RETURNING id`,
		team.ConferenceID, team.Name).Scan(&team.ID)

	return err
}
