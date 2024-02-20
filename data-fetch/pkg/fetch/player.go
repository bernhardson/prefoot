package pkg

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jackc/pgx"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"

	"github.com/bernhardson/prefoot/data-fetch/pkg/comm"
	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
)

const (
	playersURL                   = "https://api-football-v1.p.rapidapi.com/v3/players?league=%d&season=%d&page=%d"
	playersSquadURL              = "https://api-football-v1.p.rapidapi.com/v3/players/squads?team=%d"
	insertPlayer                 = `INSERT INTO players (id, team, season, firstname, lastname) VALUES ($1, $2, $3, $4, $5)`
	insertPlayerStatisticsSeason = `INSERT INTO player_statistics_season ("player", "season", "team","minutes", "position", "rating", "captain", "games","lineups", "shots_total", "shots_on",
									"goals_scored", "goals_assisted", "passes_total", "passes_key", "accuracy", "tackles", "block", "interceptions", "duels_total", "duels_won", 
									"dribbles_total", "dribbles_won", "yellow", "red", "penalty_won", "penalty_committed", "penalty_scored", "penalty_missed", "penalty_saved", "saves")
        							VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31)`
)

func GetPlayers(league int, season int) (*[]model.Player, error) {
	ret := []model.Player{}
	for i := 1; i < 100; i++ {
		data, err := comm.GetHttpBody(playersURL, season, 1)
		if err != nil {
			return nil, err
		}

		resp := &model.PlayerAPIResponse{}
		mapstructure.Decode(data, &resp)

		p := []model.Player{}
		err = mapstructure.Decode(data["response"], &p)
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("%d", i))
		}
		ret = append(ret, p...)
		//TODO remove true
		if resp.Paging.Current == resp.Paging.Total || true {
			break
		}
	}
	return &ret, nil
}

func GetPlayersByTeamId(teamId int) (*[]model.PlayerSquad, error) {
	data, err := comm.GetHttpBody(playersSquadURL, teamId)
	if err != nil {
		return nil, err
	}
	sqs := []model.Squad{}
	err = mapstructure.Decode(data["response"], &sqs)
	if err != nil {
		return nil, err
	}
	for _, sq := range sqs {
		return &sq.Players, nil
	}
	return nil, errors.New("WTF")
}

func GetPlayerIdsByTeamId(teamId int) (*[]int, *[]model.PlayerSquad, error) {
	data, err := comm.GetHttpBody(playersSquadURL, teamId)
	if err != nil {
		return nil, nil, err
	}
	ps := []model.Squad{}
	err = mapstructure.Decode(data["response"], &ps)
	if err != nil {
		return nil, nil, err
	}

	res := []int{}

	for _, s := range ps {
		for _, p := range s.Players {
			res = append(res, p.ID)
		}
	}
	return &res, &ps[0].Players, nil
}

func InsertPlayers(players *[]model.Player, conn *pgx.Conn) {
	for _, p := range *players {
		// player statistics only has one entry so there will be just one insert to player table
		for _, s := range p.Statistics {
			_, err := conn.Exec(
				insertPlayer,
				p.PlayerDetails.ID, s.Team.ID, s.League.Season, p.PlayerDetails.FirstName, p.PlayerDetails.LastName)
			if err != nil {
				log.Err(err).Msg(fmt.Sprintf("player id %d", p.PlayerDetails.ID))
			} else {
				log.Info().Msg(fmt.Sprintf("player id %d", p.PlayerDetails.ID))
			}
			rating, err := strconv.ParseFloat(*s.Games.Rating, 32)
			if err != nil {
				rating = 0
				log.Err(err).Msg("")
			}
			_, err = conn.Exec(
				insertPlayerStatisticsSeason,
				p.PlayerDetails.ID, s.League.Season, s.Team.ID, s.Games.Minutes, s.Games.Position, rating,
				s.Games.Captain, s.Games.Appearances, s.Games.Lineups, s.Shots.Total, s.Shots.On, s.Goals.Total,
				s.Goals.Assists, s.Passes.Total, s.Passes.Key, s.Passes.Accuracy, s.Tackles.Total,
				s.Tackles.Blocks, s.Tackles.Interceptions, s.Duels.Total, s.Duels.Won,
				s.Dribbles.Attempts, s.Dribbles.Success, s.Cards.Yellow, s.Cards.Red,
				s.Penalty.Won, s.Penalty.Committed, s.Penalty.Scored, s.Penalty.Missed, s.Penalty.Saved, s.Goals.Saves,
			)
			if err != nil {
				log.Err(err).Msg(fmt.Sprintf("failed inserting player season statistics %d", p.PlayerDetails.ID))
			} else {
				log.Debug().Msg(fmt.Sprintf("inserted player %d", p.PlayerDetails.ID))
			}
		}
	}
}
