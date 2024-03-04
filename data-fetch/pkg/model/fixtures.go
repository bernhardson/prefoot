package model

type FixtureResponse struct {
	Get           string            `json:"get"`
	Parameters    map[string]string `json:"parameters"`
	Errors        []interface{}     `json:"errors"`
	Results       int               `json:"results"`
	Paging        Paging            `json:"paging"`
	FixtureDetail []FixtureDetail   `json:"response"`
}

// Fixture struct represents fixture details
type FixtureFD struct {
	ID        int     `json:"id"`
	Referee   string  `json:"referee"`
	Timezone  string  `json:"timezone"`
	Date      string  `json:"date"`
	Timestamp int     `json:"timestamp"`
	Periods   Periods `json:"periods"`
	Venue     Venue   `json:"venue"`
	Status    Status  `json:"status"`
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
	Halftime  Halftime  `json:"halftime"`
	Fulltime  Fulltime  `json:"fulltime"`
	Extratime Extratime `json:"extratime"`
	Penalty   PenaltyFD `json:"penalty"`
}

// Halftime struct represents halftime score details
type Halftime struct {
	Home int `json:"home"`
	Away int `json:"away"`
}

// Fulltime struct represents fulltime score details
type Fulltime struct {
	Home int `json:"home"`
	Away int `json:"away"`
}

// Extratime struct represents extratime score details
type Extratime struct {
	Home int `json:"home"`
	Away int `json:"away"`
}

// Penalty struct represents penalty details
type PenaltyFD struct {
	Home int `json:"home"`
	Away int `json:"away"`
}

// Event struct represents event details
type Event struct {
	Time     Time     `json:"time"`
	Team     Team     `json:"team"`
	Player   PlayerFD `json:"player"`
	Assist   Assist   `json:"assist"`
	Type     string   `json:"type"`
	Detail   string   `json:"detail"`
	Comments string   `json:"comments"`
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
	Team        Team         `json:"team"`
	Coach       Coach        `json:"coach"`
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
	Team       Team               `json:"team"`
	Statistics []StatisticDetails `json:"statistics"`
}

// StatisticDetails struct represents statistic details
type StatisticDetails struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// PlayerFD struct represents player details
type PlayerFD struct {
	Team    Team                 `json:"team"`
	Players []PlayerStatisticsFD `json:"players"`
}

// PlayerStatisticsFD struct represents player details details
type PlayerStatisticsFD struct {
	Player     PlayerDetails             `json:"player"`
	Statistics []PlayerStatisticsDetails `json:"statistics"`
}

// FixtureDetail struct represents the complete fixture detail
type FixtureDetail struct {
	Fixture    FixtureFD   `json:"fixture"`
	League     League      `json:"league"`
	Teams      Teams       `json:"teams"`
	Goals      GoalsFD     `json:"goals"`
	Score      Score       `json:"score"`
	Events     []Event     `json:"events"`
	Lineups    []Lineup    `json:"lineups"`
	Statistics []Statistic `json:"statistics"`
	Players    []PlayerFD  `json:"players"`
}

// PlayerStatisticsDetails struct represents player statistics details
type PlayerStatisticsDetails struct {
	Games    GamesDetails    `json:"games"`
	Offsides int             `json:"offsides"`
	Shots    ShotsDetails    `json:"shots"`
	Goals    GoalsDetails    `json:"goals"`
	Passes   PassesDetails   `json:"passes"`
	Tackles  TacklesDetails  `json:"tackles"`
	Duels    DuelsDetails    `json:"duels"`
	Dribbles DribblesDetails `json:"dribbles"`
	Fouls    Fouls           `json:"fouls"`
	Cards    CardsDetails    `json:"cards"`
	Penalty  PenaltyDetails  `json:"penalty"`
}

// GamesDetails struct represents games details in player statistics
type GamesDetails struct {
	Minutes    int    `json:"minutes"`
	Number     int    `json:"number"`
	Position   string `json:"position"`
	Rating     string `json:"rating"`
	Captain    bool   `json:"captain"`
	Substitute bool   `json:"substitute"`
}

// ShotsDetails struct represents shots details in player statistics
type ShotsDetails struct {
	Total int `json:"total"`
	On    int `json:"on"`
}

// GoalsDetails struct represents goals details in player statistics
type GoalsDetails struct {
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

// TacklesDetails struct represents tackles details in player statistics
type TacklesDetails struct {
	Total         int `json:"total"`
	Blocks        int `json:"blocks"`
	Interceptions int `json:"interceptions"`
}

// DuelsDetails struct represents duels details in player statistics
type DuelsDetails struct {
	Total int `json:"total"`
	Won   int `json:"won"`
}

// DribblesDetails struct represents dribbles details in player statistics
type DribblesDetails struct {
	Attempts int `json:"attempts"`
	Success  int `json:"success"`
	Past     int `json:"past"`
}

// CardsDetails struct represents cards details in player statistics
type CardsDetails struct {
	Yellow int `json:"yellow"`
	Red    int `json:"red"`
}

// PenaltyDetails struct represents penalty details in player statistics
type PenaltyDetails struct {
	Won      int `json:"won"`
	Commited int `json:"commited"`
	Scored   int `json:"scored"`
	Missed   int `json:"missed"`
	Saved    int `json:"saved"`
}

// TeamStatistics struct represents the "team_statistics" table
type TeamStatistics struct {
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



