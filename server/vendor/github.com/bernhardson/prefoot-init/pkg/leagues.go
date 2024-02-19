package pkg

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

const (
	getLeaguesURL  = "https://api-football-v1.p.rapidapi.com/v3/leagues"
	getStandingURL = "https://api-football-v1.p.rapidapi.com/v3/standings?season=%d&league=%d"
	insertLeague   = `INSERT INTO leagues (id, name) VALUES ($1, $2);`
	insertVenue    = `INSERT INTO venues (id, name, city) VALUES ($1, $2, $3)`
)

type LeagueData struct {
	League  League   `json:"league"`
	Country Country  `json:"country"`
	Seasons []Season `json:"seasons"`
}

type League struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Logo string `json:"logo"`
}

type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Flag string `json:"flag"`
}

type Season struct {
	Year     int            `json:"year"`
	Start    string         `json:"start"`
	End      string         `json:"end"`
	Current  bool           `json:"current"`
	Coverage SeasonCoverage `json:"coverage"`
}

type SeasonCoverage struct {
	Fixtures    SeasonCoverageFixtures `json:"fixtures"`
	Standings   bool                   `json:"standings"`
	Players     bool                   `json:"players"`
	TopScorers  bool                   `json:"top_scorers"`
	TopAssists  bool                   `json:"top_assists"`
	TopCards    bool                   `json:"top_cards"`
	Injuries    bool                   `json:"injuries"`
	Predictions bool                   `json:"predictions"`
	Odds        bool                   `json:"odds"`
}

type SeasonCoverageFixtures struct {
	Events             bool `json:"events"`
	Lineups            bool `json:"lineups"`
	StatisticsFixtures bool `json:"statistics_fixtures"`
	StatisticsPlayers  bool `json:"statistics_players"`
}

type StandingsResponse struct {
	Get        string           `json:"get"`
	Parameters StandingsParams  `json:"parameters"`
	Errors     []interface{}    `json:"errors"`
	Results    int              `json:"results"`
	Paging     StandingsPaging  `json:"paging"`
	Response   []StandingsEntry `json:"response"`
}

type StandingsParams struct {
	League string `json:"league"`
	Season string `json:"season"`
}

type StandingsPaging struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

type StandingsEntry struct {
	League StandingsLeague `json:"league"`
}

type StandingsLeague struct {
	ID        int               `json:"id"`
	Name      string            `json:"name"`
	Country   string            `json:"country"`
	Logo      string            `json:"logo"`
	Flag      string            `json:"flag"`
	Season    int               `json:"season"`
	Standings [][]StandingsTeam `json:"standings"`
}

type StandingsTeam struct {
	Rank        int                 `json:"rank"`
	Team        StandingsTeamDetail `json:"team"`
	Points      int                 `json:"points"`
	GoalsDiff   int                 `json:"goalsDiff"`
	Group       string              `json:"group"`
	Form        string              `json:"form"`
	Status      string              `json:"status"`
	Description string              `json:"description"`
	All         StandingsStats      `json:"all"`
	Home        StandingsStats      `json:"home"`
	Away        StandingsStats      `json:"away"`
	Update      string              `json:"update"`
}

type StandingsTeamDetail struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

type StandingsStats struct {
	Played int                `json:"played"`
	Win    int                `json:"win"`
	Draw   int                `json:"draw"`
	Lose   int                `json:"lose"`
	Goals  StandingsGoalStats `json:"goals"`
}

type StandingsGoalStats struct {
	For     int `json:"for"`
	Against int `json:"against"`
}

func InsertLeagues(conn *pgx.Conn) (*[]int, *[]int, error) {
	ls, err := getLeagues()
	if err != nil {
		log.Err(err).Msg("")
	}
	insertedLeague := []int{}
	failedInsert := []int{}

	for _, l := range *ls {
		_, err = conn.Exec(
			insertLeague,
			l.League.ID, l.League.Name)
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("%d", l.League.ID))
			failedInsert = append(failedInsert, l.League.ID)
		} else {
			insertedLeague = append(insertedLeague, l.League.ID)
		}
	}
	return &insertedLeague, &failedInsert, nil
}

func getLeagues() (*[]LeagueData, error) {
	data, err := GetHttpBody(getLeaguesURL, nil)
	if err != nil {
		return nil, err
	}
	l := &[]LeagueData{}
	mapstructure.Decode(data["response"], l)

	return l, nil
}

func GetStanding(league int, season int) (*[][]StandingsTeam, error) {
	data, err := GetHttpBody(getStandingURL, season, league)
	if err != nil {
		return nil, err
	}
	ses := []StandingsEntry{}
	mapstructure.Decode(data["response"], &ses)
	for _, se := range ses {
		return &se.League.Standings, nil
	}
	return nil, errors.New("wtf")
}
