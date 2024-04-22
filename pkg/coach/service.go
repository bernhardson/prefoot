package coach

import (
	"fmt"
	"time"

	"github.com/bernhardson/prefoot/pkg/team"
	"github.com/rs/zerolog"
)

type CoachModel struct {
	Logger *zerolog.Logger
	Repo   *CoachRepo
}

func (cm *CoachModel) FetchAndInsertCoaches(teams *[]team.TeamVenue) (*[]int, *[]int, error) {
	fc := []int{}
	fcc := []int{}

	for _, t := range *teams {
		//insert coach
		cs, err := GetCoach(t.Team.ID)
		if err != nil {
			cm.Logger.Err(err).Msg(fmt.Sprintf("could not get for team id %d", t.Team.ID))
		}

		for _, c := range cs.Response {
			_, err = cm.Repo.Insert(&CoachRow{
				ID:   c.ID,
				Name: c.Name,
			})
			if err != nil {
				cm.Logger.Err(err).Msg(fmt.Sprintf("%d", c.ID))
				fc = append(fc, c.ID)
			}
			for _, cc := range c.Career {
				//insert coach career
				start, err := time.Parse("2006-01-02", cc.Start)
				if err != nil {
					cm.Logger.Err(err).Msg(fmt.Sprintf("coach_%d", c.ID))
				}
				var end time.Time
				if cc.End != "" {
					end, err = time.Parse("2006-01-02", cc.End)
					if err != nil {
						cm.Logger.Err(err).Msg(fmt.Sprintf("coach_%d", c.ID))
					}
				} else {
					end = time.Time{}
				}
				_, err = cm.Repo.InsertCareer(&CoachCareerRow{
					CoachID: c.ID,
					TeamID:  cc.Team.ID,
					Start:   &start,
					End:     &end,
				})
				if err != nil {
					cm.Logger.Err(err).Msg(fmt.Sprintf("coach_%d#team%d", c.ID, cc.Team.ID))
					fcc = append(fcc, c.ID)
				} else {
					cm.Logger.Debug().Msg(fmt.Sprintf("inserted Coach %s with id %d", c.Name, c.ID))
				}
			}
		}
	}
	return &fc, &fcc, nil
}
