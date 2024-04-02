package service

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/bernhardson/prefoot/data-fetch/pkg/database"
	"github.com/bernhardson/prefoot/data-fetch/pkg/fetch"
)

func FetchAndInsertFixtures(repo *database.Repository, league, season int) error {

	fr, err := fetch.FetchFixtures(league, season)
	if err != nil {
		log.Err(err).Msg("")
		return err
	}

	for _, f := range fr.Response {
		round := 0
		//insert league
		fd, err := fetch.GetFixtureDetail(f.Fixture.ID)

		if err != nil {
			return err
		}
		round, err = strconv.Atoi(extractDigits(f.League.Round))
		if err != nil {
			repo.Logger.Err(err).Msg("")
		}
		InsertFixture(repo, &fd.FixtureDetail, &f, league, season, round)
		repo.Logger.Info().Msg(fmt.Sprintf("Inserted fixture=%d", f.Fixture.ID))
	}

	return nil
}

// Loops at fixtures f and triggers their data base insert.
// since fixture details come with all kinds of match information such as
// lineups, player statistics etc. that are not part of the fixture table
// we insert those to database as well while the information is available
func InsertFixture(repo *database.Repository, fr *[]fetch.FixtureDetail, f *fetch.Fixture, league, season, round int) {

	for _, fd := range *fr {

		start, err := repo.Fixture.SelectTimestampFromRounds(league, season, round)
		if err != nil {
			repo.Logger.Err(err).Msg("")
		}
		if start == -1 {
			start = f.Fixture.Timestamp

		} else {
			if start > f.Fixture.Timestamp {
				start = f.Fixture.Timestamp
			}
		}
		_, err = repo.Fixture.InsertRound(&database.RoundRow{
			Start:  int64(start),
			Round:  round,
			Season: season,
			League: league})

		if err != nil {
			repo.Logger.Err(err).Msg("")
		}

		repo.Logger.Debug().Msg(fmt.Sprintf("insert fixture :%d", fd.Fixture.ID))
		//insert fixture
		_, err = repo.Fixture.Insert(&database.FixtureRow{
			ID:            f.Fixture.ID,
			League:        f.League.ID,
			Round:         round,
			Referee:       f.Fixture.Referee,
			Timezone:      fd.Fixture.Timezone,
			Timestamp:     fd.Fixture.Timestamp,
			Venue:         f.Fixture.Venue.ID,
			Season:        season,
			HomeTeam:      f.Teams.Home.ID,
			AwayTeam:      fd.Teams.Away.ID,
			HomeGoals:     f.Goals.Home,
			AwayGoals:     f.Goals.Away,
			HomeGoalsHalf: fd.Score.Halftime.Home,
			AwayGoalsHalf: fd.Score.Halftime.Away,
		})
		if err != nil {
			repo.Logger.Err(err).Msg(fmt.Sprintf("insert fixture: fixture_%d", f.Fixture.ID))
		}

		if fd.Fixture.Status.Elapsed > 0 { //calculate and insert results
			home, away := calculateResult(&fd, league, season, round)
			repo.Result.Insert(home)
			repo.Result.Insert(away)

			for i, l := range fd.Lineups {
				ts := convertTeamStatistics(i, &fd, &repo.Logger)
				_, err := repo.Fixture.InsertTeamsStats(&database.TeamStatisticsRow{
					Team:           l.Team.ID,
					Fixture:        f.Fixture.ID,
					ShotsTotal:     ts.ShotsTotal,
					ShotsOn:        ts.ShotsOn,
					ShotsOff:       ts.ShotsOff,
					ShotsBlocked:   ts.ShotsBlocked,
					ShotsBox:       ts.ShotsBox,
					ShotsOutside:   ts.ShotsOutside,
					Offsides:       ts.Offsides,
					Fouls:          ts.Fouls,
					Corners:        ts.Corners,
					Possession:     ts.Possession,
					Yellow:         ts.Yellow,
					Red:            ts.Red,
					GKSaves:        ts.GkSaves,
					PassesTotal:    ts.PassesTotal,
					PassesAccurate: ts.PassesAccurate,
					PassesPercent:  ts.PassesPercent,
					ExpectedGoals:  ts.ExpectedGoals,
				})
				if err != nil {
					repo.Logger.Err(err).Msg(fmt.Sprintf("insert team statistic: fixture_%d#team_%d", f.Fixture.ID, l.Team.ID))
				}

				if len(l.Substitutes) == 5 {

					_, err = repo.Fixture.InsertFormation(&database.FormationRow{
						Fixture:   fd.Fixture.ID,
						Team:      l.Team.ID,
						Formation: l.Formation,
						Player1:   l.StartXI[0].Player.ID,
						Player2:   l.StartXI[1].Player.ID,
						Player3:   l.StartXI[2].Player.ID,
						Player4:   l.StartXI[3].Player.ID,
						Player5:   l.StartXI[4].Player.ID,
						Player6:   l.StartXI[5].Player.ID,
						Player7:   l.StartXI[6].Player.ID,
						Player8:   l.StartXI[7].Player.ID,
						Player9:   l.StartXI[8].Player.ID,
						Player10:  l.StartXI[9].Player.ID,
						Player11:  l.StartXI[10].Player.ID,
						Sub1:      l.Substitutes[0].Player.ID,
						Sub2:      l.Substitutes[1].Player.ID,
						Sub3:      l.Substitutes[2].Player.ID,
						Sub4:      l.Substitutes[3].Player.ID,
						Sub5:      l.Substitutes[4].Player.ID,
						Coach:     l.Coach.ID,
					})
				} else if len(l.Substitutes) == 4 {
					_, err = repo.Fixture.InsertFormation(&database.FormationRow{
						Fixture:   fd.Fixture.ID,
						Team:      l.Team.ID,
						Formation: l.Formation,
						Player1:   l.StartXI[0].Player.ID,
						Player2:   l.StartXI[1].Player.ID,
						Player3:   l.StartXI[2].Player.ID,
						Player4:   l.StartXI[3].Player.ID,
						Player5:   l.StartXI[4].Player.ID,
						Player6:   l.StartXI[5].Player.ID,
						Player7:   l.StartXI[6].Player.ID,
						Player8:   l.StartXI[7].Player.ID,
						Player9:   l.StartXI[8].Player.ID,
						Player10:  l.StartXI[9].Player.ID,
						Player11:  l.StartXI[10].Player.ID,
						Sub1:      l.Substitutes[0].Player.ID,
						Sub2:      l.Substitutes[1].Player.ID,
						Sub3:      l.Substitutes[2].Player.ID,
						Sub4:      l.Substitutes[3].Player.ID,
						Coach:     l.Coach.ID,
					})
				} else if len(l.Substitutes) == 3 {
					_, err = repo.Fixture.InsertFormation(&database.FormationRow{
						Fixture:   fd.Fixture.ID,
						Team:      l.Team.ID,
						Formation: l.Formation,
						Player1:   l.StartXI[0].Player.ID,
						Player2:   l.StartXI[1].Player.ID,
						Player3:   l.StartXI[2].Player.ID,
						Player4:   l.StartXI[3].Player.ID,
						Player5:   l.StartXI[4].Player.ID,
						Player6:   l.StartXI[5].Player.ID,
						Player7:   l.StartXI[6].Player.ID,
						Player8:   l.StartXI[7].Player.ID,
						Player9:   l.StartXI[8].Player.ID,
						Player10:  l.StartXI[9].Player.ID,
						Player11:  l.StartXI[10].Player.ID,
						Sub1:      l.Substitutes[0].Player.ID,
						Sub2:      l.Substitutes[1].Player.ID,
						Sub3:      l.Substitutes[2].Player.ID,
						Coach:     l.Coach.ID,
					})
				}
				if err != nil {
					repo.Logger.Err(err).Msg(fmt.Sprintf("insert formation: fixture_%d#team_%d", f.Fixture.ID, l.Team.ID))
				}
			}
			for _, playerstats := range fd.Players {
				// insert player statistics
				t := playerstats.Team.ID
				for _, player := range playerstats.Players {
					ps := player.Statistics[0]
					defaultStringValue(&ps)
					rating, err := strconv.ParseFloat(ps.Games.Rating, 64)
					if err != nil {
						repo.Logger.Err(err).Msg("")
					}
					accuracy, err := strconv.Atoi(strings.Replace(ps.Passes.Accuracy, "%", "", 1))
					if err != nil {
						repo.Logger.Err(err).Msg("")
					}
					_, err = repo.Players.InsertStats(&database.PlayerStatsRow{
						PlayerID:         player.Player.ID,
						FixtureID:        fd.Fixture.ID,
						TeamID:           playerstats.Team.ID,
						Season:           season,
						Minutes:          ps.Games.Minutes,
						Position:         ps.Games.Position,
						Captain:          ps.Games.Captain,
						Rating:           rating,
						Substitute:       ps.Games.Substitute,
						ShotsTotal:       ps.Shots.Total,
						ShotsOn:          ps.Shots.On,
						GoalsScored:      ps.Goals.Total,
						GoalsAssisted:    ps.Goals.Assists,
						PassesTotal:      ps.Passes.Total,
						PassesKey:        ps.Passes.Key,
						Accuracy:         accuracy,
						Tackles:          ps.Tackles.Total,
						Block:            ps.Tackles.Blocks,
						Interceptions:    ps.Tackles.Interceptions,
						DuelsTotal:       ps.Duels.Total,
						DuelsWon:         ps.Duels.Won,
						DribblesTotal:    ps.Dribbles.Attempts,
						DribblesWon:      ps.Dribbles.Success,
						Yellow:           ps.Cards.Yellow,
						Red:              ps.Cards.Red,
						PenaltyWon:       ps.Penalty.Won,
						PenaltyCommitted: ps.Penalty.Commited,
						PenaltyScored:    ps.Penalty.Scored,
						PenaltyMissed:    ps.Penalty.Missed,
						PenaltySaved:     ps.Penalty.Saved,
						Saves:            ps.Goals.Saves,
					})

					if err != nil && strings.HasSuffix(err.Error(), "(SQLSTATE 23503)") {
						repo.Logger.Info().Msg(fmt.Sprintf("retrying player#%d", player.Player.ID))
						addMissingPlayer(repo, season, player.Player.ID, t, ps.Games.Rating)
					} else if err != nil {
						repo.Logger.Err(err).Msg(fmt.Sprintf(
							"player statistics#player_%d#fixture_%d",
							player.Player.ID,
							f.Fixture.ID))
					}
				}

			}
		}
	}
}

