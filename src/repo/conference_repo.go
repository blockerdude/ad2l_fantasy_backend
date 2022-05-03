package repo

import (
	"context"
	"dota2_fantasy/src/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ConferenceRepo interface {
	GetAllConferences(conn *pgxpool.Pool) ([]model.Conference, error)
}

func NewConferenceRepo() ConferenceRepo {
	return conferenceRepo{}
}

type conferenceRepo struct{}

func (c conferenceRepo) GetAllConferences(pool *pgxpool.Pool) ([]model.Conference, error) {

	rows, err := pool.Query(context.Background(), `select id, name, description from conference`)
	if err != nil {
		return nil, err
	}

	conferences := make([]model.Conference, 0)
	for rows.Next() {

		conf := model.Conference{}
		err = rows.Scan(&conf.ID, &conf.Name, &conf.Description)
		if err != nil {
			return nil, err
		}

		conferences = append(conferences, conf)
	}

	return conferences, nil
}
