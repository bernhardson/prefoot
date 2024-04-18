package fixture

import (
	"encoding/json"
	"time"

	"github.com/bernhardson/prefoot/pkg/coach"
	"github.com/bernhardson/prefoot/pkg/comm"
	"github.com/bernhardson/prefoot/pkg/leagues"
	"github.com/bernhardson/prefoot/pkg/players"
	"github.com/bernhardson/prefoot/pkg/shared"
	"github.com/bernhardson/prefoot/pkg/team"
)

const (
	fixturesURL      = "https://api-football-v1.p.rapidapi.com/v3/fixtures?league=%d&season=%d"
	fixtureDetailURL = "https://api-football-v1.p.rapidapi.com/v3/fixtures?id=%d"
)

// struct representing json returned by fixtures with params league and season
type FixtureResponse struct {
	Get        string `json:"get"`
	Parameters struct {
		League string `json:"league"`
		Date   string `json:"date"`
		Season string `json:"season"`
	} `json:"parameters"`
	Errors   []interface{} `json:"errors"`
	Results  int           `json:"results"`
	Paging   shared.Paging `json:"paging"`
	Response []Fixture     `json:"response"`
}
type FixtureMeta struct {
	ID        int       `json:"id"`
	Referee   string    `json:"referee"`
	Timezone  string    `json:"timezone"`
	Date      time.Time `json:"date"`
	Timestamp int       `json:"timestamp"`
	Periods   struct {
		First  int `json:"first"`
		Second int `json:"second"`
	} `json:"periods"`
	Venue  team.Venue `json:"venue"`
	Status Status     `json:"status"`
}

type TeamFixtures struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Logo   string `json:"logo"`
	Winner bool   `json:"winner"`
}

type Fixture struct {
	Fixture FixtureMeta    `json:"fixture"`
	League  leagues.League `json:"league"`
	Teams   struct {
		Home TeamFixtures `json:"home"`
		Away TeamFixtures `json:"away"`
	} `json:"teams"`
	Goals GoalsFD `json:"goals"`
	Score Score   `json:"score"`
}

// struct returned by rapid api querying with fixtureDetailURL
type FixtureDetailResponse struct {
	Get           string            `json:"get"`
	Parameters    map[string]string `json:"parameters"`
	Errors        []interface{}     `json:"errors"`
	Results       int               `json:"results"`
	Paging        shared.Paging     `json:"paging"`
	FixtureDetail []FixtureDetail   `json:"response"`
}

// Fixture struct represents fixture details
type FixtureFD struct {
	ID        int        `json:"id"`
	Referee   string     `json:"referee"`
	Timezone  string     `json:"timezone"`
	Date      string     `json:"date"`
	Timestamp int        `json:"timestamp"`
	Periods   Periods    `json:"periods"`
	Venue     team.Venue `json:"venue"`
	Status    Status     `json:"status"`
}
type Periods struct {
	First  int `json:"first"`
	Second int `json:"second"`
}

// Status struct represents fixture status details
type Status struct {
	Long    string `json:"long"`
	Short   string `json:"short"`
	Elapsed int    `json:"elapsed"`
}

// Teams struct represents teams details
type Teams struct {
	Home TeamFD `json:"home"`
	Away TeamFD `json:"away"`
}

// Team struct represents team details
type TeamFD struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Logo   string `json:"logo"`
	Winner bool   `json:"winner"`
}

// Goals struct represents goals details
type GoalsFD struct {
	Home int `json:"home"`
	Away int `json:"away"`
}

// Score struct represents score details
type Score struct {
	Halftime  GoalsFD `json:"halftime"`
	Fulltime  GoalsFD `json:"fulltime"`
	Extratime GoalsFD `json:"extratime"`
	Penalty   GoalsFD `json:"penalty"`
}

// Event struct represents event details
type Event struct {
	Time     Time      `json:"time"`
	Team     team.Team `json:"team"`
	Player   PlayerFD  `json:"player"`
	Assist   Assist    `json:"assist"`
	Type     string    `json:"type"`
	Detail   string    `json:"detail"`
	Comments string    `json:"comments"`
}

// Time struct represents time details in an event
type Time struct {
	Elapsed int `json:"elapsed"`
	Extra   int `json:"extra"`
}

// Assist struct represents assist details in an event
type Assist struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Lineup struct represents lineup details
type Lineup struct {
	Team        team.Team    `json:"team"`
	Coach       coach.Coach  `json:"coach"`
	Formation   string       `json:"formation"`
	StartXI     []StartXI    `json:"startXI"`
	Substitutes []Substitute `json:"substitutes"`
}

type PlayerLineup struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Number int    `json:"number"`
	Pos    string `json:"pos"`
	Grid   string `json:"grid"`
}

// StartXI struct represents starting XI details
type StartXI struct {
	Player PlayerLineup `json:"player"`
}

// Substitute struct represents substitute details
type Substitute struct {
	Player PlayerLineup `json:"player"`
}

// Statistic struct represents statistic details
type Statistic struct {
	Team       team.Team            `json:"team"`
	Statistics []StatisticDetailsFD `json:"statistics"`
}

// StatisticDetailsFD struct represents statistic details
type StatisticDetailsFD struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// PlayerFD struct represents player details
type PlayerFD struct {
	Team    team.Team            `json:"team"`
	Players []PlayerStatisticsFD `json:"players"`
}

