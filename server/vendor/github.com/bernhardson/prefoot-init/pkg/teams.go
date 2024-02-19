package pkg

import (
	"fmt"
	"io"
	"net/http"

	"github.com/jackc/pgx"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

const (
	insertTeam = `INSERT INTO teams (id, name)VALUES ($1, $2)`
	teamsURL   = "https://api-football-v1.p.rapidapi.com/v3/teams?league=%d&season=%d"
)

type TeamVenue struct {
	Team  Team  `json:"team"`
	Venue Venue `json:"venue"`
}

type Team struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Country  string `json:"country"`
	Founded  int    `json:"founded"`
	National bool   `json:"national"`
	Logo     string `json:"logo"`
}

type Venue struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	City     string `json:"city"`
	Capacity int    `json:"capacity"`
	Surface  string `json:"surface"`
	Image    string `json:"image"`
}

func InsertTeams(teams *[]TeamVenue, conn *pgx.Conn) {
	//insert first the teams and venues
	//to ensure foreign key are set already
	for _, t := range *teams {

		//insert venue
		_, err := conn.Exec(
			insertVenue,
			t.Venue.ID, t.Venue.Name, t.Venue.City,
		)
		if err != nil {
			log.Err(err).Msg("")
		}
		//insert team
		_, err = conn.Exec(
			insertTeam,
			t.Team.ID,
			t.Team.Name,
		)
		if err != nil {
			log.Err(err).Msg("")
		}
	}
}

func GetTeams(league int, season int) (*[]TeamVenue, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(teamsURL, league, season), nil)
	if err != nil {
		return nil, err
	}

	req = AddRequestHeader(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	data, err := UnmarshalData(body)
	if err != nil {
		return nil, err
	}

	tv := &[]TeamVenue{}
	mapstructure.Decode(data["response"], &tv)

	return tv, nil
}
