package database

import (
	"context"

	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

const (
	insertTeam  = `INSERT INTO teams (id, name)VALUES ($1, $2)`
	insertVenue = `INSERT INTO venues (id, name, city) VALUES ($1, $2, $3)`
)

func InsertTeams(pool *pgxpool.Pool, teams *[]model.TeamVenue) {
	for _, t := range *teams {
		//insert venue
		_, err := pool.Exec(
			context.Background(),
			insertVenue,
			t.Venue.ID, t.Venue.Name, t.Venue.City,
		)
		if err != nil {
			log.Err(err).Msg("")
		}
		//insert team
		_, err = pool.Exec(
			context.Background(),
			insertTeam,
			t.Team.ID,
			t.Team.Name,
		)
		if err != nil {
			log.Err(err).Msg("")
		}
	}
}
