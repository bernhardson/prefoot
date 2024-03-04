package database

import (
	"context"
	"fmt"
	"time"

	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

const (
	insertCoach       = `INSERT INTO coaches (id, name) VALUES ($1, $2)`
	insertCoachCareer = `INSERT INTO coach_careers (coach, team, start, "end") VALUES ($1, $2, $3, $4)`
)

func InsertCoaches(pool *pgxpool.Pool, cs *[]model.Coach, logger zerolog.Logger) (*[]int, *[]int) {

	insertSuccess := []int{}
	insertFail := []int{}

	for _, c := range *cs {
		_, err := pool.Exec(
			context.Background(),
			insertCoach,
			c.ID, c.Name,
		)
		if err != nil {
			logger.Err(err).Msg(fmt.Sprintf("%d", c.ID))
		}
		for _, cc := range c.Career {
			//insert coach career
			start, err := time.Parse("2006-01-02", cc.Start)
			if err != nil {
				logger.Err(err).Msg(fmt.Sprintf("coach_%d", c.ID))
			}
			var end time.Time
			if cc.End != "" {
				end, err = time.Parse("2006-01-02", cc.End)
				if err != nil {
					logger.Err(err).Msg(fmt.Sprintf("coach_%d", c.ID))
				}
			} else {
				end = time.Time{}
			}

			_, err = pool.Exec(
				context.Background(),
				insertCoachCareer,
				c.ID, cc.Team.ID, start, end,
			)
			if err != nil {
				logger.Err(err).Msg(fmt.Sprintf("coach_%d#team%d", c.ID, cc.Team.ID))
				insertFail = append(insertFail, c.ID)
			} else {
				logger.Debug().Msg(fmt.Sprintf("inserted Coach %s with id %d", c.Name, c.ID))
				insertSuccess = append(insertSuccess, c.ID)
			}
		}
	}
	return &insertSuccess, &insertFail
}
