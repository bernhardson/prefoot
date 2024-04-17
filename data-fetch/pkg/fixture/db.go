package fixture

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	insertFixture = `INSERT INTO fixtures (id, league, round, referee, timezone, timestamp, venue, season, home_team, away_team,
						home_goals, away_goals, home_goals_half, away_goals_half)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);`
	deleteFixture        = `DELETE FROM fixtures WHERE id = $1;`
	insertTeamStatistics = "INSERT INTO team_statistics " +
		"(team, fixture, shots_total, shots_on, shots_off, shots_blocked, " +
		"shots_box, shots_outside, offsides, fouls, corners, possession, yellow, red, " +
		"gk_saves, passes_total, passes_accurate, passes_percent, expected_goals) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)"

	insertFormation = `INSERT INTO formations (fixture, team, formation, player1, player2, player3, player4, player5, player6, player7, player8, player9, player10, player11, sub1, sub2, sub3, sub4, sub5, coach)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`

	selectFixturesByRound             = "SELECT * FROM fixtures WHERE round = $1"
	selectFixturesByLeagueSeasonRound = `SELECT * FROM "fixtures" WHERE "league" = $1 AND "season" = $2 AND "round" = $3`
	selectFixturesByLastNRounds       = `SELECT id FROM fixtures WHERE league=$1 AND season=$2 AND round BETWEEN $3 and $4`

	selectLastNFixturesByTeams = `SELECT * FROM fixtures WHERE (home_team = $1 AND away_team = $2) OR (home_team = $2 AND away_team = $1) ORDER BY timestamp DESC LIMIT $3;`
	selectLastNFixturesByTeam  = `SELECT * FROM fixtures WHERE (home_team = $1 OR away_team = $1) AND timestamp < $2 ORDER BY timestamp DESC LIMIT $3;`
)

type FixtureRepo struct {
	Pool *pgxpool.Pool
}

type FixtureRow struct {
	ID            int    `json:"id"`
	League        int    `json:"league"`
	Round         int    `json:"round"`
	Referee       string `json:"referee"`
	Timezone      string `json:"timezone"`
	Timestamp     int    `json:"timestamp"`
	Venue         int    `json:"venue"`
	Season        int    `json:"season"`
	HomeTeam      int    `json:"home_team"`
	AwayTeam      int    `json:"away_team"`
	HomeGoals     int    `json:"home_goals"`
	AwayGoals     int    `json:"away_goals"`
	HomeGoalsHalf int    `json:"home_goals_half"`
	AwayGoalsHalf int    `json:"away_goals_half"`
}

func (fm *FixtureRepo) Insert(f *FixtureRow) (int64, error) {

	row, err := fm.Pool.Exec(
		context.Background(),
		insertFixture,
		f.ID, f.League, f.Round, f.Referee, f.Timezone,
		f.Timestamp, f.Venue, f.Season, f.HomeTeam,
		f.AwayTeam, f.HomeGoals, f.AwayGoals, f.HomeGoalsHalf,
		f.AwayGoalsHalf)

	return row.RowsAffected(), err
}

type TeamStatisticsRow struct {
	Team           int     `json:"team"`
	Fixture        int     `json:"fixture"`
	ShotsTotal     int     `json:"shots_total"`
	ShotsOn        int     `json:"shots_on"`
	ShotsOff       int     `json:"shots_off"`
	ShotsBlocked   int     `json:"shots_blocked"`
	ShotsBox       int     `json:"shots_box"`
	ShotsOutside   int     `json:"shots_outside"`
	Offsides       int     `json:"offsides"`
	Fouls          int     `json:"fouls"`
	Corners        int     `json:"corners"`
	Possession     int     `json:"possession"`
	Yellow         int     `json:"yellow"`
	Red            int     `json:"red"`
	GKSaves        int     `json:"gk_saves"`
	PassesTotal    int     `json:"passes_total"`
	PassesAccurate int     `json:"passes_accurate"`
	PassesPercent  int     `json:"passes_percent"`
	ExpectedGoals  float64 `json:"expected_goals"`
}

