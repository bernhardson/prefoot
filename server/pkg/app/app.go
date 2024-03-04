package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/bernhardson/prefoot/server/pkg/transport/rest"
)

func Run() {

	// Create a context with the connection pool

	connConfig, err := pgxpool.ParseConfig("postgres://peterson:123@localhost/prefoot")
	if err != nil {
		log.Err(err).Msg("")
	}

	pool, err := pgxpool.New(context.Background(), connConfig.ConnString())
	if err != nil {
		log.Err(err).Msg("")
	}
	defer pool.Close()

	rest.Pool = pool
	//ctx := context.WithValue(context.Background(), "pgxpool", pool)

	rest.StartServer()

}
