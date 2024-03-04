package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Player struct {
	ID               int
	Team             int
	Season           int
	FirstName        string
	LastName         string
	BirthPlace       string
	BirthCountry     string
	BirthDate        string
	Fixture          int
	Minutes          int
	Position         string
	Rating           float64
	Captain          bool
	Substitute       bool
	ShotsTotal       int
	ShotsOn          int
	GoalsScored      int
	GoalsAssisted    int
	PassesTotal      int
	PassesKey        int
	Accuracy         int
	Tackles          int
	Block            int
	Interceptions    int
	DuelsTotal       int
	DuelsWon         int
	DribblesTotal    int
	DribblesWon      int
	Yellow           int
	Red              int
	PenaltyWon       int
	PenaltyCommitted int
	PenaltyScored    int
	PenaltyMissed    int
	PenaltySaved     int
	Saves            int
}

type PlayerStatistic struct {
	Player           int     `db:"player"`
	Fixture          int     `db:"fixture"`
	Team             int     `db:"team"`
	Season           int     `db:"season"`
	Minutes          int     `db:"minutes"`
	Position         string  `db:"position"`
	Rating           float64 `db:"rating"`
	Captain          bool    `db:"captain"`
	Substitute       bool    `db:"substitute"`
	ShotsTotal       int     `db:"shots_total"`
	ShotsOn          int     `db:"shots_on"`
	GoalsScored      int     `db:"goals_scored"`
	GoalsAssisted    int     `db:"goals_assisted"`
	PassesTotal      int     `db:"passes_total"`
	PassesKey        int     `db:"passes_key"`
	Accuracy         int     `db:"accuracy"`
	Tackles          int     `db:"tackles"`
	Block            int     `db:"block"`
	Interceptions    int     `db:"interceptions"`
	DuelsTotal       int     `db:"duels_total"`
	DuelsWon         int     `db:"duels_won"`
	DribblesTotal    int     `db:"dribbles_total"`
	DribblesWon      int     `db:"dribbles_won"`
	Yellow           int     `db:"yellow"`
	Red              int     `db:"red"`
	PenaltyWon       int     `db:"penalty_won"`
	PenaltyCommitted int     `db:"penalty_committed"`
	PenaltyScored    int     `db:"penalty_scored"`
	PenaltyMissed    int     `db:"penalty_missed"`
	PenaltySaved     int     `db:"penalty_saved"`
	Saves            int     `db:"saves"`
}

func SelectPlayerStatistics(ids *[]int, pool *pgxpool.Pool) ([]*PlayerStatistic, error) {

	rows, err := pool.Query(context.Background(), "SELECT * FROM player_statistics WHERE player = ANY($1)", ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		log.Err(err).Msg("")
	}
	players, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[PlayerStatistic])
	if err != nil {
		return nil, err
	}
	return players, nil
}

func SelectPlayers(ids *[]int, pool *pgxpool.Pool) ([]*Player, error) {

	rows, err := pool.Query(context.Background(), "SELECT * FROM players WHERE player = ANY($1)", ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		log.Err(err).Msg("")
	}
	players, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Player])
	if err != nil {
		return nil, err
	}
	return players, nil
}

func SelectPlayersByTeamId(id int, pool *pgxpool.Pool) ([]*Player, error) {

	rows, err := pool.Query(context.Background(), "SELECT * FROM players WHERE team = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		log.Err(err).Msg("")
	}
	players, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Player])
	if err != nil {
		return nil, err
	}
	return players, nil
}

func SelectPlayersAndStatisticsByTeamId(id int, pool *pgxpool.Pool) (*[]*Player, error) {

	// Perform a JOIN query to retrieve data from both tables
	rows, err := pool.Query(
		context.Background(),
		`SELECT p.*, ps.*
        FROM players p
        JOIN player_statistics ps ON p.id = ps.player
		WHERE p.team = $1`, id)

	if err != nil {
		return nil, err
	}

	pls, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Player])
	return &pls, err
}