func (fm *FixtureRepo) InsertTeamsStats(t *TeamStatisticsRow) (int64, error) {
	row, err := fm.Pool.Exec(
		context.Background(),
		insertTeamStatistics,
		t.Team, t.Fixture, t.ShotsTotal, t.ShotsOn, t.ShotsOff, t.ShotsBlocked,
		t.ShotsBox, t.ShotsOutside, t.Offsides, t.Fouls, t.Corners, t.Possession, t.Yellow, t.Red,
		t.GKSaves, t.PassesTotal, t.PassesAccurate, t.PassesPercent, t.ExpectedGoals,
	)
	return row.RowsAffected(), err
}

type FormationRow struct {
	Fixture   int    `json:"fixture"`
	Team      int    `json:"team"`
	Formation string `json:"formation"`
	Player1   int    `json:"player1"`
	Player2   int    `json:"player2"`
	Player3   int    `json:"player3"`
	Player4   int    `json:"player4"`
	Player5   int    `json:"player5"`
	Player6   int    `json:"player6"`
	Player7   int    `json:"player7"`
	Player8   int    `json:"player8"`
	Player9   int    `json:"player9"`
	Player10  int    `json:"player10"`
	Player11  int    `json:"player11"`
	Sub1      int    `json:"sub1"`
	Sub2      int    `json:"sub2"`
	Sub3      int    `json:"sub3"`
	Sub4      int    `json:"sub4"`
	Sub5      int    `json:"sub5"`
	Coach     int    `json:"coach"`
}

func (fm *FixtureRepo) InsertFormation(f *FormationRow) (int64, error) {
	// insert formation
	row, err := fm.Pool.Exec(
		context.Background(),
		insertFormation,
		f.Fixture, f.Team, f.Formation,
		f.Player1, f.Player2, f.Player3, f.Player4, f.Player5,
		f.Player6, f.Player7, f.Player8, f.Player9, f.Player10,
		f.Player11, f.Sub1, f.Sub2, f.Sub3, f.Sub4, f.Sub5,
		f.Coach,
	)

	return row.RowsAffected(), err
}

func (pm *FixtureRepo) SelectFixturesByRound(round int) ([]*FixtureRow, error) {

	rows, err := pm.Pool.Query(
		context.Background(), selectFixturesByRound, round)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fixtures, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[FixtureRow])
	if err != nil {
		return nil, err
	}
	return fixtures, nil
}

func (pm *FixtureRepo) SelectFixtureByLeagueSeasonRound(league, season, round int) ([]*FixtureRow, error) {

	rows, err := pm.Pool.Query(
		context.Background(), selectFixturesByLeagueSeasonRound, league, season, round)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fixture, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[FixtureRow])
	if err != nil {
		return nil, err
	}
	return fixture, nil
}

func (fm *FixtureRepo) DeleteFixture(id int) (int64, error) {

	result, err := fm.Pool.Exec(context.Background(), deleteFixture, id)
	if err != nil {
		return -1, err
	}

	return result.RowsAffected(), nil
}

func (fm *FixtureRepo) SelectFixtureIdsForLastNRounds(league, season, round, n int) (*[]int, error) {

	var ret []int
	rows, err := fm.Pool.Query(context.Background(), selectFixturesByLastNRounds, league, season, round-n, round)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ret = append(ret, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &ret, err
}

func (pm *FixtureRepo) SelectLastNMatchups(team1, team2, n int) ([]*FixtureRow, error) {

	rows, err := pm.Pool.Query(
		context.Background(), selectLastNFixturesByTeams, team1, team2, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fixtures, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[FixtureRow])
	if err != nil {
		return nil, err
	}
	return fixtures, nil
}

func (pm *FixtureRepo) SelectLastNFixturesByTeam(team, ts, n int) ([]*FixtureRow, error) {

	rows, err := pm.Pool.Query(
		context.Background(), selectLastNFixturesByTeam, team, ts, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fixtures, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[FixtureRow])
	if err != nil {
		return nil, err
	}
	return fixtures, nil
}
