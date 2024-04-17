package team

import (
	"encoding/json"

	"github.com/bernhardson/prefoot/data-fetch/pkg/comm"
	"github.com/bernhardson/prefoot/data-fetch/pkg/shared"
)

const (
	teamsURL = "https://api-football-v1.p.rapidapi.com/v3/teams?league=%d&season=%d"
)

type TeamsResponse struct {
	Get        string        `json:"get"`
	Parameters Parameters    `json:"parameters"`
	Errors     []interface{} `json:"errors"`
	Results    int           `json:"results"`
	Paging     shared.Paging `json:"paging"`
	TeamVenues []TeamVenue   `json:"response"`
}

type TeamVenue struct {
	Team  Team  `json:"team"`
	Venue Venue `json:"venue"`
}

type Team struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Country  string `json:"country"`
	Founded  int    `json:"founded"`
	National bool   `json:"national"`
	Logo     string `json:"logo"`
}

type Venue struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	City     string `json:"city"`
	Capacity int    `json:"capacity"`
	Surface  string `json:"surface"`
	Image    string `json:"image"`
}

type Parameters struct {
	League string `json:"league"`
	Season string `json:"season"`
}

func GetTeams(league int, season int) (*TeamsResponse, error) {

	data, err := comm.GetHttpBody(teamsURL, league, season)
	if err != nil {
		return nil, err
	}

	tv := &TeamsResponse{}
	err = json.Unmarshal(data, tv)
	if err != nil {
		return nil, err
	}

	return tv, nil
}
