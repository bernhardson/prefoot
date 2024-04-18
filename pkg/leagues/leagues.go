package leagues

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LeaguesModel struct {
	Logger *zerolog.Logger
	Repo   interface {
		Insert(*League) (int64, error)
	}
}

func (lm *LeaguesModel) FetchAndInsertLeagues() (*[]int, error) {
	ls, err := GetLeagues()
	var failed []int
	if err != nil {
		log.Err(err).Msg("")
	}

	for _, l := range ls.Response {
		_, err := lm.Repo.Insert(&l.League)
		if err != nil {
			failed = append(failed, l.League.ID)
		}
	}
	return &failed, nil
}

func (lm *LeaguesModel) FetchAndInsertLeague(league int) (*LeagueResponse, *[]int, error) {
	ls, err := GetLeague(league)
	var failed []int
	if err != nil {
		log.Err(err).Msg("")
	}

	for _, l := range ls.Response {
		_, err := lm.Repo.Insert(&l.League)
		if err != nil {
			failed = append(failed, l.League.ID)
		}
	}
	return ls, &failed, nil
}
