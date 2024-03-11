package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (app *application) getPlayersByTeamId(w http.ResponseWriter, r *http.Request) {

	teamId, err := strconv.Atoi(r.URL.Query().Get("teamId"))

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	res, err := app.repo.Players.SelectPlayersAndStatisticsByTeamId(teamId)

	if err != nil {
		app.clientError(w, http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(res)

}

/* func getLeagueStanding(c *gin.Context) {
	l, err := strconv.Atoi(c.Param("league"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	s, err := strconv.Atoi(c.Param("season"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	ss, err := service.GetStandings(l, s)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	c.IndentedJSON(http.StatusOK, ss)
}

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
