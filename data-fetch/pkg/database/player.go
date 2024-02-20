package database

import (
	"context"

	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	insertPlayerStatistics = "INSERT INTO player_statistics " +
		"(p, fixture, season, minutes, position, rating, captain, substitute, " +
		"shots_total, shots_on, goals_scored, goals_assisted, passes_total, passes_key, " +
		"accuracy, tackles, block, interceptions, duels_total, duels_won, dribbles_total, " +
		"dribbles_won, yellow, red, penalty_won, penalty_committed, penalty_scored, " +
		"penalty_missed, penalty_saved, saves)" +
		"VALUES " +
		"($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, " +
		"$18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30)"
)

func InsertPlayerStatistic(pool *pgxpool.Pool, p *model.PlayerStatisticsFD, ps *model.PlayerStatisticsDetails, f *model.Fixture, season int) error {

	_, err := pool.Exec(
		context.Background(),
		insertPlayerStatistics,
		p.Player.ID, f.ID, season, ps.Games.Minutes, ps.Games.Position, ps.Games.Rating,
		ps.Games.Captain, ps.Games.Substitute, ps.Shots.Total, ps.Shots.On, ps.Goals.Total,
		ps.Goals.Assists, ps.Passes.Total, ps.Passes.Key, ps.Passes.Accuracy, ps.Tackles.Total,
		ps.Tackles.Blocks, ps.Tackles.Interceptions, ps.Duels.Total, ps.Duels.Won,
		ps.Dribbles.Attempts, ps.Dribbles.Success, ps.Cards.Yellow, ps.Cards.Red, ps.Penalty.Won,
		ps.Penalty.Commited, ps.Penalty.Scored, ps.Penalty.Missed, ps.Penalty.Saved, ps.Goals.Saves,
	)

	return err
}
