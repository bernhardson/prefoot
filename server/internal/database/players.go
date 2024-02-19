package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Player struct {
	Player           int
	Fixture          int
	Season           int
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

func SelectPlayerStatistics(ids *[]int, pool *pgxpool.Pool) ([]*Player, error) {

	rows, err := pool.Query(context.Background(), "SELECT * FROM player_statistics WHERE player = ANY($1)", ids)
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
