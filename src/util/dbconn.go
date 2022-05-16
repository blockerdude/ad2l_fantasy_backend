package util

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

var dbPool *pgxpool.Pool

type DBConnection interface {
	EstablishConnection() error
	GetPool() *pgxpool.Pool
}

func NewDBConnection(secrets Secrets) DBConnection {
	return dbConnection{
		secrets: secrets,
	}
}

type dbConnection struct {
	secrets Secrets
}

func (dbc dbConnection) EstablishConnection() error {
	pool, err := pgxpool.Connect(context.Background(), dbc.secrets.DBConnectionString)

	if err != nil {
		return err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return err
	}

	dbPool = pool

	return nil
}

func (dbc dbConnection) GetPool() *pgxpool.Pool {
	return dbPool
}
