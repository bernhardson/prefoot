package database

import (
	"context"

	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	insertFixture = `INSERT INTO fixtures (id, league, referee, timezone, timestamp, venue, season, home_team, away_team,
						home_goals, away_goals, home_goals_half, away_goals_half)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);`
)

func InsertFixture(pool *pgxpool.Pool, f model.Match, fd model.FixtureDetail, season int) error {

	_, err := pool.Exec(
		context.Background(),
		insertFixture,
		f.Fixture.ID, f.League.ID, f.Fixture.Referee, f.Fixture.Timezone,
		f.Fixture.Timestamp, f.Fixture.Venue.ID, season, f.Teams.Home.ID,
		f.Teams.Away.ID, f.Goals.Home, f.Goals.Away, fd.Score.Halftime.Home,
		fd.Score.Halftime.Away)

	return err
}
