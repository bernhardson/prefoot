package service

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/bernhardson/prefoot/data-fetch/pkg/database"
	"github.com/bernhardson/prefoot/data-fetch/pkg/fetch"
)

func FetchAndInsertFixtures(pool *pgxpool.Pool, league, season int, logger zerolog.Logger) error {

	fixtures, err := fetch.GetMatches(league, season)
	if err != nil {
		log.Err(err).Msg("")
		return err
	}

	for _, f := range *fixtures {
		//insert league
		fd, err := fetch.GetFixtureDetail(f.Fixture.ID)
		if err != nil {
			return err
		}
		database.InsertFixtures(pool, fd, f, season, logger)
	}
	return nil
}
