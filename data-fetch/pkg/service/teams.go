package service

import (
	"strconv"

	"github.com/bernhardson/prefoot/data-fetch/pkg/database"
	"github.com/bernhardson/prefoot/data-fetch/pkg/fetch"
)

func FetchAndInsertTeams(repo *database.Repository, league int, season int) (*[]fetch.TeamVenue, error) {

	resp, err := fetch.GetTeams(league, season)
	if err != nil {
		return nil, err
	}
	//each team-venue struct
	for _, tv := range resp.TeamVenues {
		t := &database.TeamRow{
			Id:      tv.Team.ID,
			Name:    tv.Team.Name,
			Country: tv.Team.Country,
			Code:    tv.Team.Code,
		}
		row, err := repo.Teams.Insert(t)
		if err != nil {
			repo.Logger.Err(err).Msg(strconv.FormatInt(row, 10))
		}
		v := &database.VenueRow{
			Id:   tv.Venue.ID,
			Name: tv.Team.Name,
			City: tv.Venue.City,
		}
		row, err = repo.Venues.Insert(v)
		if err != nil {
			repo.Logger.Err(err).Msg(strconv.FormatInt(row, 10))
		}
	}
	return &resp.TeamVenues, nil
}
