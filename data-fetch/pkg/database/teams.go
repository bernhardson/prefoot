package database

import (
	"context"

	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	insertTeamStatistics = "INSERT INTO team_statistics " +
		"(team, fixture, shots_total, shots_on, shots_off, shots_blocked, " +
		"shots_box, shots_outside, offsides, fouls, corners, possession, yellow, red, " +
		"gk_saves, passes_total, passes_accurate, passes_percent, expected_goals) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)"

	insertFormation = `INSERT INTO formations (fixture, team, formation, player1, player2, player3, player4, player5, player6, player7, player8, player9, player10, player11, sub1, sub2, sub3, sub4, sub5, coach)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`
)

func InsertTeams(pool *pgxpool.Pool, l *model.Lineup, f *model.Fixture, ts *model.TeamStatistics) error {
	_, err := pool.Exec(
		context.Background(),
		insertTeamStatistics,
		l.Team.ID, f.ID, ts.ShotsTotal, ts.ShotsOn, ts.ShotsOff, ts.ShotsBlocked,
		ts.ShotsBox, ts.ShotsOutside, ts.Offsides, ts.Fouls, ts.Corners, ts.Possession, ts.Yellow, ts.Red,
		ts.GkSaves, ts.PassesTotal, ts.PassesAccurate, ts.PassesPercent, ts.ExpectedGoals,
	)
	return err
}

func InsertFormation(pool *pgxpool.Pool, l *model.Lineup, f *model.Fixture, ts *model.TeamStatistics) error {
	// insert formation
	_, err := pool.Exec(
		context.Background(),
		insertFormation,
		f.ID, l.Team.ID, l.Formation,
		l.StartXI[0].Player.ID, l.StartXI[1].Player.ID,
		l.StartXI[2].Player.ID, l.StartXI[3].Player.ID,
		l.StartXI[4].Player.ID, l.StartXI[5].Player.ID,
		l.StartXI[6].Player.ID, l.StartXI[7].Player.ID,
		l.StartXI[8].Player.ID, l.StartXI[9].Player.ID,
		l.StartXI[10].Player.ID,
		l.Substitutes[0].Player.ID, l.Substitutes[1].Player.ID,
		l.Substitutes[2].Player.ID, l.Substitutes[3].Player.ID,
		l.Substitutes[4].Player.ID,
		l.Coach.ID,
	)

	return err
}
