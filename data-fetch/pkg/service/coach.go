package service

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jackc/pgx"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"

	"github.com/bernhardson/prefoot/data-fetch/pkg/comm"
	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
)

const (
	coachURL          = "https://api-football-v1.p.rapidapi.com/v3/coachs?team=%d"
	insertCoach       = `INSERT INTO coaches (id, name) VALUES ($1, $2)`
	insertCoachCareer = `INSERT INTO coach_careers (id, team, start, "end") VALUES ($1, $2, $3, $4)`
)

func InsertCoaches(teams *[]model.TeamVenue, conn *pgx.Conn) (*[]int, *[]int, error) {
	insertSuccess := []int{}
	insertFail := []int{}

	for _, t := range *teams {
		//insert coach
		cs, err := GetCoach(t.Team.ID)
		if err != nil {
			return nil, nil, err
		}
		for _, c := range *cs {
			_, err = conn.Exec(
				insertCoach,
				c.ID, c.Name,
			)
			if err != nil {
				log.Err(err).Msg(fmt.Sprintf("%d", c.ID))
			}
			for _, cc := range c.Career {
				//insert coach career
				start, err := time.Parse("2006-01-02", cc.Start)
				if err != nil {
					log.Err(err).Msg(fmt.Sprintf("coach_%d", c.ID))
				}
				var end time.Time
				if cc.End != "" {
					end, err = time.Parse("2006-01-02", cc.End)
					if err != nil {
						log.Err(err).Msg(fmt.Sprintf("coach_%d", c.ID))
					}
				} else {
					end = time.Time{}
				}

				_, err = conn.Exec(

					insertCoachCareer,
					c.ID, cc.Team.ID, start, end,
				)
				if err != nil {
					log.Err(err).Msg(fmt.Sprintf("coach_%d#team%d", c.ID, cc.Team.ID))
					insertFail = append(insertFail, c.ID)
				} else {
					log.Debug().Msg(fmt.Sprintf("inserted Coach %s with id %d", c.Name, c.ID))
					insertSuccess = append(insertSuccess, c.ID)
				}
			}
		}
	}
	return &insertSuccess, &insertFail, nil
}

func GetCoach(id int) (*[]model.Coach, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(coachURL, id), nil)
	if err != nil {
		return nil, err
	}

	req = comm.AddRequestHeader(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	data, err := comm.UnmarshalData(body)
	if err != nil {
		return nil, err
	}

	coach := &[]model.Coach{}
	mapstructure.Decode(data["response"], &coach)

	return coach, nil
}
