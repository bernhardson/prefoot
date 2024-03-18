package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bernhardson/prefoot/data-fetch/pkg/database"
	"github.com/bernhardson/prefoot/data-fetch/pkg/fetch"
)

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

	w.Header().Set("Content-Type", "application/json")
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(res)
}

func (app *application) getLeagueStanding(w http.ResponseWriter, r *http.Request) {

	league, err := strconv.Atoi(r.URL.Query().Get("league"))
	if err != nil {
		app.serverError(w, err)
	}

	season, err := strconv.Atoi(r.URL.Query().Get("season"))
	if err != nil {
		app.serverError(w, err)
	}

	res, err := fetch.GetStanding(league, season)
	if err != nil {
		app.serverError(w, err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res.League.Standings[0])
	}
}

type fixture struct {
	Fixture *database.FixtureRow `json:"fixture"`
	Home    *database.TeamRow    `json:"home"`
	Away    *database.TeamRow    `json:"away"`
}

func (app *application) getFixture(w http.ResponseWriter, r *http.Request) {

	/* league, err := strconv.Atoi(r.URL.Query().Get("league"))
	if err != nil {
		app.serverError(w, err)
	}

	season, err := strconv.Atoi(r.URL.Query().Get("season"))
	if err != nil {
		app.serverError(w, err)
	} */

	round, err := strconv.Atoi(r.URL.Query().Get("round"))
	if err != nil {
		app.serverError(w, err)
	}

	fixtures, err := app.repo.Fixture.SelectFixturesByRound(round)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(body)

}

/*
func getFixtures(c *gin.Context) {
	l, err := strconv.Atoi(c.Param("league"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	s, err := strconv.Atoi(c.Param("season"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	r, err := strconv.Atoi(c.Param("round"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	ss, err := service.GetFixtures(l, s, r, Pool, Logger)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	c.IndentedJSON(http.StatusOK, ss)
}
*/
