package leagues

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	insertLeague = `INSERT INTO leagues (id, name) VALUES ($1, $2);`
)

type LeagueRepo struct {
	Pool *pgxpool.Pool
}

type SeasonCoverage struct {
	Fixtures    SeasonCoverageFixtures `json:"fixtures"`
	Standings   bool                   `json:"standings"`
	Players     bool                   `json:"players"`
	TopScorers  bool                   `json:"top_scorers"`
	TopAssists  bool                   `json:"top_assists"`
	TopCards    bool                   `json:"top_cards"`
	Injuries    bool                   `json:"injuries"`
	Predictions bool                   `json:"predictions"`
	Odds        bool                   `json:"odds"`
}

type SeasonCoverageFixtures struct {
	Events             bool `json:"events"`
	Lineups            bool `json:"lineups"`
	StatisticsFixtures bool `json:"statistics_fixtures"`
	StatisticsPlayers  bool `json:"statistics_players"`
}

type StandingsPaging struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

func (lm *LeagueRepo) Insert(l *League) (int64, error) {
	row, err := lm.Pool.Exec(
		context.Background(),
		insertLeague,
		l.ID, l.Name)
	return row.RowsAffected(), err
}
