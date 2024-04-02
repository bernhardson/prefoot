package fetch

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bernhardson/prefoot/data-fetch/pkg/comm"
)

const (
	getLeaguesURL  = "https://api-football-v1.p.rapidapi.com/v3/leagues"
	getLeagueURL   = "https://api-football-v1.p.rapidapi.com/v3/leagues?id=%d"
	getStandingURL = "https://api-football-v1.p.rapidapi.com/v3/standings?league=%d&season=%d"
)

type LeagueResponse struct {
	Get        string        `json:"get"`
	Parameters LeagueParams  `json:"parameters"`
	Errors     []interface{} `json:"errors"`
	Results    int           `json:"results"`
	Paging     Paging        `json:"paging"`
	Response   []LeagueData  `json:"response"`
}

type LeagueParams struct {
	ID string `json:"id"`
}

type LeagueData struct {
	League  League   `json:"league"`
	Country Country  `json:"country"`
	Seasons []Season `json:"seasons"`
}

type League struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Logo    string `json:"logo"`
	Flag    string `json:"flag"`
	Season  int    `json:"season"`
	Type    string `json:"type"`
	Round   string `json:"round"`
}

type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Flag string `json:"flag"`
}

type Season struct {
	Year     int            `json:"year"`
	Start    string         `json:"start"`
	End      string         `json:"end"`
	Current  bool           `json:"current"`
	Coverage SeasonCoverage `json:"coverage"`
}

type SeasonCoverage struct {
	Fixtures    SeasonCoverageFixtures `json:"fixtures"`
	Standings   bool                   `json:"standings"`
	Players     bool                   `json:"players"`
	TopScorers  bool                   `json:"top_scorers"`
	TopAssists  bool                   `json:"top_assists"`
	TopCards    bool                   `json:"top_cards"`
	Injuries    bool                   `json:"injuries"`
	Predictions bool                   `json:"predictions"`
	Odds        bool                   `json:"odds"`
}

type SeasonCoverageFixtures struct {
	Events             bool `json:"events"`
	Lineups            bool `json:"lineups"`
	StatisticsFixtures bool `json:"statistics_fixtures"`
	StatisticsPlayers  bool `json:"statistics_players"`
}

func GetLeagues() (*LeagueResponse, error) {

	data, err := comm.GetHttpBodyRaw(getLeaguesURL)
	if err != nil {
		return nil, err
	}

	l := &LeagueResponse{}
	err = json.Unmarshal(data, l)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func GetLeague(id int) (*LeagueResponse, error) {

	data, err := comm.GetHttpBodyRaw(fmt.Sprintf(getLeagueURL, id))
	if err != nil {
		return nil, err
	}

	l := &LeagueResponse{}
	err = json.Unmarshal(data, l)
	if err != nil {
		return nil, err
	}

	return l, nil
}

type StandingsResponse struct {
	Get        string           `json:"get"`
	Parameters StandingsParams  `json:"parameters"`
	Errors     interface{}      `json:"errors"`
	Results    int              `json:"results"`
	Paging     Paging           `json:"paging"`
	Response   []StandingsEntry `json:"response"`
}

type StandingsParams struct {
	League string `json:"league"`
	Season string `json:"season"`
}

type StandingsEntry struct {
	League StandingsLeague `json:"league"`
}

type StandingsLeague struct {
	ID        int               `json:"id"`
	Name      string            `json:"name"`
	Country   string            `json:"country"`
	Logo      string            `json:"logo"`
	Flag      string            `json:"flag"`
	Season    int               `json:"season"`
	Standings [][]StandingsTeam `json:"standings"`
}

type StandingsTeam struct {
	Rank        int                 `json:"rank"`
	Team        StandingsTeamDetail `json:"team"`
	Points      int                 `json:"points"`
	GoalsDiff   int                 `json:"goalsDiff"`
	Group       string              `json:"group"`
	Form        string              `json:"form"`
	Status      string              `json:"status"`
	Description string              `json:"description"`
	All         StandingsStats      `json:"all"`
	Home        StandingsStats      `json:"home"`
	Away        StandingsStats      `json:"away"`
	Update      string              `json:"update"`
}

type StandingsTeamDetail struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

type StandingsStats struct {
	Played int                `json:"played"`
	Win    int                `json:"win"`
	Draw   int                `json:"draw"`
	Lose   int                `json:"lose"`
	Goals  StandingsGoalStats `json:"goals"`
}

type StandingsGoalStats struct {
	For     int `json:"for"`
	Against int `json:"against"`
}

func GetStanding(league int, season int) (*StandingsEntry, error) {

	data, err := comm.GetHttpBodyRaw(getStandingURL, league, season)
	if err != nil {
		return nil, err
	}
	resp := &StandingsResponse{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	_, ok := resp.Errors.(map[string]interface{})
	if ok {
		return nil, errors.New(fmt.Sprint(resp.Errors))
	}
	return &resp.Response[0], nil
}
