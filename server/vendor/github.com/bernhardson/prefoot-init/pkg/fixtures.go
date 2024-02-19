package pkg

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

const (
	fixturesURL = "https://api-football-v1.p.rapidapi.com/v3/fixtures?league=%d&season=%d"

	fixtureDetailURL = "https://api-football-v1.p.rapidapi.com/v3/fixtures?id=%d"

	insertFixture = `INSERT INTO fixtures (id, league, referee, timezone, timestamp, venue, home_team, away_team,
						home_goals, away_goals, home_goals_half, away_goals_half)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`

	insertFormation = `INSERT INTO formations (fixture, team, formation, player1, player2, player3, player4, player5, player6, player7, player8, player9, player10, player11, sub1, sub2, sub3, sub4, sub5, coach)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`

	insertPlayerStatistics = "INSERT INTO player_statistics " +
		"(player, fixture, season, minutes, position, rating, captain, substitute, " +
		"shots_total, shots_on, goals_scored, goals_assisted, passes_total, passes_key, " +
		"accuracy, tackles, block, interceptions, duels_total, duels_won, dribbles_total, " +
		"dribbles_won, yellow, red, penalty_won, penalty_committed, penalty_scored, " +
		"penalty_missed, penalty_saved, saves)" +
		"VALUES " +
		"($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, " +
		"$18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30)"

	insertTeamStatistics = "INSERT INTO team_statistics " +
		"(team, fixture, shots_total, shots_on, shots_off, shots_blocked, " +
		"shots_box, shots_outside, offsides, fouls, corners, possession, yellow, red, " +
		"gk_saves, passes_total, passes_accurate, passes_percent, expected_goals) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)"
)