// uses goals to calculate win, draw, loss and adds the given points
func calculateResult(fd *fetch.FixtureDetail, league, season, round int) (*database.ResultRow, *database.ResultRow) {
	hPoints := 0
	aPoints := 0
	if fd.Fixture.Status.Elapsed != 0 {
		if fd.Teams.Home.Winner {
			hPoints = 3
		} else if fd.Teams.Away.Winner {
			aPoints = 3
		} else {
			hPoints = 1
			aPoints = 1
		}
	}

	sHome := &database.ResultRow{
		Team:         fd.Teams.Home.ID,
		League:       league,
		Round:        round,
		Season:       season,
		Points:       hPoints,
		GoalsFor:     fd.Goals.Home,
		GoalsAgainst: fd.Goals.Away,
		Modus:        1,
		Elapsed:      fd.Fixture.Status.Elapsed,
	}

	sAway := &database.ResultRow{
		Team:         fd.Teams.Away.ID,
		League:       league,
		Round:        round,
		Season:       season,
		Points:       aPoints,
		GoalsFor:     fd.Goals.Away,
		GoalsAgainst: fd.Goals.Home,
		Modus:        2,
		Elapsed:      fd.Fixture.Status.Elapsed,
	}

	return sHome, sAway
}

// overrides empty string to "0"
func defaultStringValue(ps *fetch.PlayerStatisticsDetailsFD) {
	if ps.Games.Rating == "" {
		ps.Games.Rating = "0"
	}
	if ps.Passes.Accuracy == "" {
		ps.Passes.Accuracy = "0"
	}
}

