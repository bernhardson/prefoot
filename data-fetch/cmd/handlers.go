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

func (app *application) getPlayerStats(w http.ResponseWriter, r *http.Request) {

	league, err := strconv.Atoi(r.URL.Query().Get("league"))
	if err != nil {
		app.serverError(w, err)
	}

	season, err := strconv.Atoi(r.URL.Query().Get("season"))
	if err != nil {
		app.serverError(w, err)
	}

	team, err := strconv.Atoi(r.URL.Query().Get("team"))
	if err != nil {
		app.serverError(w, err)
	}

	round, err := strconv.Atoi(r.URL.Query().Get("round"))
	if err != nil {
		app.serverError(w, err)
	}

	players, err := app.repo.Players.SelectPlayersByTeamLeagueSeason(season, team)
	if err != nil {
		app.serverError(w, err)
		return
	}

	fixtures, err := app.repo.Fixture.SelectFixtureIdsForLastNRounds(league, season, round, 7)
	if err != nil {
		app.serverError(w, err)
		return
	}
	var pIds, fIds []int
	for _, p := range players {
		pIds = append(pIds, p.Id)
	}
	for _, f := range *fixtures {
		fIds = append(fIds, *f)
	}

	playerStats, err := app.repo.Players.SelectPlayerStatisticsByPlayersFixturesTeam(pIds, fIds, team)
	if err != nil {
		app.serverError(w, err)
		return
	}

	type player struct {
		Player *database.PlayerRow      `json:"player"`
		Stats  *database.PlayerStatsRow `json:"stats"`
	}

	type response struct {
		Data map[int]*player `json:"data"`
	}

	data := make(map[int]*player)
	resp := &response{Data: data}
	for _, pl := range players {
		resp.Data[pl.Id] = &player{
			Player: pl,
		}
	}

	games := make(map[int]int)
	for _, stat := range *playerStats {

		statsAgg := resp.Data[stat.Player].Stats
		if stat.Minutes == 0 {
			continue
		}
		if statsAgg == nil {
			statsAgg = &database.PlayerStatsRow{}
		}
		// Aggregate field values into statsAgg
		statsAgg.Minutes += stat.Minutes
		statsAgg.ShotsTotal += stat.ShotsTotal
		statsAgg.ShotsOn += stat.ShotsOn
		statsAgg.GoalsScored += stat.GoalsScored
		statsAgg.GoalsAssisted += stat.GoalsAssisted
		statsAgg.PassesTotal += stat.PassesTotal
		statsAgg.PassesKey += stat.PassesKey
		statsAgg.Accuracy += stat.Accuracy
		statsAgg.Tackles += stat.Tackles
		statsAgg.Block += stat.Block
		statsAgg.Interceptions += stat.Interceptions
		statsAgg.DuelsTotal += stat.DuelsTotal
		statsAgg.DuelsWon += stat.DuelsWon
		statsAgg.DribblesTotal += stat.DribblesTotal
		statsAgg.DribblesWon += stat.DribblesWon
		statsAgg.Yellow += stat.Yellow
		statsAgg.Red += stat.Red
		statsAgg.PenaltyWon += stat.PenaltyWon
		statsAgg.PenaltyCommitted += stat.PenaltyCommitted
		statsAgg.PenaltyScored += stat.PenaltyScored
		statsAgg.PenaltyMissed += stat.PenaltyMissed
		statsAgg.PenaltySaved += stat.PenaltySaved
		statsAgg.Saves += stat.Saves
		statsAgg.Rating = statsAgg.Rating + stat.Rating

		resp.Data[stat.Player].Stats = statsAgg
		games[stat.Player]++
	}

	for playerID, entry := range resp.Data {

		if games[playerID] == 0 {
			delete(resp.Data, playerID)
			continue
		}

		statsAgg := entry.Stats
		numGames := games[playerID]
		if numGames > 0 {
			statsAgg.Accuracy /= numGames
			statsAgg.Rating /= float64(numGames)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}
