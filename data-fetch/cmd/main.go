package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bernhardson/prefoot/data-fetch/pkg/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	connConfig, err := pgxpool.ParseConfig("postgres://peterson:123@localhost/prefoot")
	if err != nil {
		log.Err(err).Msg("")
	}

	pool, err := pgxpool.New(context.Background(), connConfig.ConnString())
	if err != nil {
		log.Err(err).Msg("")
	}

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.ErrorLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	season := 2023
	league := 78

	ils, ifs, err := service.FetchAndInsertLeagues(pool, logger)
	logger.Err(err).Msg(fmt.Sprintf("insert leagues: success=%v # failed=%v", *ils, *ifs))

	ts, err := service.FetchAndInsertTeams(pool, league, season, logger)
	logger.Err(err).Msg(fmt.Sprintf("insert teams: league:%d#season=%d", league, season))

	err = service.FetchAndInsertPlayers(pool, league, season, logger)
	logger.Err(err).Msg(fmt.Sprintf("insert players: league:%d#season=%d", league, season))

	err = service.FetchAndInsertFixtures(pool, league, season, logger)
	logger.Err(err).Msg(fmt.Sprintf("insert fixtures: league:%d#season=%d", league, season))

	ils, ifs, err = service.FetchAndInsertCoaches(pool, ts, logger)
	logger.Err(err).Msg(fmt.Sprintf("insert coaches: league:%d#season=%d#success=%v#failed=%v", league, season, *ils, *ifs))

}
