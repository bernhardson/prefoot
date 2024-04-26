package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/players/", app.getPlayers)
	router.HandlerFunc(http.MethodGet, "/statistics/", app.getStatistics)
	// ui standings table
	router.HandlerFunc(http.MethodGet, "/standings/", app.getLeagueStanding)
	// fetchCurrentRound
	router.HandlerFunc(http.MethodGet, "/rounds/", app.getRounds)
	// key player stats
	router.HandlerFunc(http.MethodGet, "/statistics/players/", app.getPlayerStats)
	// last matchups
	router.HandlerFunc(http.MethodGet, "/fixtures/matchups/", app.getLastNMatchups)
	// last matches per team
	router.HandlerFunc(http.MethodGet, "/fixtures/last/", app.getLastNFixturesByTeam)
	// round list
	router.HandlerFunc(http.MethodGet, "/fixtures/", app.getFixture)
	router.HandlerFunc(http.MethodGet, "/init/", app.initDB)
	router.HandlerFunc(http.MethodGet, "/updateDb/", app.updateDb)

	router.HandlerFunc(http.MethodPost, "/user/signup", app.userSignupPost)
	router.HandlerFunc(http.MethodPost, "/user/login", app.userLoginPost)
	router.HandlerFunc(http.MethodPost, "/user/logout", app.userLogoutPost)

	standard := alice.New(app.sessionManager.LoadAndSave, app.recoverPanic, app.logRequest, secureHeaders)
	//protected := standard.Append(app.requireAuthentication)
	return standard.Then(router)
}
