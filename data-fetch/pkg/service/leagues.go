package service

import (
	"github.com/bernhardson/prefoot/data-fetch/pkg/database"
	"github.com/bernhardson/prefoot/data-fetch/pkg/fetch"
	"github.com/rs/zerolog/log"
)

func FetchAndInsertLeagues(repo *database.Repository) (*[]int, error) {
	ls, err := fetch.GetLeagues()
	var failed []int
	if err != nil {
		log.Err(err).Msg("")
	}

	for _, l := range ls.Response {
		_, err := repo.League.Insert(&l.League)
		if err != nil {
			failed = append(failed, l.League.ID)
		}
	}
	return &failed, nil
}

func FetchAndInsertLeague(repo *database.Repository, league int) (*fetch.LeagueResponse, *[]int, error) {
	ls, err := fetch.GetLeague(league)
	var failed []int
	if err != nil {
		log.Err(err).Msg("")
	}

	for _, l := range ls.Response {
		_, err := repo.League.Insert(&l.League)
		if err != nil {
			failed = append(failed, l.League.ID)
		}
	}
	return ls, &failed, nil
}
