package pkg

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jackc/pgx"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

const (
	playersURL                   = "https://api-football-v1.p.rapidapi.com/v3/players?league=%d&season=%d&page=%d"
	playersSquadURL              = "https://api-football-v1.p.rapidapi.com/v3/players/squads?team=%d"
	insertPlayer                 = `INSERT INTO players (id, team, season, firstname, lastname) VALUES ($1, $2, $3, $4, $5)`
	insertPlayerStatisticsSeason = `INSERT INTO player_statistics_season ("player", "season", "team","minutes", "position", "rating", "captain", "games","lineups", "shots_total", "shots_on",
									"goals_scored", "goals_assisted", "passes_total", "passes_key", "accuracy", "tackles", "block", "interceptions", "duels_total", "duels_won", 
									"dribbles_total", "dribbles_won", "yellow", "red", "penalty_won", "penalty_committed", "penalty_scored", "penalty_missed", "penalty_saved", "saves")
        							VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31)`
)

type PlayerSeason struct {
	Player     Player       `json:"player"`
	Statistics []Statistics `json:"statistics"`
}

type Player struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Age       int    `json:"age"`
	Birth     struct {
		Date    string `json:"date"`
		Place   string `json:"place"`
		Country string `json:"country"`
	} `json:"birth"`
	Nationality string `json:"nationality"`
	Height      string `json:"height"`
	Weight      string `json:"weight"`
	Injured     bool   `json:"injured"`
	Photo       string `json:"photo"`
}

type TeamPl struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

type LeaguePl struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Logo    string `json:"logo"`
	Flag    string `json:"flag"`
	Season  int    `json:"season"`
}

type Statistics struct {
	Team   Team     `json:"team"`
	League LeaguePl `json:"league"`
	Games  struct {
		Appearances int    `json:"appearances"`
		Lineups     int    `json:"lineups"`
		Minutes     int    `json:"minutes"`
		Number      int    `json:"number"`
		Position    string `json:"position"`
		Rating      string `json:"rating"`
		Captain     bool   `json:"captain"`
	} `json:"games"`
	Substitutes struct {
		In    int `json:"in"`
		Out   int `json:"out"`
		Bench int `json:"bench"`
	} `json:"substitutes"`
	Shots struct {
		Total int `json:"total"`
		On    int `json:"on"`
	} `json:"shots"`
	Goals struct {
		Total    int `json:"total"`
		Conceded int `json:"conceded"`
		Assists  int `json:"assists"`
		Saves    int `json:"saves"`
	} `json:"goals"`
	Passes struct {
		Total    int     `json:"total"`
		Key      int     `json:"key"`
		Accuracy float64 `json:"accuracy"`
	} `json:"passes"`
	Tackles struct {
		Total         int `json:"total"`
		Blocks        int `json:"blocks"`
		Interceptions int `json:"interceptions"`
	} `json:"tackles"`
	Duels struct {
		Total int `json:"total"`
		Won   int `json:"won"`
	} `json:"duels"`
	Dribbles struct {
		Attempts int `json:"attempts"`
		Success  int `json:"success"`
		Past     int `json:"past"`
	} `json:"dribbles"`
	Fouls struct {
		Drawn     int `json:"drawn"`
		Committed int `json:"committed"`
	} `json:"fouls"`
	Cards struct {
		Yellow    int `json:"yellow"`
		Yellowred int `json:"yellowred"`
		Red       int `json:"red"`
	} `json:"cards"`
	Penalty struct {
		Won       int `json:"won"`
		Committed int `json:"committed"`
		Scored    int `json:"scored"`
		Missed    int `json:"missed"`
		Saved     int `json:"saved"`
	} `json:"penalty"`
}

type Paging struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

// simple list of current players in a squad
type Squad struct {
	Team    Team     `json:"team"`
	Players []Player `json:"players"`
}

func GetPlayers(league int, season int) (*[]PlayerSeason, error) {
	ret := []PlayerSeason{}
	for i := 1; i < 100; i++ {
		data, err := GetHttpBody(playersURL, season, 1)
		if err != nil {
			return nil, err
		}

		pg := &Paging{}
		mapstructure.Decode(data["paging"], &pg)

		p := []PlayerSeason{}
		err = mapstructure.Decode(data["response"], &p)
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("%d", i))
		}
		ret = append(ret, p...)
		//TODO remove true
		if pg.Current == pg.Total || true {
			break
		}
	}
	return &ret, nil
}

func GetPlayersByTeamId(teamId int) (*[]Player, error) {
	data, err := GetHttpBody(playersSquadURL, teamId)
	if err != nil {
		return nil, err
	}
	sqs := []Squad{}
	err = mapstructure.Decode(data["response"], &sqs)
	if err != nil {
		return nil, err
	}
	for _, sq := range sqs {
		return &sq.Players, nil
	}
	return nil, errors.New("WTF")
}

func GetPlayerIdsByTeamId(teamId int) (*[]int, *[]Player, error) {
	data, err := GetHttpBody(playersSquadURL, teamId)
	if err != nil {
		return nil, nil, err
	}
	ps := []Squad{}
	err = mapstructure.Decode(data["response"], &ps)
	if err != nil {
		return nil, nil, err
	}

	res := []int{}

	for _, s := range ps {
		for _, p := range s.Players {
			res = append(res, p.ID)
		}
	}
	return &res, &ps[0].Players, nil
}

func InsertPlayers(players *[]PlayerSeason, conn *pgx.Conn) {
	for _, p := range *players {
		// player statistics only has one entry so there will be just one insert to player table
		for _, s := range p.Statistics {
			_, err := conn.Exec(
				insertPlayer,
				p.Player.ID, s.Team.ID, s.League.Season, p.Player.Firstname, p.Player.Lastname)
			if err != nil {
				log.Err(err).Msg(fmt.Sprintf("player id %d", p.Player.ID))
			} else {
				log.Info().Msg(fmt.Sprintf("player id %d", p.Player.ID))
			}
			rating, err := strconv.ParseFloat(s.Games.Rating, 32)
			if err != nil {
				rating = 0
				log.Err(err).Msg("")
			}
			_, err = conn.Exec(
				insertPlayerStatisticsSeason,
				p.Player.ID, s.League.Season, s.Team.ID, s.Games.Minutes, s.Games.Position, rating,
				s.Games.Captain, s.Games.Appearances, s.Games.Lineups, s.Shots.Total, s.Shots.On, s.Goals.Total,
				s.Goals.Assists, s.Passes.Total, s.Passes.Key, s.Passes.Accuracy, s.Tackles.Total,
				s.Tackles.Blocks, s.Tackles.Interceptions, s.Duels.Total, s.Duels.Won,
				s.Dribbles.Attempts, s.Dribbles.Success, s.Cards.Yellow, s.Cards.Red,
				s.Penalty.Won, s.Penalty.Committed, s.Penalty.Scored, s.Penalty.Missed, s.Penalty.Saved, s.Goals.Saves,
			)
			if err != nil {
				log.Err(err).Msg(fmt.Sprintf("failed inserting player season statistics %d", p.Player.ID))
			} else {
				log.Debug().Msg(fmt.Sprintf("inserted player %d", p.Player.ID))
			}
		}
	}
}
