package postgres

import (
	"context"
	"fmt"

	"github.com/Edbeer/restapi/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPsqlDB(c *config.Config) (*pgxpool.Pool, error) {
	dataSourceName := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		c.Postgres.PostgresqlUser,
		c.Postgres.PostgresqlPassword,
		c.Postgres.PostgresqlHost,
		c.Postgres.PostgresqlPort,
		c.Postgres.PostgresqlDbname,
	)

	db, err := pgxpool.Connect(context.Background(), dataSourceName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(context.Background()); err != nil {
		return nil, err
	}

	return db, nil
}