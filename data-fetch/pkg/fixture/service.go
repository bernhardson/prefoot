package fixture

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/bernhardson/prefoot/data-fetch/pkg/players"
	"github.com/bernhardson/prefoot/data-fetch/pkg/result"
	"github.com/bernhardson/prefoot/data-fetch/pkg/rounds"
)

type FixtureModel struct {
	Logger *zerolog.Logger
	Repo   interface {
		Insert(*FixtureRow) (int64, error)
		InsertTeamsStats(*TeamStatisticsRow) (int64, error)
		InsertFormation(*FormationRow) (int64, error)
		SelectFixturesByRound(int) ([]*FixtureRow, error)
		SelectFixtureByLeagueSeasonRound(int, int, int) ([]*FixtureRow, error)
		SelectFixtureIdsForLastNRounds(int, int, int, int) (*[]int, error)
		SelectLastNMatchups(int, int, int) ([]*FixtureRow, error)
		SelectLastNFixturesByTeam(int, int, int) ([]*FixtureRow, error)
		DeleteFixture(int) (int64, error)
	}

	RoundRepo  rounds.Repo
	PlayerRepo players.Repo
	ResultRepo result.ResultRepo
}

// Initialize fixtures, formations, team_statistics, player_statistics, rounds tables.
// Queries Rapid API then insert into local postgres.
// Some data manipulation is done on the fly.
func (fm *FixtureModel) FetchAndInsertFixtures(league, season int) error {

	fr, err := FetchFixtures(league, season)
	if err != nil {
		log.Err(err).Msg("")
		return err
	}

	for _, f := range fr.Response {
		round := 0
		//insert league
		fd, err := GetFixtureDetail(f.Fixture.ID)

		if err != nil {
			return err
		}
		round, err = strconv.Atoi(extractDigits(f.League.Round))
		if err != nil {
			fm.Logger.Err(err).Msg("")
		}
		fm.InsertFixture(&fd.FixtureDetail, league, season, round)
		fm.Logger.Info().Msg(fmt.Sprintf("Inserted fixture=%d", f.Fixture.ID))
	}

	return nil
}

func (fm *FixtureModel) UpdateFixture(league, season int) error {

	ts := time.Now().Unix()
	row, err := fm.RoundRepo.SelectLatestFinishedRound(league, season, ts)
	if err != nil {
		return err
	}
	fixtures, err := fm.Repo.SelectFixtureByLeagueSeasonRound(league, season, row.Round)
	if err != nil {
		return err
	}
	for _, f := range fixtures {
		fD, err := GetFixtureDetail(f.ID)
		if err != nil {
			return err
		}
		_, err = fm.Repo.DeleteFixture(f.ID)
		if err != nil {
			return err
		}
		fm.InsertFixture(&fD.FixtureDetail, league, season, f.Round)
	}
	return nil
}

