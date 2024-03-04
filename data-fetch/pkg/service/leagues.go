package service

import (
	"github.com/bernhardson/prefoot/data-fetch/pkg/database"
	"github.com/bernhardson/prefoot/data-fetch/pkg/fetch"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func FetchAndInsertLeagues(pool *pgxpool.Pool, logger zerolog.Logger) (*[]int, *[]int, error) {
	ls, err := fetch.GetLeagues()
	if err != nil {
		log.Err(err).Msg("")
	}
	insertedLeague, failedInsert, err := database.InsertLeagues(pool, ls)
	if err != nil {
		return nil, nil, err
	}
	return insertedLeague, failedInsert, nil
}

func GetStandings (league int, season int ,logger zerolog.Logger){
	
}
