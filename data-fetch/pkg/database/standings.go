package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	insertStanding = `INSERT INTO "standings" ("team", "league", "round", "season", "points", "goals_for", "goals_against", "modus") VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`
)

type StandingModel struct {
	Pool *pgxpool.Pool
}

type StandingRow struct {
	Team         int `json:"team"`
	League       int `json:"league"`
	Round        int `json:"round"`
	Season       int `json:"season"`
	Points       int `json:"points"`
	GoalsFor     int `json:"goals_for"`
	GoalsAgainst int `json:"goals_against"`
	Modus        int `json:"modus"`
}

func (lm *StandingModel) Insert(s *StandingRow) (int64, error) {
	row, err := lm.Pool.Exec(
		context.Background(),
		insertStanding,
		s.Team, s.League, s.Round, s.Season, s.Points, s.GoalsFor, s.GoalsAgainst, s.Modus)
	return row.RowsAffected(), err
}
