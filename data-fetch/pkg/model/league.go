package model

type LeagueData struct {
	League  League   `json:"league"`
	Country Country  `json:"country"`
	Seasons []Season `json:"seasons"`
}

type League struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Logo    string `json:"logo"`
	Flag    string `json:"flag"`
	Season  int    `json:"season"`
	Type    string `json:"type"`
	Round   string `json:"round"`
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
