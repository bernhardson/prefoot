package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bernhardson/prefoot/data-fetch/pkg/comm"
	"github.com/bernhardson/prefoot/data-fetch/pkg/fetch"
)

const (
	getLeaguesURL  = "https://api-football-v1.p.rapidapi.com/v3/leagues"
	getStandingURL = "https://api-football-v1.p.rapidapi.com/v3/standings?league=%d&season=%d"
)

type Fixture struct {
	ID            int    `json:"id"`
	League        int    `json:"league"`
	Round         int    `json:"round"`
	Referee       string `json:"referee"`
	Timezone      string `json:"timezone"`
	Timestamp     int64  `json:"timestamp"`
	Venue         int    `json:"venue"`
	Season        int    `json:"season"`
	HomeTeam      int    `json:"home_team"`
	AwayTeam      int    `json:"away_team"`
	HomeGoals     int    `json:"home_goals"`
	AwayGoals     int    `json:"away_goals"`
	HomeGoalsHalf int    `json:"home_goals_half"`
	AwayGoalsHalf int    `json:"away_goals_half"`
}

func (app *application) SelectFixturesByRound(w http.ResponseWriter, r *http.Request) {
	round, err := strconv.Atoi(r.URL.Query().Get("round"))
	if err != nil {
		app.clientError(w, http.StatusNotAcceptable)
	}
	rows, err := app.repo.Fixture.SelectFixturesByRound(round)
	if err != nil {
		app.serverError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(rows)
}

func GetStandings(league int, season int) (*[]fetch.StandingsTeam, error) {

	data, err := comm.GetHttpBodyRaw(getStandingURL, league, season)
	if err != nil {
		return nil, err
	}
	resp := fetch.StandingsResponse{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	for _, s := range resp.Response {
		for _, st := range s.League.Standings {
			return &st, nil
		}
	}
	return nil, err
}
