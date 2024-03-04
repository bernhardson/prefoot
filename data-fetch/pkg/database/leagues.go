package database

import (
	"context"
	"fmt"

	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

const (
	insertLeague = `INSERT INTO leagues (id, name) VALUES ($1, $2);`
)

func InsertLeagues(pool *pgxpool.Pool, ls *[]model.LeagueData) (*[]int, *[]int, error) {
	insertedLeague := []int{}
	failedInsert := []int{}
	var err error
	for _, l := range *ls {
		_, err := pool.Exec(
			context.Background(),
			insertLeague,
			l.League.ID, l.League.Name)
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("leagueId_%d", l.League.ID))
			failedInsert = append(failedInsert, l.League.ID)
		} else {
			insertedLeague = append(insertedLeague, l.League.ID)
		}
	}
	return &insertedLeague, &failedInsert, err
}
