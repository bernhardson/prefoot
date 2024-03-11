package database

import (
	"github.com/bernhardson/prefoot/data-fetch/pkg/fetch"
	"github.com/rs/zerolog"
)

type Repository struct {
	Logger zerolog.Logger
	Teams  interface {
		Insert(*TeamRow) (int64, error)
	}
	Venues interface {
		Insert(*VenueRow) (int64, error)
	}
	Players interface {
		Insert(*PlayerRow) (int64, error)
		InsertSeasonStats(*PlayerStatisticsRow) (int64, error)
		InsertStats(*PlayerStatsRow) (int64, error)
		SelectPlayersAndStatisticsByTeamId(int) (*[]*PlayersJoinOnPlayerStatsRow, error)
	}
	Fixture interface {
		Insert(*FixtureRow) (int64, error)
		InsertTeamsStats(*TeamStatisticsRow) (int64, error)
		InsertFormation(*FormationRow) (int64, error)
		SelectFixturesByRound(int) ([]*FixtureRow, error)
	}
	League interface {
		Insert(*fetch.League) (int64, error)
	}
	Coach interface {
		Insert(*CoachRow) (int64, error)
		InsertCareer(*CoachCareerRow) (int64, error)
	}
}
