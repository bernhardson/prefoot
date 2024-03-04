package service

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/bernhardson/prefoot/data-fetch/pkg/database"
	"github.com/bernhardson/prefoot/data-fetch/pkg/fetch"
	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
)

func FetchAndInsertCoaches(pool *pgxpool.Pool, teams *[]model.TeamVenue, logger zerolog.Logger) (*[]int, *[]int, error) {
	insertSuccess := []int{}
	insertFail := []int{}

	for _, t := range *teams {
		//insert coach
		cs, err := fetch.GetCoach(t.Team.ID)
		if err != nil {
			logger.Err(err).Msg(fmt.Sprintf("could not get for team id %d", t.Team.ID))
		}

		ins, ifs := database.InsertCoaches(pool, cs, logger)
		insertSuccess = append(insertSuccess, *ins...)
		insertFail = append(insertFail, *ifs...)

	}
	return &insertSuccess, &insertFail, nil
}
