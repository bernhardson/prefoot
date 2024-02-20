package app

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/bernhardson/prefoot/server/pkg/transport/rest"
)

func Run() {

	connConfig, err := pgx.ParseConfig("postgres://peterson:123@localhost/prefoot")
	if err != nil {
		log.Err(err).Msg("")
	}

	rest.Pool, err = pgxpool.New(context.Background(), connConfig.ConnString())
	if err != nil {
		log.Err(err).Msg("")
	}
	defer rest.Pool.Close()

	rest.StartServer()

}
