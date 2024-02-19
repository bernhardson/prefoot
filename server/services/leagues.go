package services

import (
	"github.com/bernhardson/prefoot-init/pkg"
)

func GetStandings(league int, season int) (*[][]pkg.StandingsTeam, error) {
	s, err := pkg.GetStanding(league, season)
	if err != nil {
		return nil, err
	}
	return s, nil
}
