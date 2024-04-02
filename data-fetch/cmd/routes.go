package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	//mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/players/", app.getPlayers)
	mux.HandleFunc("/statistics/", app.getStatistics)
	mux.HandleFunc("/standings/", app.getLeagueStanding)
	mux.HandleFunc("/fixtures/", app.getFixture)
	mux.HandleFunc("/rounds/", app.getRounds)

	return mux
}
