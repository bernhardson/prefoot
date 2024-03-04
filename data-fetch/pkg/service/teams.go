package service

import (
	"github.com/bernhardson/prefoot/data-fetch/pkg/database"
	"github.com/bernhardson/prefoot/data-fetch/pkg/fetch"
	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

func FetchAndInsertTeams(pool *pgxpool.Pool, league int, season int, logger zerolog.Logger) (*[]model.TeamVenue, error) {

	ts, err := fetch.GetTeams(league, season)
	if err != nil {
		return nil, err
	}

	database.InsertTeams(pool, ts)
	return ts, nil
}
