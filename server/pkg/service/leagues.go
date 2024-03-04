package service

import (
	"encoding/json"

	"github.com/bernhardson/prefoot/data-fetch/pkg/comm"
	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
)

const (
	getLeaguesURL  = "https://api-football-v1.p.rapidapi.com/v3/leagues"
	getStandingURL = "https://api-football-v1.p.rapidapi.com/v3/standings?season=%d&league=%d"
)

func GetStandings(league int, season int) (*[]model.StandingsEntry, error) {

	data, err := comm.GetHttpBodyRaw(getStandingURL, league, season)

	resp := model.StandingsResponse{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Response, nil

}
