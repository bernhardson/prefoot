package fetch

import (
	"github.com/bernhardson/prefoot/data-fetch/pkg/comm"
	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
	"github.com/mitchellh/mapstructure"
)

const (
	teamsURL = "https://api-football-v1.p.rapidapi.com/v3/teams?league=%d&season=%d"
)

func GetTeams(league int, season int) (*[]model.TeamVenue, error) {
	data, err := comm.GetHttpBody(teamsURL, league, season)
	if err != nil {
		return nil, err
	}
	tv := &[]model.TeamVenue{}
	mapstructure.Decode(data["response"], tv)

	return tv, nil
}