// Players are inserted into the database before fixtures.
// However the rapid api plaeyer endpoint is missing player entries.
// If those appear during fixture insertion. Get the player detail and insert it
// separately.
func addMissingPlayer(env *database.Repository, season, id, team int, rating string) {
	p, err := fetch.GetPlayerById(id, season)
	if err != nil {
		env.Logger.Err(err).Msg(fmt.Sprintf("player_id=%d", id))
		return
	} else {
		_, err := env.Players.Insert(
			&database.PlayerRow{
				Id:           p.PlayerDetails.ID,
				TeamID:       team,
				Season:       season,
				FirstName:    p.PlayerDetails.FirstName,
				LastName:     p.PlayerDetails.LastName,
				BirthPlace:   p.PlayerDetails.Birth.Place,
				BirthCountry: p.PlayerDetails.Birth.Country,
				BirthDate:    p.PlayerDetails.Birth.Date,
			},
		)
		if err != nil {
			env.Logger.Err(err).Msg("")
		}
		//catch empty string ratin
		rating, err := strconv.ParseFloat(rating, 32)
		if err != nil {
			rating = 0
			env.Logger.Debug().Msg(err.Error())
		}
		//insert season stats
		_, err = env.Players.InsertSeasonStats(&database.PlayerSeasonStatsRow{
			PlayerID:           p.PlayerDetails.ID,
			Season:             season,
			TeamID:             team,
			Minutes:            p.Statistics[0].Games.Minutes,
			Position:           p.Statistics[0].Games.Position,
			Rating:             rating,
			Captain:            p.Statistics[0].Games.Captain,
			Appearances:        p.Statistics[0].Games.Number,
			Lineups:            p.Statistics[0].Games.Lineups,
			TotalShots:         p.Statistics[0].Shots.Total,
			ShotsOnTarget:      p.Statistics[0].Shots.On,
			TotalGoals:         p.Statistics[0].Goals.Total,
			Assists:            p.Statistics[0].Goals.Assists,
			TotalPasses:        p.Statistics[0].Passes.Total,
			KeyPasses:          p.Statistics[0].Passes.Key,
			PassAccuracy:       p.Statistics[0].Passes.Accuracy,
			TotalTackles:       p.Statistics[0].Tackles.Total,
			TackleBlocks:       p.Statistics[0].Tackles.Blocks,
			Interceptions:      p.Statistics[0].Tackles.Interceptions,
			TotalDuels:         p.Statistics[0].Duels.Total,
			DuelsWon:           p.Statistics[0].Duels.Won,
			DribbleAttempts:    p.Statistics[0].Dribbles.Attempts,
			DribbleSuccess:     p.Statistics[0].Dribbles.Success,
			YellowCards:        p.Statistics[0].Cards.Yellow,
			RedCards:           p.Statistics[0].Cards.Red,
			PenaltiesWon:       p.Statistics[0].Penalty.Won,
			PenaltiesCommitted: p.Statistics[0].Penalty.Committed,
			PenaltiesScored:    p.Statistics[0].Penalty.Scored,
			PenaltiesMissed:    p.Statistics[0].Penalty.Missed,
			PenaltiesSaved:     p.Statistics[0].Penalty.Saved,
			GoalkeeperSaves:    p.Statistics[0].Goals.Saves,
		})
		if err != nil {
			env.Logger.Err(err).Msg("")
		}
	}
}

