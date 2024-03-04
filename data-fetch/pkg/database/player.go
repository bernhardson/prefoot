package database

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	insertPlayerStatistics       = `INSERT INTO player_statistics (player, fixture, team, season, minutes, position, rating, captain, substitute, shots_total, shots_on, goals_scored, goals_assisted, passes_total, passes_key, accuracy, tackles, block, interceptions, duels_total, duels_won, dribbles_total,dribbles_won, yellow, red, penalty_won, penalty_committed, penalty_scored, penalty_missed, penalty_saved, saves)	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,$18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31)`
	insertPlayer                 = `INSERT INTO players (id, team, season, firstname, lastname, birthplace, birthcountry, birthdate) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	insertPlayerStatisticsSeason = `INSERT INTO player_statistics_season ("player","season", "team", "minutes", "position", "rating","captain", "games", "lineups", "shots_total", "shots_on", "goals_scored", "goals_assisted", "passes_total", "passes_key","accuracy", "tackles", "block", "interceptions", "duels_total", "duels_won", "dribbles_total", "dribbles_won", "yellow", "red", "penalty_won","penalty_committed", "penalty_scored", "penalty_missed", "penalty_saved", "saves") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31)`
)

func InsertPlayerStatistic(pool *pgxpool.Pool, p *model.PlayerStatisticsFD, ps *model.PlayerStatisticsDetails, f *model.FixtureMeta, team, season int, logger zerolog.Logger) error {

	_, err := pool.Exec(
		context.Background(),
		insertPlayerStatistics,
		p.Player.ID, f.ID, team, season, ps.Games.Minutes, ps.Games.Position, ps.Games.Rating,
		ps.Games.Captain, ps.Games.Substitute, ps.Shots.Total, ps.Shots.On, ps.Goals.Total,
		ps.Goals.Assists, ps.Passes.Total, ps.Passes.Key, ps.Passes.Accuracy, ps.Tackles.Total,
		ps.Tackles.Blocks, ps.Tackles.Interceptions, ps.Duels.Total, ps.Duels.Won,
		ps.Dribbles.Attempts, ps.Dribbles.Success, ps.Cards.Yellow, ps.Cards.Red, ps.Penalty.Won,
		ps.Penalty.Commited, ps.Penalty.Scored, ps.Penalty.Missed, ps.Penalty.Saved, ps.Goals.Saves,
	)
	return err
}

func InsertPlayers(players *[]model.Player, pool *pgxpool.Pool, logger zerolog.Logger) ([]int, []int) {
	var failedP, failedS []int
	for _, p := range *players {
		// player statistics only has one entry so there will be just one insert to player table
		for _, s := range p.Statistics {
			_, err := pool.Exec(
				context.Background(),
				insertPlayer,
				p.PlayerDetails.ID, s.Team.ID, s.League.Season, p.PlayerDetails.FirstName,
				p.PlayerDetails.LastName, p.PlayerDetails.Birth.Place,
				p.PlayerDetails.Birth.Country, p.PlayerDetails.Birth.Date)
			if err != nil {
				if !strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
					failedP = append(failedP, p.PlayerDetails.ID)
					log.Debug().Msg(err.Error())
				} else {
					logger.Err(err).Msg(fmt.Sprintf("player id %d", p.PlayerDetails.ID))
				}

			} else {
				logger.Debug().Msg(fmt.Sprintf("inserted player id %d", p.PlayerDetails.ID))
			}
			rating, err := strconv.ParseFloat(s.Games.Rating, 32)
			if err != nil {
				rating = 0
				logger.Info().Msg(err.Error())
			}
			_, err = pool.Exec(
				context.Background(),
				insertPlayerStatisticsSeason,
				p.PlayerDetails.ID, s.League.Season, s.Team.ID, s.Games.Minutes,
				s.Games.Position, rating, s.Games.Captain, s.Games.Appearances,
				s.Games.Lineups, s.Shots.Total, s.Shots.On, s.Goals.Total,
				s.Goals.Assists, s.Passes.Total, s.Passes.Key, s.Passes.Accuracy,
				s.Tackles.Total, s.Tackles.Blocks, s.Tackles.Interceptions, s.Duels.Total,
				s.Duels.Won, s.Dribbles.Attempts, s.Dribbles.Success, s.Cards.Yellow, s.Cards.Red,
				s.Penalty.Won, s.Penalty.Committed, s.Penalty.Scored, s.Penalty.Missed,
				s.Penalty.Saved, s.Goals.Saves,
			)
			if err != nil {
				logger.Err(err).Msg(fmt.Sprintf("failed inserting player season statistics %d", p.PlayerDetails.ID))
				failedS = append(failedS, p.PlayerDetails.ID)
			} else {
				logger.Debug().Msg(fmt.Sprintf("inserted player season statistics %d", p.PlayerDetails.ID))
			}
		}
	}
	return failedP, failedS
}

