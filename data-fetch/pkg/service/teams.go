package service

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"

	"github.com/bernhardson/prefoot/data-fetch/pkg/comm"
	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
)

const (
	insertTeam  = `INSERT INTO teams (id, name)VALUES ($1, $2)`
	insertVenue = `INSERT INTO venues (id, name, city) VALUES ($1, $2, $3)`
	teamsURL    = "https://api-football-v1.p.rapidapi.com/v3/teams?league=%d&season=%d"
)

func InsertTeams(teams *[]model.TeamVenue, conn *pgxpool.Pool) {
	//insert first the teams and venues
	//to ensure foreign key are set already
	for _, t := range *teams {

		//insert venue
		_, err := conn.Exec(
			context.Background(),
			insertVenue,
			t.Venue.ID, t.Venue.Name, t.Venue.City,
		)
		if err != nil {
			log.Err(err).Msg("")
		}
		//insert team
		_, err = conn.Exec(
			context.Background(),
			insertTeam,
			t.Team.ID,
			t.Team.Name,
		)
		if err != nil {
			log.Err(err).Msg("")
		}
	}
}

func GetTeams(league int, season int) (*[]model.TeamVenue, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(teamsURL, league, season), nil)
	if err != nil {
		return nil, err
	}

	req = comm.AddRequestHeader(req)
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
	data, err := comm.UnmarshalData(body)
	if err != nil {
		return nil, err
	}

	tv := &[]model.TeamVenue{}
	mapstructure.Decode(data["response"], &tv)

	return tv, nil
}
