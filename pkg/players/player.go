package players

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bernhardson/prefoot/pkg/comm"
	"github.com/bernhardson/prefoot/pkg/leagues"
	"github.com/bernhardson/prefoot/pkg/shared"
	"github.com/bernhardson/prefoot/pkg/team"
)

const (
	playersURL      = "https://api-football-v1.p.rapidapi.com/v3/players?league=%d&season=%d&page=%d"
	playersSquadURL = "https://api-football-v1.p.rapidapi.com/v3/players/squads?team=%d"
	playersIdURL    = "https://api-football-v1.p.rapidapi.com/v3/players?id=%d&season=%d"
)

type PlayerAPIResponse struct {
	Get        string           `json:"get"`
	Parameters PlayerParameters `json:"parameters"`
	Errors     []interface{}    `json:"errors"`
	Results    float64          `json:"results"`
	Paging     shared.Paging    `json:"paging"`
	Response   []Player         `json:"response"`
}

type Player struct {
	PlayerDetails PlayerDetails      `json:"player"`
	Statistics    []PlayerStatistics `json:"statistics"`
}

type PlayerParameters struct {
	League string `json:"league"`
	Page   string `json:"page"`
	Season string `json:"season"`
}

type PlayerDetails struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	FirstName   string       `json:"firstname"`
	LastName    string       `json:"lastname"`
	Age         int          `json:"age"`
	Birth       shared.Birth `json:"birth"`
	Nationality string       `json:"nationality"`
	Height      string       `json:"height"`
	Weight      string       `json:"weight"`
	Injured     bool         `json:"injured"`
	Photo       string       `json:"photo"`
}

type BirthInfo struct {
	Date    string `json:"date"`
	Place   string `json:"place"`
	Country string `json:"country"`
}

type Games struct {
	Appearances int    `json:"appearances"`
	Lineups     int    `json:"lineups"`
	Minutes     int    `json:"minutes"`
	Number      int    `json:"number"`
	Position    string `json:"position"`
	Rating      string `json:"rating"`
	Captain     bool   `json:"captain"`
}

type Substitutes struct {
	In    int `json:"in"`
	Out   int `json:"out"`
	Bench int `json:"bench"`
}

type Shots struct {
	Total int `json:"total"`
	On    int `json:"on"`
}

type Goals struct {
	Total    int `json:"total"`
	Conceded int `json:"conceded"`
	Assists  int `json:"assists"`
	Saves    int `json:"saves"`
}

type Passes struct {
	Total    int `json:"total"`
	Key      int `json:"key"`
	Accuracy int `json:"accuracy"`
}

type Tackles struct {
	Total         int `json:"total"`
	Blocks        int `json:"blocks"`
	Interceptions int `json:"interceptions"`
}

type Duels struct {
	Total int `json:"total"`
	Won   int `json:"won"`
}

type Dribbles struct {
	Attempts int `json:"attempts"`
	Success  int `json:"success"`
	Past     int `json:"past"`
}

type Fouls struct {
	Drawn     int `json:"drawn"`
	Committed int `json:"committed"`
}

type Cards struct {
	Yellow    int `json:"yellow"`
	YellowRed int `json:"yellowred"`
	Red       int `json:"red"`
}

type Penalty struct {
	Won       int `json:"won"`
	Committed int `json:"committed"`
	Scored    int `json:"scored"`
	Missed    int `json:"missed"`
	Saved     int `json:"saved"`
}

type PlayerStatistics struct {
	Team        team.Team      `json:"team"`
	League      leagues.League `json:"league"`
	Games       Games          `json:"games"`
	Substitutes Substitutes    `json:"substitutes"`
	Shots       Shots          `json:"shots"`
	Goals       Goals          `json:"goals"`
	Passes      Passes         `json:"passes"`
	Tackles     Tackles        `json:"tackles"`
	Duels       Duels          `json:"duels"`
	Dribbles    Dribbles       `json:"dribbles"`
	Fouls       Fouls          `json:"fouls"`
	Cards       Cards          `json:"cards"`
	Penalty     Penalty        `json:"penalty"`
}

func GetPlayers(league int, season int, paging int) (*[]Player, *shared.Paging, error) {

	data, err := comm.GetHttpBody(playersURL, league, season, paging)
	if err != nil {
		return nil, nil, err
	}

	p := PlayerAPIResponse{}
	err = json.Unmarshal(data, &p)
	if err != nil {
		return nil, nil, err
	}
	return &p.Response, &p.Paging, nil
}

type SquadResponse struct {
	Get        string        `json:"get"`
	Parameters interface{}   `json:"parameters"`
	Errors     []interface{} `json:"errors"`
	Results    int           `json:"results"`
	Response   []Squad       `json:"response"`
}

type PlayerSquad struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Number   int    `json:"number"`
	Position string `json:"position"`
	Photo    string `json:"photo"`
}

type Squad struct {
	Team    team.Team     `json:"team"`
	Players []PlayerSquad `json:"players"`
}

func GetPlayersByTeamId(teamId int) (*[]PlayerSquad, error) {
	data, err := comm.GetHttpBody(playersSquadURL, teamId)
	if err != nil {
		return nil, err
	}

	sqs := SquadResponse{}
	err = json.Unmarshal(data, &sqs)
	if err != nil {
		return nil, err
	}

	for _, sq := range sqs.Response {
		return &sq.Players, nil
	}
	return nil, errors.New("WTF")
}

func GetPlayerIdsByTeamId(teamId int) (*[]int, *[]PlayerSquad, error) {
	data, err := comm.GetHttpBody(playersSquadURL, teamId)
	if err != nil {
		return nil, nil, err
	}
	sqs := SquadResponse{}
	err = json.Unmarshal(data, &sqs)
	if err != nil {
		return nil, nil, err
	}

	res := []int{}

	for _, s := range sqs.Response {
		for _, p := range s.Players {
			res = append(res, p.ID)
		}
	}
	return &res, &sqs.Response[0].Players, nil
}

func GetPlayerById(id, season int) (*Player, error) {
	data, err := comm.GetHttpBody(playersIdURL, id, season)
	if err != nil {
		return nil, err
	}

	p := PlayerAPIResponse{}
	err = json.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}
	if len(p.Response) == 0 {
		return nil, fmt.Errorf("Player %d missing in Rapid API", id)
	}
	return &p.Response[0], nil
}
