package team

import (
	"strconv"

	"github.com/rs/zerolog"
)

type TeamModel struct {
	Logger   *zerolog.Logger
	TeamRepo interface {
		Insert(*TeamRow) (int64, error)
		Select(int) (*TeamRow, error)
		InsertTeamSeason(*TeamSeasonRow) (int64, error)
		SelectTeamsSeason(int, int) (*[]*TeamIds, error)
		SelectTeamsByIds(*[]int) (*[]*TeamRow, error)
	}
	VenuesRepo interface {
		Insert(*VenueRow) (int64, error)
	}
}

func (tm *TeamModel) FetchAndInsertTeams(league int, season int) (*[]TeamVenue, error) {

	resp, err := GetTeams(league, season)
	if err != nil {
		return nil, err
	}
	//each team-venue struct
	for _, tv := range resp.TeamVenues {
		t := &TeamRow{
			Id:      tv.Team.ID,
			Name:    tv.Team.Name,
			Country: tv.Team.Country,
			Code:    tv.Team.Code,
		}
		row, err := tm.TeamRepo.Insert(t)
		if err != nil {
			tm.Logger.Err(err).Msg(strconv.FormatInt(row, 10))
		}
		v := &VenueRow{
			Id:   tv.Venue.ID,
			Name: tv.Team.Name,
			City: tv.Venue.City,
		}
		row, err = tm.VenuesRepo.Insert(v)
		if err != nil {
			tm.Logger.Err(err).Msg(strconv.FormatInt(row, 10))
		}
		ts := &TeamSeasonRow{
			League: league,
			Season: season,
			Team:   tv.Team.ID,
		}
		row, err = tm.TeamRepo.InsertTeamSeason(ts)
		if err != nil {
			tm.Logger.Err(err).Msg(strconv.FormatInt(row, 10))
		}
	}
	return &resp.TeamVenues, nil
}
