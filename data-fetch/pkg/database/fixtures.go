package database

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

const (
	insertFixture = `INSERT INTO fixtures (id, league, referee, timezone, timestamp, venue, season, home_team, away_team,
						home_goals, away_goals, home_goals_half, away_goals_half)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);`
)

func InsertFixture(pool *pgxpool.Pool, f *model.Fixture, fd *model.FixtureDetail, season int) error {

	_, err := pool.Exec(
		context.Background(),
		insertFixture,
		f.Fixture.ID, f.League.ID, f.Fixture.Referee, f.Fixture.Timezone,
		f.Fixture.Timestamp, f.Fixture.Venue.ID, season, f.Teams.Home.ID,
		f.Teams.Away.ID, f.Goals.Home, f.Goals.Away, fd.Score.Halftime.Home,
		fd.Score.Halftime.Away)

	return err
}

const (
	insertTeamStatistics = "INSERT INTO team_statistics " +
		"(team, fixture, shots_total, shots_on, shots_off, shots_blocked, " +
		"shots_box, shots_outside, offsides, fouls, corners, possession, yellow, red, " +
		"gk_saves, passes_total, passes_accurate, passes_percent, expected_goals) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)"

	insertFormation = `INSERT INTO formations (fixture, team, formation, player1, player2, player3, player4, player5, player6, player7, player8, player9, player10, player11, sub1, sub2, sub3, sub4, sub5, coach)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`
)

func InsertTeamsStatistics(pool *pgxpool.Pool, l *model.Lineup, f *model.FixtureMeta, ts *model.TeamStatistics) error {
	_, err := pool.Exec(
		context.Background(),
		insertTeamStatistics,
		l.Team.ID, f.ID, ts.ShotsTotal, ts.ShotsOn, ts.ShotsOff, ts.ShotsBlocked,
		ts.ShotsBox, ts.ShotsOutside, ts.Offsides, ts.Fouls, ts.Corners, ts.Possession, ts.Yellow, ts.Red,
		ts.GkSaves, ts.PassesTotal, ts.PassesAccurate, ts.PassesPercent, ts.ExpectedGoals,
	)
	return err
}

func InsertFormation(pool *pgxpool.Pool, l *model.Lineup, f *model.FixtureMeta, ts *model.TeamStatistics, logger zerolog.Logger) error {
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

	if err != nil {
		logger.Err(err).Msg(fmt.Sprintf("fixture_%d", f.ID))
	}

	return err
}
func InsertFixtures(pool *pgxpool.Pool, fr *[]model.FixtureDetail, f *model.Fixture, season int, logger zerolog.Logger) {
	for _, fd := range *fr {
		//insert fixture
		err := InsertFixture(pool, f, &fd, season)
		if err != nil {
			logger.Err(err).Msg(fmt.Sprintf("insert fixture: fixture_%d", f.Fixture.ID))
		}
		for i, l := range fd.Lineups {
			ts := getTeamStatistics(i, &fd, logger)
			err = InsertTeamsStatistics(pool, &l, &f.Fixture, ts)
			if err != nil {
				logger.Err(err).Msg(fmt.Sprintf("insert team statistic: fixture_%d#team_%d", f.Fixture.ID, l.Team.ID))
			}
			err = InsertFormation(pool, &l, &f.Fixture, ts, logger)
			if err != nil {
				logger.Err(err).Msg(fmt.Sprintf("insert formation fixture_%d#team_%d", f.Fixture.ID, l.Team.ID))
			}

		}
		for _, teamStats := range fd.Players {
			// insert player statistics
			t := teamStats.Team.ID
			for _, player := range teamStats.Players {
				ps := player.Statistics[0]
				if ps.Games.Rating == "" {
					ps.Games.Rating = "0"
				}
				if ps.Passes.Accuracy == "" {
					ps.Passes.Accuracy = "0"
				}
				err = InsertPlayerStatistic(pool, &player, &ps, &f.Fixture, t, season, logger)
				if err != nil {
					logger.Err(err).Msg(fmt.Sprintf(
						"player statistics#player_%d#fixture_%d",
						player.Player.ID,
						f.Fixture.ID))
				}
			}

		}
	}
}

// convert team statistics since they come as mixed type json
func getTeamStatistics(home int, fd *model.FixtureDetail, logger zerolog.Logger) *model.TeamStatistics {
	var sd model.TeamStatistics
	for _, s := range fd.Statistics[home].Statistics {
		val := 0
		f, ok := s.Value.(float64)
		if ok && s.Value != nil {
			val = int(f)
		}
		switch t := s.Type; t {
		case "Shots on Goal":
			sd.ShotsOn = val
		case "Shots off Goal":
			sd.ShotsOff = val
		case "Total Shots":
			sd.ShotsTotal = val
		case "Blocked Shots":
			sd.ShotsBlocked = val
		case "Shots insidebox":
			sd.ShotsBox = val
		case "Shots outsidebox":
			sd.ShotsOutside = val
		case "Fouls":
			sd.Fouls = val
		case "Corner Kicks":
			sd.Corners = val
		case "Offsides":
			sd.Offsides = val
		case "Ball Possession":
			p := strings.Replace(s.Value.(string), "%", "", -1)
			sd.Possession, _ = strconv.Atoi(p)
		case "Yellow Cards":
			sd.Yellow = val
		case "Red Cards":
			sd.Red = val
		case "Goalkeeper Saves":
			sd.GkSaves = val
		case "Total passes":
			sd.PassesTotal = val
		case "Passes accurate":
			sd.PassesAccurate = val
		case "Passes %":
			p := strings.Replace(s.Value.(string), "%", "", -1)
			sd.PassesPercent, _ = strconv.Atoi(p)
		case "expected_goals":
			fl, err := strconv.ParseFloat(s.Value.(string), 32)
			if err != nil {
				logger.Err(err).Msg("")
				sd.ExpectedGoals = 0
			} else {
				sd.ExpectedGoals = fl
			}
		}
	}
	return &sd
}
