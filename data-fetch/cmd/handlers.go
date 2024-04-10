package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bernhardson/prefoot/data-fetch/pkg/database"
)

func (app *application) getRounds(w http.ResponseWriter, r *http.Request) {

	league, err := strconv.Atoi(r.URL.Query().Get("league"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	season, err := strconv.Atoi(r.URL.Query().Get("season"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	ts, err := strconv.ParseInt(r.URL.Query().Get("ts"), 10, 0)
	var res interface{}

	if err != nil {
		app.serverError(w, err)
	}

	if err == nil {
		res, err = app.repo.Fixture.SelectRoundByTimestamp(league, season, ts)
	}
	if err != nil {
		app.serverError(w, err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (app *application) getPlayers(w http.ResponseWriter, r *http.Request) {

	teamId, err := strconv.Atoi(r.URL.Query().Get("teamId"))
	var res interface{}

	if err != nil {
		app.serverError(w, err)
	}
	//playerId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err == nil {
		res, err = app.repo.Players.SelectPlayersByTeamId(teamId)
	}
	if err != nil {
		app.serverError(w, err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(res)

}

func (app *application) getStatistics(w http.ResponseWriter, r *http.Request) {
	teamId, err := strconv.Atoi(r.URL.Query().Get("teamId"))
	var res interface{}
	if err == nil {
		res, err = app.repo.Players.SelectPlayersAndStatisticsByTeamId(teamId)
	}

	if err != nil {
		app.clientError(w, http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(res)
}

type standingResponse struct {
	Standings *[]*database.ResultRow    `json:"standings"`
	Teams     map[int]*database.TeamRow `json:"teams"`
}

func (app *application) getLeagueStanding(w http.ResponseWriter, r *http.Request) {

	league, err := strconv.Atoi(r.URL.Query().Get("league"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	season, err := strconv.Atoi(r.URL.Query().Get("season"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	standings, err := app.repo.Result.SelectByLeagueAndSeason(league, season)
	if err != nil {
		app.serverError(w, err)
		return
	}

	teamSeason, err := app.repo.Teams.SelectTeamsSeason(league, season)
	if err != nil {
		app.serverError(w, err)
		return
	}
	var ids = make([]int, 20)
	for _, t := range *teamSeason {
		ids = append(ids, t.Team)
	}

	teams, err := app.repo.Teams.SelectTeamsByIds(&ids)
	if err != nil {
		app.serverError(w, err)
		return
	}
	ts := make(map[int]*database.TeamRow)
	for _, t := range *teams {
		ts[t.Id] = t
	}
	resp := &standingResponse{
		Teams:     ts,
		Standings: standings,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}

type fixture struct {
	Fixture *database.FixtureRow `json:"fixture"`
	Home    *database.TeamRow    `json:"home"`
	Away    *database.TeamRow    `json:"away"`
}

func (app *application) getFixture(w http.ResponseWriter, r *http.Request) {

	league, err := strconv.Atoi(r.URL.Query().Get("league"))
	if err != nil {
		app.serverError(w, err)
	}

	season, err := strconv.Atoi(r.URL.Query().Get("season"))
	if err != nil {
		app.serverError(w, err)
	}

	round, err := strconv.Atoi(r.URL.Query().Get("round"))
	if err != nil {
		app.serverError(w, err)
	}

	fixtures, err := app.repo.Fixture.SelectFixtureByLeagueSeasonRound(league, season, round)
	if err != nil {
		app.serverError(w, err)
		return
	}

	body := []fixture{}
	for _, f := range fixtures {
		ht, err := app.repo.Teams.Select(f.HomeTeam)
		if err != nil {
			app.serverError(w, err)
			return
		}

		at, err := app.repo.Teams.Select(f.AwayTeam)
		if err != nil {
			app.serverError(w, err)
			return
		}
		fixture := fixture{Fixture: f, Home: ht, Away: at}
		body = append(body, fixture)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(body)

}

// returns object to show key players of selected match
func (app *application) getPlayerStats(w http.ResponseWriter, r *http.Request) {

	league, err := strconv.Atoi(r.URL.Query().Get("league"))
	if err != nil {
		app.serverError(w, err)
	}

	season, err := strconv.Atoi(r.URL.Query().Get("season"))
	if err != nil {
		app.serverError(w, err)
	}

	team, err := strconv.Atoi(r.URL.Query().Get("homeTeam"))
	if err != nil {
		app.serverError(w, err)
	}

	awayTeam, err := strconv.Atoi(r.URL.Query().Get("awayTeam"))
	if err != nil {
		app.serverError(w, err)
	}

	round, err := strconv.Atoi(r.URL.Query().Get("round"))
	if err != nil {
		app.serverError(w, err)
	}

	players, err := app.repo.Players.SelectPlayerIdsBySeasonAndTeamId(season, team)
	if err != nil {
		app.serverError(w, err)
		return
	}

	aplayers, err := app.repo.Players.SelectPlayerIdsBySeasonAndTeamId(season, awayTeam)
	if err != nil {
		app.serverError(w, err)
		return
	}

	fixtures, err := app.repo.Fixture.SelectFixtureIdsForLastNRounds(league, season, round, 7)
	if err != nil {
		app.serverError(w, err)
		return
	}

	players = append(players, aplayers...)

	stats, err := app.repo.Players.SelectPlayerStatisticsByPlayersFixturesTeam(players, fixtures)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)

}
