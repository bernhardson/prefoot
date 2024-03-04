package fetch

import (
	"encoding/json"
	"errors"

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

func GetPlayers(league int, season int, paging int) (*[]model.Player, *model.Paging, error) {

	data, err := comm.GetHttpBodyRaw(playersURL, league, season, paging)
	if err != nil {
		return nil, nil, err
	}

	p := model.PlayerAPIResponse{}
	err = json.Unmarshal(data, &p)
	if err != nil {
		return nil, nil, err
	}
	return &p.Response, &p.Paging, nil
}

func GetPlayersByTeamId(teamId int) (*[]model.PlayerSquad, error) {
	data, err := comm.GetHttpBodyRaw(playersSquadURL, teamId)
	if err != nil {
		return nil, err
	}

	sqs := model.SquadResponse{}
	err = json.Unmarshal(data, &sqs)
	if err != nil {
		return nil, err
	}

	for _, sq := range sqs.Response {
		return &sq.Players, nil
	}
	return nil, errors.New("WTF")
}

func GetPlayerIdsByTeamId(teamId int) (*[]int, *[]model.PlayerSquad, error) {
	data, err := comm.GetHttpBodyRaw(playersSquadURL, teamId)
	if err != nil {
		return nil, nil, err
	}
	sqs := model.SquadResponse{}
	err = json.Unmarshal(data, &sqs)
	if err != nil {
		return nil, nil, err
	}

	res := []int{}

	for _, s := range sqs.Response {
		for _, p := range s.Players {
			res = append(res, p.ID)
		}
	}
	return &res, &sqs.Response[0].Players, nil
}
