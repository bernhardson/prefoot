package model

import (
	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
)

type PlayerSeason struct {
	Player     Player       `json:"player"`
	Statistics []Statistics `json:"statistics"`
}

type Birth struct {
	Date    string `json:"date"`
	Place   string `json:"place"`
	Country string `json:"country"`
}

type Player struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Age         int    `json:"age"`
	Birth       Birth  `json:"birth"`
	Nationality string `json:"nationality"`
	Height      string `json:"height"`
	Weight      string `json:"weight"`
	Injured     bool   `json:"injured"`
	Photo       string `json:"photo"`
}

type League struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Logo    string `json:"logo"`
	Flag    string `json:"flag"`
	Season  int    `json:"season"`
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

type PlayerGoals struct {
	Total    int `json:"total"`
	Conceded int `json:"conceded"`
	Assists  int `json:"assists"`
	Saves    int `json:"saves"`
}

type Passes struct {
	Total    int     `json:"total"`
	Key      int     `json:"key"`
	Accuracy float64 `json:"accuracy"`
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

type Fouls struct {
	Drawn     int `json:"drawn"`
	Committed int `json:"committed"`
}

type Cards struct {
	Yellow    int `json:"yellow"`
	Yellowred int `json:"yellowred"`
	Red       int `json:"red"`
}

type Statistics struct {
	Team        model.Team           `json:"team"`
	League      League               `json:"league"`
	Games       Games                `json:"games"`
	Substitutes Substitutes          `json:"substitutes"`
	Shots       Shots                `json:"shots"`
	Goals       PlayerGoals          `json:"goals"`
	Passes      Passes               `json:"passes"`
	Tackles     Tackles              `json:"tackles"`
	Duels       Duels                `json:"duels"`
	Dribbles    model.Dribbles       `json:"dribbles"`
	Fouls       Fouls                `json:"fouls"`
	Cards       Cards                `json:"cards"`
	Penalty     model.PenaltyDetails `json:"penalty"`
}

type Paging struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

// simple list of current players in a squad
type Squad struct {
	Team    model.Team `json:"team"`
	Players []Player   `json:"players"`
}
