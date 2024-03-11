package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	//mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/players/", app.getPlayersByTeamId)

	return mux
}