// PlayerStatisticsFD struct represents player details details
type PlayerStatisticsFD struct {
	Player     players.PlayerDetails       `json:"player"`
	Statistics []PlayerStatisticsDetailsFD `json:"statistics"`
}

// FixtureDetail struct represents the complete fixture detail
type FixtureDetail struct {
	Fixture    FixtureFD      `json:"fixture"`
	League     leagues.League `json:"league"`
	Teams      Teams          `json:"teams"`
	Goals      GoalsFD        `json:"goals"`
	Score      Score          `json:"score"`
	Events     []Event        `json:"events"`
	Lineups    []Lineup       `json:"lineups"`
	Statistics []Statistic    `json:"statistics"`
	Players    []PlayerFD     `json:"players"`
}

// PlayerStatisticsDetailsFD struct represents player statistics details
type PlayerStatisticsDetailsFD struct {
	Games    GamesDetailsFD    `json:"games"`
	Offsides int               `json:"offsides"`
	Shots    ShotsDetailsFD    `json:"shots"`
	Goals    GoalsDetailsFD    `json:"goals"`
	Passes   PassesDetails     `json:"passes"`
	Tackles  TacklesDetailsFD  `json:"tackles"`
	Duels    DuelsDetailsFD    `json:"duels"`
	Dribbles DribblesDetailsFD `json:"dribbles"`
	Fouls    players.Fouls     `json:"fouls"`
	Cards    CardsDetailsFD    `json:"cards"`
	Penalty  PenaltyDetailsFD  `json:"penalty"`
}

// GamesDetailsFD struct represents games details in player statistics
type GamesDetailsFD struct {
	Minutes    int    `json:"minutes"`
	Number     int    `json:"number"`
	Position   string `json:"position"`
	Rating     string `json:"rating"`
	Captain    bool   `json:"captain"`
	Substitute bool   `json:"substitute"`
}

// ShotsDetailsFD struct represents shots details in player statistics
type ShotsDetailsFD struct {
	Total int `json:"total"`
	On    int `json:"on"`
}

// GoalsDetailsFD struct represents goals details in player statistics
type GoalsDetailsFD struct {
	Total    int `json:"total"`
	Conceded int `json:"conceded"`
	Assists  int `json:"assists"`
	Saves    int `json:"saves"`
}

// PassesDetails struct represents passes details in player statistics
type PassesDetails struct {
	Total    int    `json:"total"`
	Key      int    `json:"key"`
	Accuracy string `json:"accuracy"`
}

// TacklesDetailsFD struct represents tackles details in player statistics
type TacklesDetailsFD struct {
	Total         int `json:"total"`
	Blocks        int `json:"blocks"`
	Interceptions int `json:"interceptions"`
}

// DuelsDetailsFD struct represents duels details in player statistics
type DuelsDetailsFD struct {
	Total int `json:"total"`
	Won   int `json:"won"`
}

// DribblesDetailsFD struct represents dribbles details in player statistics
type DribblesDetailsFD struct {
	Attempts int `json:"attempts"`
	Success  int `json:"success"`
	Past     int `json:"past"`
}

// CardsDetailsFD struct represents cards details in player statistics
type CardsDetailsFD struct {
	Yellow int `json:"yellow"`
	Red    int `json:"red"`
}

// PenaltyDetailsFD struct represents penalty details in player statistics
type PenaltyDetailsFD struct {
	Won      int `json:"won"`
	Commited int `json:"commited"`
	Scored   int `json:"scored"`
	Missed   int `json:"missed"`
	Saved    int `json:"saved"`
}

// TeamStatisticsFD struct represents the "team_statistics" table
type TeamStatisticsFD struct {
	Team           int     `json:"team"`
	Fixture        int     `json:"fixture"`
	ShotsTotal     int     `json:"shots_total"`
	ShotsOn        int     `json:"shots_on"`
	ShotsOff       int     `json:"shots_off"`
	ShotsBlocked   int     `json:"shots_blocked"`
	ShotsBox       int     `json:"shots_box"`
	ShotsOutside   int     `json:"shots_outside"`
	Offsides       int     `json:"offsides"`
	Fouls          int     `json:"fouls"`
	Corners        int     `json:"corners"`
	Possession     int     `json:"possession"`
	Yellow         int     `json:"yellow"`
	Red            int     `json:"red"`
	GkSaves        int     `json:"gk_saves"`
	PassesTotal    int     `json:"passes_total"`
	PassesAccurate int     `json:"passes_accurate"`
	PassesPercent  int     `json:"passes_percent"`
	ExpectedGoals  float64 `json:"expected_goals"`
}

func FetchFixtures(league int, season int) (*FixtureResponse, error) {

	data, err := comm.GetHttpBody(fixturesURL, league, season)
	if err != nil {
		return nil, err
	}

	matches := &FixtureResponse{}
	err = json.Unmarshal(data, matches)
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func GetFixtureDetail(id int) (*FixtureDetailResponse, error) {

	data, err := comm.GetHttpBody(fixtureDetailURL, id)
	if err != nil {
		return nil, err
	}

	fd := &FixtureDetailResponse{}
	err = json.Unmarshal(data, fd)
	if err != nil {
		return nil, err
	}

	return fd, nil
}
