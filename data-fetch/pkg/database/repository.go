package database

import (
	"github.com/bernhardson/prefoot/data-fetch/pkg/fetch"
	"github.com/rs/zerolog"
)

type Repository struct {
	Logger zerolog.Logger
	Teams  interface {
		Insert(*TeamRow) (int64, error)
		Select(int) (*TeamRow, error)
		InsertTeamSeason(*TeamSeasonRow) (int64, error)
		SelectTeamsSeason(int, int) (*[]*TeamIds, error)
		SelectTeamsByIds(*[]int) (*[]*TeamRow, error)
	}
	Venues interface {
		Insert(*VenueRow) (int64, error)
	}
	Players interface {
		Insert(*PlayerRow) (int64, error)
		InsertSeasonStats(*PlayerSeasonStatsRow) (int64, error)
		InsertStats(*PlayerStatsRow) (int64, error)
		SelectPlayersAndStatisticsByTeamId(int) (*[]*PlayersJoinOnPlayerStatsRow, error)
		SelectPlayersByTeamId(int) ([]*PlayerRow, error)
		SelectPlayersByTeamLeagueSeason(int, int) ([]*PlayerRow, error)
		SelectPlayerStatisticsByPlayersFixturesTeam([]int, *[]int) ([]*KeyPlayerStats, error)
		SelectPlayerIdsBySeasonAndTeamId(int, int) ([]int, error)
	}
	Fixture interface {
		Insert(*FixtureRow) (int64, error)
		InsertTeamsStats(*TeamStatisticsRow) (int64, error)
		InsertFormation(*FormationRow) (int64, error)
		InsertRound(*RoundRow) (int64, error)
		SelectFixturesByRound(int) ([]*FixtureRow, error)
		SelectFixtureByLeagueSeasonRound(int, int, int) ([]*FixtureRow, error)
		SelectTimestampFromRounds(int, int, int) (int, error)
		SelectRoundByTimestamp(int, int, int64) (*RoundRow, error)
		SelectLatestFinishedRound(int, int, int64) (*RoundRow, error)
		SelectFixtureIdsForLastNRounds(int, int, int, int) (*[]int, error)
		DeleteFixture(int) (int64, error)
	}
	League interface {
		Insert(*fetch.League) (int64, error)
	}
	Coach interface {
		Insert(*CoachRow) (int64, error)
		InsertCareer(*CoachCareerRow) (int64, error)
	}

	Result interface {
		Insert(*ResultRow) (int64, error)
		Select(int) (*ResultRow, error)
		SelectByLeagueAndSeason(int, int) (*[]*ResultRow, error)
	}
}
