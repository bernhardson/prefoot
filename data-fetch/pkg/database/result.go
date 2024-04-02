package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	insertResult                  = `INSERT INTO "results" ("team", "league", "round", "season", "points", "goals_for", "goals_against", "modus", elapsed) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`
	selectResult                  = `SELECT * FROM results WHERE team=$1`
	selectResultByLeagueAndSeason = `SELECT * FROM results WHERE league=$1 AND season=$2`
)

type ResultModel struct {
	Pool *pgxpool.Pool
}

type ResultRow struct {
	Team         int `json:"team"`
	League       int `json:"league"`
	Round        int `json:"round"`
	Season       int `json:"season"`
	Points       int `json:"points"`
	GoalsFor     int `json:"goals_for"`
	GoalsAgainst int `json:"goals_against"`
	Modus        int `json:"modus"`
	Elapsed      int `json:"elapsed"`
}

func (sm *ResultModel) Insert(s *ResultRow) (int64, error) {
	row, err := sm.Pool.Exec(
		context.Background(),
		insertResult,
		s.Team, s.League, s.Round, s.Season, s.Points, s.GoalsFor, s.GoalsAgainst, s.Modus, s.Elapsed)
	return row.RowsAffected(), err
}

func (sm *ResultModel) Select(id int) (*ResultRow, error) {
	s := &ResultRow{}
	err := sm.Pool.QueryRow(context.Background(), selectResult, id).Scan(&s.Team, &s.League, &s.Round, &s.Season, &s.Points, &s.GoalsFor, &s.GoalsAgainst, &s.Modus, &s.Elapsed)
	return s, err
}

func (sm *ResultModel) SelectByLeagueAndSeason(league, season int) (*[]*ResultRow, error) {
	rows, err := sm.Pool.Query(
		context.Background(),
		selectResultByLeagueAndSeason, league, season)
	if err != nil {
		return nil, err
	}

	pls, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[ResultRow])
	return &pls, err
}
