package fetch

import (
	"encoding/json"

	"github.com/bernhardson/prefoot/data-fetch/pkg/comm"
)

const (
	coachURL = "https://api-football-v1.p.rapidapi.com/v3/coachs?team=%d"
)

type CoachResponse struct {
	Get        string        `json:"get"`
	Parameters interface{}   `json:"parameters"`
	Errors     []interface{} `json:"errors"`
	Results    int           `json:"results"`
	Paging     Paging        `json:"paging"`
	Response   []Coach       `json:"response"`
}

type Coach struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Age         int    `json:"age"`
	Birth       Birth  `json:"birth"`
	Nationality string `json:"nationality"`
	Height      string `json:"height"`
	Weight      string `json:"weight"`
	Photo       string `json:"photo"`
	Team        Team   `json:"team"`
	Career      Career `json:"career"`
}

type Career []struct {
	Team  Team   `json:"team"`
	Start string `json:"start"`
	End   string `json:"end"`
}

func GetCoach(id int) (*CoachResponse, error) {

	data, err := comm.GetHttpBodyRaw(coachURL, id)
	if err != nil {
		return nil, err
	}

	coach := CoachResponse{}
	err = json.Unmarshal(data, &coach)
	if err != nil {
		return nil, err
	}

	return &coach, nil
}
