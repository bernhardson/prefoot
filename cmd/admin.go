package main

import (
	"net/http"
)

func (app *application) calculateStandings(w http.ResponseWriter, r *http.Request) {

	/* league, err := strconv.Atoi(r.URL.Query().Get("league"))
	if err != nil {
		app.serverError(w, err)
	}

	season, err := strconv.Atoi(r.URL.Query().Get("season"))
	if err != nil {
		app.serverError(w, err)
	} */

	/* round, err := strconv.Atoi(r.URL.Query().Get("round"))
	if err != nil {
		app.serverError(w, err)
	}

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
	json.NewEncoder(w).Encode(body) */

}
