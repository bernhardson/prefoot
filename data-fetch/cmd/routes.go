package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	//mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/players/", app.getPlayers)
	mux.HandleFunc("/statistics/", app.getStatistics)
	// ui standings table
	mux.HandleFunc("/standings/", app.getLeagueStanding)
	// fetchCurrentRound
	mux.HandleFunc("/rounds/", app.getRounds)
	// key player stats
	mux.HandleFunc("/statistics/players/", app.getPlayerStats)
	// last matchups
	mux.HandleFunc("/fixtures/matchups/", app.getLastNMatchups)
	// last matches per team
	mux.HandleFunc("/fixtures/last/", app.getLastNFixturesByTeam)
	// round list
	mux.HandleFunc("/fixtures/", app.getFixture)

	return mux
}