// Fixture struct represents fixture details
type FixtureFD struct {
	ID        int    `json:"id"`
	Referee   string `json:"referee"`
	Timezone  string `json:"timezone"`
	Date      string `json:"date"`
	Timestamp int    `json:"timestamp"`
	Periods   struct {
		First  int `json:"first"`
		Second int `json:"second"`
	} `json:"periods"`
	Venue  Venue  `json:"venue"`
	Status Status `json:"status"`
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
type Goals struct {
	Home int `json:"home"`
	Away int `json:"away"`
}

// Score struct represents score details
type Score struct {
	Halftime  Halftime  `json:"halftime"`
	Fulltime  Fulltime  `json:"fulltime"`
	Extratime Extratime `json:"extratime"`
	Penalty   Penalty   `json:"penalty"`
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
type Penalty struct {
	Home int `json:"home"`
	Away int `json:"away"`
}

// Event struct represents event details
type Event struct {
	Time     Time   `json:"time"`
	Team     Team   `json:"team"`
	Player   Player `json:"player"`
	Assist   Assist `json:"assist"`
	Type     string `json:"type"`
	Detail   string `json:"detail"`
	Comments string `json:"comments"`
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
	Coach       CoachFD      `json:"coach"`
	Formation   string       `json:"formation"`
	StartXI     []StartXI    `json:"startXI"`
	Substitutes []Substitute `json:"substitutes"`
}

// Coach struct represents coach details
type CoachFD struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
}

// StartXI struct represents starting XI details
type StartXI struct {
	Player Player `json:"player"`
}

// Substitute struct represents substitute details
type Substitute struct {
	Player Player `json:"player"`
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

// PlayerDetails struct represents player details
type PlayerDetails struct {
	Team    Team                   `json:"team"`
	Players []PlayerDetailsDetails `json:"players"`
}

// PlayerDetailsDetails struct represents player details details
type PlayerDetailsDetails struct {
	Player     Player                    `json:"player"`
	Statistics []PlayerStatisticsDetails `json:"statistics"`
}

// FixtureDetail struct represents the complete fixture detail
type FixtureDetail struct {
	Fixture    FixtureFD       `json:"fixture"`
	League     League          `json:"league"`
	Teams      Teams           `json:"teams"`
	Goals      Goals           `json:"goals"`
	Score      Score           `json:"score"`
	Events     []Event         `json:"events"`
	Lineups    []Lineup        `json:"lineups"`
	Statistics []Statistic     `json:"statistics"`
	Players    []PlayerDetails `json:"players"`
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
	Fouls    FoulsDetails    `json:"fouls"`
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

// FoulsDetails struct represents fouls details in player statistics
type FoulsDetails struct {
	Drawn     int `json:"drawn"`
	Committed int `json:"committed"`
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
type Fixture struct {
	ID        int    `json:"id"`
	Referee   string `json:"referee"`
	Timezone  string `json:"timezone"`
	Date      string `json:"date"`
	Timestamp int64  `json:"timestamp"`
	Periods   struct {
		First  int64 `json:"first"`
		Second int64 `json:"second"`
	} `json:"periods"`
	Venue struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		City string `json:"city"`
	} `json:"venue"`
	Status struct {
		Long    string `json:"long"`
		Short   string `json:"short"`
		Elapsed int    `json:"elapsed"`
	} `json:"status"`
}

type Match struct {
	Fixture Fixture `json:"fixture"`
	League  League  `json:"league"`
	Teams   Teams   `json:"teams"`
	Goals   Goals   `json:"goals"`
	Score   Score   `json:"score"`
}

func InsertFixtures(fixtures *[]Match, conn *pgx.Conn, season int) {
	for _, f := range *fixtures {
		//insert league
		fixtureDetail, err := GetFixtureDetail(f.Fixture.ID)
		if err != nil {
			log.Err(err).Msg("")
		}
		for _, fd := range *fixtureDetail {
			//insert fixture
			_, err := conn.Exec(insertFixture, f.Fixture.ID, f.League.ID, f.Fixture.Referee, f.Fixture.Timezone, f.Fixture.Timestamp, f.Fixture.Venue.ID, f.Teams.Home.ID, f.Teams.Away.ID,
				f.Goals.Home, f.Goals.Away, fd.Score.Halftime.Home, fd.Score.Halftime.Away)
			if err != nil {
				log.Err(err).Msg(fmt.Sprintf("fixture_%d", f.Fixture.ID))
			}
			for i, l := range fd.Lineups {
				teamStatistic := getTeamStatistics(i, &fd)
				_, err = conn.Exec(
					insertTeamStatistics,
					l.Team.ID, f.Fixture.ID, teamStatistic.ShotsTotal, teamStatistic.ShotsOn, teamStatistic.ShotsOff, teamStatistic.ShotsBlocked,
					teamStatistic.ShotsBox, teamStatistic.ShotsOutside, teamStatistic.Offsides, teamStatistic.Fouls, teamStatistic.Corners, teamStatistic.Possession, teamStatistic.Yellow, teamStatistic.Red,
					teamStatistic.GkSaves, teamStatistic.PassesTotal, teamStatistic.PassesAccurate, teamStatistic.PassesPercent, teamStatistic.ExpectedGoals,
				)
				if err != nil {
					log.Err(err).Msg(fmt.Sprintf("teamstatistics#team_%d#fixture_%d", l.Team.ID, f.Fixture.ID))
				}

				// insert formation
				_, err = conn.Exec(insertFormation,
					f.Fixture.ID, l.Team.ID, l.Formation,
					l.StartXI[0].Player.ID, l.StartXI[1].Player.ID, l.StartXI[2].Player.ID,
					l.StartXI[3].Player.ID, l.StartXI[4].Player.ID, l.StartXI[5].Player.ID,
					l.StartXI[6].Player.ID, l.StartXI[7].Player.ID, l.StartXI[8].Player.ID,
					l.StartXI[9].Player.ID, l.StartXI[10].Player.ID,
					l.Substitutes[0].Player.ID, l.Substitutes[1].Player.ID, l.Substitutes[2].Player.ID,
					l.Substitutes[3].Player.ID, l.Substitutes[4].Player.ID, l.Coach.ID,
				)
				if err != nil {
					log.Err(err).Msg("")
				}
			}
			for _, teamStats := range fd.Players {
				// insert player statistics
				for _, player := range teamStats.Players {
					ps := player.Statistics[0]
					if ps.Games.Rating == "" {
						ps.Games.Rating = "0"
					}
					if ps.Passes.Accuracy == "" {
						ps.Passes.Accuracy = "0"
					}
					_, err = conn.Exec(
						insertPlayerStatistics,
						player.Player.ID, f.Fixture.ID, season, ps.Games.Minutes, ps.Games.Position, ps.Games.Rating,
						ps.Games.Captain, ps.Games.Substitute, ps.Shots.Total, ps.Shots.On, ps.Goals.Total,
						ps.Goals.Assists, ps.Passes.Total, ps.Passes.Key, ps.Passes.Accuracy, ps.Tackles.Total,
						ps.Tackles.Blocks, ps.Tackles.Interceptions, ps.Duels.Total, ps.Duels.Won,
						ps.Dribbles.Attempts, ps.Dribbles.Success, ps.Cards.Yellow, ps.Cards.Red, ps.Penalty.Won,
						ps.Penalty.Commited, ps.Penalty.Scored, ps.Penalty.Missed, ps.Penalty.Saved, ps.Goals.Saves,
					)
					if err != nil {
						log.Err(err).Msg(fmt.Sprintf("playerstatistics#player_%d#fixture_%d", player.Player.ID, f.Fixture.ID))
					}
				}

			}
		}
	}
}

func GetMatches(league int, season int) (*[]Match, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(fixturesURL, league, season), nil)
	if err != nil {
		return nil, err
	}
	req = AddRequestHeader(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	data, err := UnmarshalData(body)
	if err != nil {
		return nil, err
	}

	matches := &[]Match{}
	err = mapstructure.Decode(data["response"], matches)
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func GetFixtureDetail(id int) (*[]FixtureDetail, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(fixtureDetailURL, id), nil)
	if err != nil {
		return nil, err
	}
	req = AddRequestHeader(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	data, err := UnmarshalData(body)
	if err != nil {
		return nil, err
	}

	fd := &[]FixtureDetail{}
	err = mapstructure.Decode(data["response"], fd)
	if err != nil {
		return nil, err
	}
	return fd, nil

}

func getTeamStatistics(home int, fd *FixtureDetail) *TeamStatistics {
	var sd TeamStatistics
	for _, s := range fd.Statistics[home].Statistics {
		val := 0
		f, ok := s.Value.(float64)
		if ok && s.Value != nil {
			val = int(f)
		}
		switch t := s.Type; t {
		case "Shots on Goal":
			sd.ShotsOn = val
		case "Shots off Goal":
			sd.ShotsOff = val
		case "Total Shots":
			sd.ShotsTotal = val
		case "Blocked Shots":
			sd.ShotsBlocked = val
		case "Shots insidebox":
			sd.ShotsBox = val
		case "Shots outsidebox":
			sd.ShotsOutside = val
		case "Fouls":
			sd.Fouls = val
		case "Corner Kicks":
			sd.Corners = val
		case "Offsides":
			sd.Offsides = val
		case "Ball Possession":
			p := strings.Replace(s.Value.(string), "%", "", -1)
			sd.Possession, _ = strconv.Atoi(p)
		case "Yellow Cards":
			sd.Yellow = val
		case "Red Cards":
			sd.Red = val
		case "Goalkeeper Saves":
			sd.GkSaves = val
		case "Total passes":
			sd.PassesTotal = val
		case "Passes accurate":
			sd.PassesAccurate = val
		case "Passes %":
			p := strings.Replace(s.Value.(string), "%", "", -1)
			sd.PassesPercent, _ = strconv.Atoi(p)
		case "expected_goals":
			fl, err := strconv.ParseFloat(s.Value.(string), 32)
			if err != nil {
				log.Err(err).Msg("")
				sd.ExpectedGoals = 0
			} else {
				sd.ExpectedGoals = fl
			}
		}
	}
	return &sd
}