// convert team statistics since they come as mixed type json
func convertTeamStatistics(home int, fd *fetch.FixtureDetail, logger *zerolog.Logger) *fetch.TeamStatisticsFD {
	var sd fetch.TeamStatisticsFD
	if len(fd.Statistics) > 0 {
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
				if s.Value != nil {
					p := strings.Replace(s.Value.(string), "%", "", -1)
					sd.Possession, _ = strconv.Atoi(p)
				} else {
					logger.Err(fmt.Errorf("Ball possession of team=%d#fixture=%d is null", home, fd.Fixture.ID))
				}
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
				if s.Value != nil {
					p := strings.Replace(s.Value.(string), "%", "", 1)
					sd.PassesPercent, _ = strconv.Atoi(p)
				} else {
					logger.Err(fmt.Errorf("Passes percentage of team=%d#fixture=%d is null", home, fd.Fixture.ID))
				}
			case "expected_goals":
				if s.Value != nil {
					fl, err := strconv.ParseFloat(s.Value.(string), 32)
					if err != nil {
						sd.ExpectedGoals = 0
					} else {
						sd.ExpectedGoals = fl
					}
				}
			}
		}
	}
	return &sd
}

func extractDigits(input string) string {
	// Regular expression to match digits
	regex := regexp.MustCompile(`\d+`)

	// Find all matches in the input string
	matches := regex.FindAllString(input, -1)

	// Join the matches into a single string
	digits := strings.Join(matches, "")

	return digits
}
