package service

import (
	"errors"
	"fmt"

	"github.com/bernhardson/prefoot/data-fetch/pkg/comm"
	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
	"github.com/jackc/pgx"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

const (
	getLeaguesURL  = "https://api-football-v1.p.rapidapi.com/v3/leagues"
	getStandingURL = "https://api-football-v1.p.rapidapi.com/v3/standings?season=%d&league=%d"
	insertLeague   = `INSERT INTO leagues (id, name) VALUES ($1, $2);`
)

func InsertLeagues(conn *pgx.Conn) (*[]int, *[]int, error) {
	ls, err := getLeagues()
	if err != nil {
		log.Err(err).Msg("")
	}
	insertedLeague := []int{}
	failedInsert := []int{}

	for _, l := range *ls {
		_, err = conn.Exec(
			insertLeague,
			l.League.ID, l.League.Name)
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("%d", l.League.ID))
			failedInsert = append(failedInsert, l.League.ID)
		} else {
			insertedLeague = append(insertedLeague, l.League.ID)
		}
	}
	return &insertedLeague, &failedInsert, nil
}

func getLeagues() (*[]model.LeagueData, error) {
	data, err := comm.GetHttpBody(getLeaguesURL)
	if err != nil {
		return nil, err
	}
	l := &[]model.LeagueData{}
	mapstructure.Decode(data["response"], l)

	return l, nil
}

func GetStanding(league int, season int) (*[][]model.StandingsTeam, error) {
	data, err := comm.GetHttpBody(getStandingURL, season, league)
	if err != nil {
		return nil, err
	}
	ses := []model.StandingsEntry{}
	mapstructure.Decode(data["response"], &ses)
	for _, se := range ses {
		return &se.League.Standings, nil
	}
	return nil, errors.New("wtf")
}
