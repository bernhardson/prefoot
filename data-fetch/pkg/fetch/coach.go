package fetch

import (
	"github.com/bernhardson/prefoot/data-fetch/pkg/comm"
	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
	"github.com/mitchellh/mapstructure"
)

const (
	coachURL = "https://api-football-v1.p.rapidapi.com/v3/coachs?team=%d"
)

func GetCoach(id int) (*[]model.Coach, error) {

	data, err := comm.GetHttpBody(coachURL, id)
	if err != nil {
		return nil, err
	}
	coach := &[]model.Coach{}
	mapstructure.Decode(data["response"], &coach)

	return coach, nil
}
