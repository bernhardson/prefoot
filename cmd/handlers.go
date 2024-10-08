package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bernhardson/prefoot/internal/models"
	"github.com/bernhardson/prefoot/internal/validator"
	"github.com/bernhardson/prefoot/pkg/fixture"
	"github.com/bernhardson/prefoot/pkg/result"
	"github.com/bernhardson/prefoot/pkg/team"
	"github.com/julienschmidt/httprouter"
)

// get last round by timestamp
// used in the round list screen to get the initially displayed round
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
		return
	}

	res, err = app.fixture.RoundRepo.SelectRoundByTimestamp(league, season, ts)

	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// unused right now
func (app *application) getPlayers(w http.ResponseWriter, r *http.Request) {

	teamId, err := strconv.Atoi(r.URL.Query().Get("teamId"))

	if err != nil {
		app.serverError(w, err)
		return
	}
	res, err := app.player.Repo.SelectPlayersByTeamId(teamId)

	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

}

// unused right now
func (app *application) getStatistics(w http.ResponseWriter, r *http.Request) {
	teamId, err := strconv.Atoi(r.URL.Query().Get("teamId"))
	if err != nil {
		app.clientError(w, http.StatusNotFound)
		return
	}

	res, err := app.player.Repo.SelectPlayersAndStatisticsByTeamId(teamId)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

type standingResponse struct {
	Standings *[]*result.ResultRow  `json:"standings"`
	Teams     map[int]*team.TeamRow `json:"teams"`
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

	results, err := app.fixture.ResultRepo.SelectByLeagueSeason(league, season)
	if err != nil {
		app.serverError(w, err)
		return
	}

	teamSeason, err := app.team.TeamRepo.SelectTeamsSeason(league, season)
	if err != nil {
		app.serverError(w, err)
		return
	}
	var ids = make([]int, 20)
	for _, t := range *teamSeason {
		ids = append(ids, t.Team)
	}

	teams, err := app.team.TeamRepo.SelectTeamsByIds(&ids)
	if err != nil {
		app.serverError(w, err)
		return
	}
	ts := make(map[int]*team.TeamRow)
	for _, t := range *teams {
		ts[t.Id] = t
	}
	resp := &standingResponse{
		Teams:     ts,
		Standings: results,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}

type fixtureResp struct {
	Fixture *fixture.FixtureRow `json:"fixture"`
	Home    *team.TeamRow       `json:"home"`
	Away    *team.TeamRow       `json:"away"`
}

func (app *application) getFixture(w http.ResponseWriter, r *http.Request) {

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

	round, err := strconv.Atoi(r.URL.Query().Get("round"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	fixtures, err := app.fixture.Repo.SelectFixtureByLeagueSeasonRound(league, season, round)
	if err != nil {
		app.serverError(w, err)
		return
	}

	body := []fixtureResp{}
	for _, f := range fixtures {
		ht, err := app.team.TeamRepo.Select(f.HomeTeam)
		if err != nil {
			app.serverError(w, err)
			return
		}

		at, err := app.team.TeamRepo.Select(f.AwayTeam)
		if err != nil {
			app.serverError(w, err)
			return
		}
		fixture := fixtureResp{Fixture: f, Home: ht, Away: at}
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
		return
	}

	season, err := strconv.Atoi(r.URL.Query().Get("season"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	team, err := strconv.Atoi(r.URL.Query().Get("homeTeam"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	awayTeam, err := strconv.Atoi(r.URL.Query().Get("awayTeam"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	round, err := strconv.Atoi(r.URL.Query().Get("round"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	players, err := app.player.Repo.SelectPlayerIdsBySeasonAndTeamId(season, team)
	if err != nil {
		app.serverError(w, err)
		return
	}

	aplayers, err := app.player.Repo.SelectPlayerIdsBySeasonAndTeamId(season, awayTeam)
	if err != nil {
		app.serverError(w, err)
		return
	}

	fixtures, err := app.fixture.Repo.SelectFixtureIdsForLastNRounds(league, season, round, 7)
	if err != nil {
		app.serverError(w, err)
		return
	}

	players = append(players, aplayers...)

	stats, err := app.player.Repo.SelectPlayerStatisticsByPlayersFixturesTeam(players, fixtures)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)

}

type matchups struct {
	Fixture []*fixture.FixtureRow `json:"fixture"`
	Teams   map[int]team.TeamRow  `json:"teams"`
}

func (app *application) getLastNMatchups(w http.ResponseWriter, r *http.Request) {

	team1, err := strconv.Atoi(r.URL.Query().Get("team1"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	team2, err := strconv.Atoi(r.URL.Query().Get("team2"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	n, err := strconv.Atoi(r.URL.Query().Get("n"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	team1Row, err := app.team.TeamRepo.Select(team1)
	if err != nil {
		app.serverError(w, err)
		return
	}
	team2Row, err := app.team.TeamRepo.Select(team2)
	if err != nil {
		app.serverError(w, err)
		return
	}

	fixtures, err := app.fixture.Repo.SelectLastNMatchups(team1, team2, n)
	if err != nil {
		app.serverError(w, err)
		return
	}

	resp := &matchups{
		Fixture: fixtures,
		Teams:   map[int]team.TeamRow{team1: *team1Row, team2: *team2Row},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

type lastNFixture struct {
	Fixture *fixture.FixtureRow
	Result  *result.ResultRow
}

func (app *application) getLastNFixturesByTeam(w http.ResponseWriter, r *http.Request) {

	team1, err := strconv.Atoi(r.URL.Query().Get("team"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	n, err := strconv.Atoi(r.URL.Query().Get("n"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	ts := int(time.Now().Unix())

	fixtures, err := app.fixture.Repo.SelectLastNFixturesByTeam(team1, ts, n)
	if err != nil {
		app.serverError(w, err)
		return
	}

	var resp []*lastNFixture
	for _, f := range fixtures {
		app.logger.Debug().Msg(fmt.Sprintf("query results for league %d, season %d, team %d, round %d", f.League, f.Season, team1, f.Round))
		result, err := app.fixture.ResultRepo.SelectResultByLeagueSeasonTeamRound(f.League, f.Season, team1, f.Round)
		if err != nil {
			if strings.HasPrefix(err.Error(), "no rows in result set") {
				continue
			}
			app.serverError(w, err)
			return
		}
		resp = append(resp, &lastNFixture{Fixture: f, Result: result})
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (app *application) initDB(w http.ResponseWriter, r *http.Request) {

	var params struct {
		Leagues []int `json:"leagues"`
	}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		app.logger.Err(err).Msg("failed to decode JSON payload")
		http.Error(w, "failed to decode JSON payload", http.StatusBadRequest)
		return
	}

	//leagues := []int{71, 137}

	for _, l := range params.Leagues {
		lresp, fs, err := app.league.FetchAndInsertLeague(l)
		app.logger.Err(err).Msg(fmt.Sprintf("insert leagues: failed=%v", *fs))

		for _, s := range lresp.Response[0].Seasons {
			year := s.Year

			app.logger.Info().Msg(fmt.Sprintf("inserting league=%d#season=%d", l, year))
			ts, err := app.team.FetchAndInsertTeams(l, year)
			app.logger.Err(err).Msg(fmt.Sprintf("insert teams: league=%d#season=%d", l, year))

			fp, fs, err := app.player.FetchAndInsertPlayers(l, year)
			app.logger.Err(err).Msg(fmt.Sprintf("insert players: league:%d#season=%d", l, year))
			app.logger.Err(err).Msg(fmt.Sprintf("insert players: failedP=%v # failedS=%v", *fp, *fs))

			err = app.fixture.FetchAndInsertFixtures(l, year)
			app.logger.Err(err).Msg(fmt.Sprintf("insert fixtures: league:%d#season=%d", l, year))

			fc, fcc, err := app.coach.FetchAndInsertCoaches(ts)
			app.logger.Err(err).Msg(fmt.Sprintf("insert coaches: league:%d#season=%d#failedC=%v#failedCC=%v", l, year, *fc, *fcc))
		}
	}
}

func (app *application) updateDb(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	league, err := strconv.Atoi(params.ByName("league"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	season, err := strconv.Atoi(params.ByName("season"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.fixture.UpdateFixture(league, season)
}

type userSignupForm struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	// Try to create a new user record in the database. If the email already // exists then add an error message to the form and re-display it.

	params := httprouter.ParamsFromContext(r.Context())

	name := params.ByName("name")

	email := params.ByName("email")

	password := params.ByName("password")

	err := app.users.Insert(name, email, password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			//return error
		} else {
			app.serverError(w, err)
		}
		return
	}
}

type user struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) { // Declare an zero-valued instance of our userSignupForm struct.

	body, err := io.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	var user user
	err = json.Unmarshal(body, &user)
	if err != nil {
		app.serverError(w, err)
	}

	valid := validator.NotBlank(user.Name)
	valid = valid && validator.NotBlank(user.Email)
	valid = valid && validator.Matches(user.Email, validator.EmailRX)
	valid = valid && validator.NotBlank(user.Password)
	valid = valid && validator.MinChars(user.Password, 8)

	// If there are any errors, redisplay the signup form along with a 422 // status code.

	if !valid {
		app.serverError(w, errors.New("user email and/or password are invalid"))
		return
	}

	// Try to create a new user record in the database. If the email already
	// exists then add an error message to the form and re-display it.
	err = app.users.Insert(user.Name, user.Email, user.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			//todo return proper error message
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Otherwise add a confirmation flash message to the session confirming that
	// their signup worked.
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {

	var user user

	valid := validator.NotBlank(user.Email)
	valid = valid && validator.Matches(user.Email, validator.EmailRX)
	valid = valid && validator.NotBlank(user.Password)

	if !valid {
		app.serverError(w, errors.New("user email or password wrong"))
		return
	}

	id, err := app.users.Authenticate(user.Email, user.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			app.serverError(w, models.ErrInvalidCredentials)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	// path := app.sessionManager.PopString(r.Context(), "redirectPathAfterLogin")
	// if path != "" {
	// 	http.Redirect(w, r, path, http.StatusSeeOther)
	// 	return
	// }

	// http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) { // Use the RenewToken() method on the current session to change the session // ID again.

	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	// http.Redirect(w, r, "/", http.StatusSeeOther)
}
