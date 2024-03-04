package fetch

import (
	"github.com/bernhardson/prefoot/data-fetch/pkg/comm"
	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
	"github.com/mitchellh/mapstructure"
)

const (
	fixturesURL = "https://api-football-v1.p.rapidapi.com/v3/fixtures?league=%d&season=%d"

	fixtureDetailURL = "https://api-football-v1.p.rapidapi.com/v3/fixtures?id=%d"
)

func GetMatches(league int, season int) (*[]*model.Fixture, error) {
	data, err := comm.GetHttpBody(fixturesURL, league, season)
	if err != nil {
		return nil, err
	}

	matches := &[]*model.Fixture{}
	err = mapstructure.Decode(data["response"], matches)
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func GetFixtureDetail(id int) (*[]model.FixtureDetail, error) {
	data, err := comm.GetHttpBody(fixtureDetailURL, id)
	if err != nil {
		return nil, err
	}

	fd := &[]model.FixtureDetail{}
	err = mapstructure.Decode(data["response"], fd)
	if err != nil {
		return nil, err
	}
	return fd, nil

}