// Loops at fixtures f and triggers their data base insert.
// since fixture details come with all kinds of match information such as
// lineups, player statistics etc. that are not part of the fixture table
// we insert those to database as well while the information is available
func (fm *FixtureModel) InsertFixture(fr *[]FixtureDetail, league, season, round int) {

	for _, fd := range *fr {
		start, err := fm.RoundRepo.SelectTimestampFromRounds(league, season, round)
		end := -1
		if err != nil {
			fm.Logger.Err(err).Msg("")
		}

		if start == -1 {
			start = fd.Fixture.Timestamp

		} else {
			if start > fd.Fixture.Timestamp {
				start = fd.Fixture.Timestamp
			}
		}

		if end == -1 {
			end = fd.Fixture.Timestamp

		} else {
			if end < fd.Fixture.Timestamp {
				end = fd.Fixture.Timestamp
			}
		}

		_, err = fm.RoundRepo.Insert(&rounds.RoundRow{
			Start:  int64(start),
			Round:  round,
			Season: season,
			League: league,
			End:    int64(end)})

		if err != nil {
			fm.Logger.Err(err).Msg("")
		}

		fm.Logger.Debug().Msg(fmt.Sprintf("insert fixture :%d", fd.Fixture.ID))
		//insert fixture
		_, err = fm.Repo.Insert(&FixtureRow{
			ID:            fd.Fixture.ID,
			League:        fd.League.ID,
			Round:         round,
			Referee:       fd.Fixture.Referee,
			Timezone:      fd.Fixture.Timezone,
			Timestamp:     fd.Fixture.Timestamp,
			Venue:         fd.Fixture.Venue.ID,
			Season:        season,
			HomeTeam:      fd.Teams.Home.ID,
			AwayTeam:      fd.Teams.Away.ID,
			HomeGoals:     fd.Goals.Home,
			AwayGoals:     fd.Goals.Away,
			HomeGoalsHalf: fd.Score.Halftime.Home,
			AwayGoalsHalf: fd.Score.Halftime.Away,
		})
		if err != nil {
			fm.Logger.Err(err).Msg(fmt.Sprintf("insert fixture: fixture_%d", fd.Fixture.ID))
		}

		if fd.Fixture.Status.Elapsed > 0 { //calculate and insert results
			home, away := calculateResult(&fd, league, season, round)
			fm.ResultRepo.Insert(home)
			fm.ResultRepo.Insert(away)

			for i, l := range fd.Lineups {
				ts := convertTeamStatistics(i, &fd, fm.Logger)
				_, err := fm.Repo.InsertTeamsStats(&TeamStatisticsRow{
					Team:           l.Team.ID,
					Fixture:        fd.Fixture.ID,
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
					fm.Logger.Err(err).Msg(fmt.Sprintf("insert team statistic: fixture_%d#team_%d", fd.Fixture.ID, l.Team.ID))
				}

				if len(l.Substitutes) == 5 {

					_, err = fm.Repo.InsertFormation(&FormationRow{
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
					_, err = fm.Repo.InsertFormation(&FormationRow{
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
					_, err = fm.Repo.InsertFormation(&FormationRow{
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
					fm.Logger.Err(err).Msg(fmt.Sprintf("insert formation: fixture_%d#team_%d", fd.Fixture.ID, l.Team.ID))
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
						fm.Logger.Err(err).Msg("")
					}
					accuracy, err := strconv.Atoi(strings.Replace(ps.Passes.Accuracy, "%", "", 1))
					if err != nil {
						fm.Logger.Err(err).Msg("")
					}
					_, err = fm.PlayerRepo.InsertStats(&players.PlayerStatsRow{
						Player:           player.Player.ID,
						Fixture:          fd.Fixture.ID,
						Team:             playerstats.Team.ID,
						League:           league,
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
						fm.Logger.Info().Msg(fmt.Sprintf("retrying player#%d", player.Player.ID))
						addMissingPlayer(fm.PlayerRepo, fm.Logger,
							season, player.Player.ID, t, ps.Games.Rating)
					} else if err != nil {
						fm.Logger.Err(err).Msg(fmt.Sprintf(
							"player statistics#player_%d#fixture_%d",
							player.Player.ID,
							fd.Fixture.ID))
					}
				}

			}
		}
	}
}

// uses goals to calculate win, draw, loss and adds the given points
func calculateResult(fd *FixtureDetail, league, season, round int) (*result.ResultRow, *result.ResultRow) {
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

	sHome := &result.ResultRow{
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

	sAway := &result.ResultRow{
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
func defaultStringValue(ps *PlayerStatisticsDetailsFD) {
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
func addMissingPlayer(repo players.Repo, logger *zerolog.Logger, season, id, team int, rating string) error {
	p, err := players.GetPlayerById(id, season)
	if err != nil {
		return err
	} else {
		_, err := repo.Insert(
			&players.PlayerRow{
				Id:           p.PlayerDetails.ID,
				Team:         team,
				Season:       season,
				FirstName:    p.PlayerDetails.FirstName,
				LastName:     p.PlayerDetails.LastName,
				BirthPlace:   p.PlayerDetails.Birth.Place,
				BirthCountry: p.PlayerDetails.Birth.Country,
				BirthDate:    p.PlayerDetails.Birth.Date,
			},
		)
		if err != nil {
			return err
		}
		//catch empty string ratin
		rating, err := strconv.ParseFloat(rating, 32)
		if err != nil {
			rating = 0
			logger.Debug().Msg(err.Error())
		}
		//insert season stats
		_, err = repo.InsertSeasonStats(&players.PlayerSeasonStatsRow{
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
			logger.Err(err).Msg("")
		}
	}
	return nil
}

// convert team statistics since they come as mixed type json
func convertTeamStatistics(home int, fd *FixtureDetail, logger *zerolog.Logger) *TeamStatisticsFD {
	var sd TeamStatisticsFD
	if len(fd.Statistics) > 0 {
		for _, s := range fd.Statistics[home].Statistics {
			val := 0
			floating, ok := s.Value.(float64)
			if ok && s.Value != nil {
				val = int(floating)
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
