package players

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

type PlayerModel struct {
	Logger *zerolog.Logger
	Repo   interface {
		Insert(*PlayerRow) (int64, error)
		InsertSeasonStats(*PlayerSeasonStatsRow) (int64, error)
		InsertStats(*PlayerStatsRow) (int64, error)
		SelectPlayersAndStatisticsByTeamId(int) (*[]*PlayersJoinOnPlayerStatsRow, error)
		SelectPlayersByTeamId(int) ([]*PlayerRow, error)
		SelectPlayersByTeamLeagueSeason(int, int) ([]*PlayerRow, error)
		SelectPlayerStatisticsByPlayersFixturesTeam([]int, *[]int) ([]*KeyPlayerStats, error)
		SelectPlayerIdsBySeasonAndTeamId(int, int) ([]int, error)
	}
}

func (pm *PlayerModel) FetchAndInsertPlayers(league int, season int) (*[]int, *[]int, error) {

	pgTotal := 1
	pgCurrent := 1
	var failedP, failedS []int
	for i := 1; pgCurrent <= pgTotal; i++ {
		ps, pg, err := GetPlayers(league, season, pgCurrent)

		if err != nil {
			return nil, nil, err
		}

		for _, p := range *ps {
			// player statistics only has one entry so there will be just one insert to player table
			for _, s := range p.Statistics {
				row, err := pm.Repo.Insert(
					&PlayerRow{
						Id:           p.PlayerDetails.ID,
						Team:         s.Team.ID,
						Season:       season,
						FirstName:    p.PlayerDetails.FirstName,
						LastName:     p.PlayerDetails.LastName,
						BirthPlace:   p.PlayerDetails.Birth.Place,
						BirthCountry: p.PlayerDetails.Birth.Country,
						BirthDate:    p.PlayerDetails.Birth.Date,
					},
				)
				//error?
				if err != nil {
					//sth more serious
					if !strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
						failedP = append(failedP, p.PlayerDetails.ID)
						pm.Logger.Err(err).Msg(err.Error())
					} else {
						pm.Logger.Debug().Msg(err.Error())
					}
					//all good
				} else {
					pm.Logger.Debug().Msg(fmt.Sprintf("inserted player_%d#row_%d ", p.PlayerDetails.ID, row))
				}
				//catch empty string ratin
				rating, err := strconv.ParseFloat(s.Games.Rating, 32)
				if err != nil {
					rating = 0
					pm.Logger.Debug().Msg(err.Error())
				}
				//insert season stats
				_, err = pm.Repo.InsertSeasonStats(&PlayerSeasonStatsRow{
					PlayerID:           p.PlayerDetails.ID,
					Season:             season,
					TeamID:             s.Team.ID,
					Minutes:            s.Games.Minutes,
					Position:           s.Games.Position,
					Rating:             rating,
					Captain:            s.Games.Captain,
					Appearances:        s.Games.Appearances,
					Lineups:            s.Games.Lineups,
					TotalShots:         s.Shots.Total,
					ShotsOnTarget:      s.Shots.On,
					TotalGoals:         s.Goals.Total,
					Assists:            s.Goals.Assists,
					TotalPasses:        s.Passes.Total,
					KeyPasses:          s.Passes.Key,
					PassAccuracy:       s.Passes.Accuracy,
					TotalTackles:       s.Tackles.Total,
					TackleBlocks:       s.Tackles.Blocks,
					Interceptions:      s.Tackles.Interceptions,
					TotalDuels:         s.Duels.Total,
					DuelsWon:           s.Duels.Won,
					DribbleAttempts:    s.Dribbles.Attempts,
					DribbleSuccess:     s.Dribbles.Success,
					YellowCards:        s.Cards.Yellow,
					RedCards:           s.Cards.Red,
					PenaltiesWon:       s.Penalty.Won,
					PenaltiesCommitted: s.Penalty.Committed,
					PenaltiesScored:    s.Penalty.Scored,
					PenaltiesMissed:    s.Penalty.Missed,
					PenaltiesSaved:     s.Penalty.Saved,
					GoalkeeperSaves:    s.Goals.Saves,
				})
				//process any error but duplicate
				if err != nil {
					if !strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
						pm.Logger.Err(err).Msg(fmt.Sprintf("failed inserting player season statistics %d", p.PlayerDetails.ID))
						failedS = append(failedS, p.PlayerDetails.ID)
					}
				} else {
					pm.Logger.Debug().Msg(fmt.Sprintf("inserted player season statistics %d", p.PlayerDetails.ID))
				}
			}
		}
		pgTotal = pg.Total
		pgCurrent = pg.Current + 1

	}
	return &failedP, &failedS, nil
}

/* func writeToCSV(ids []int, filename string, env.Logger zerolog.Logger) {
	// Create a new CSV file
	file, err := os.Create(filename)
	if err != nil {
		env.Logger.Err(err).Msg("")
	}
	defer file.Close()

	// Create a new CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Convert integers to strings and write to CSV file
	var strArray []string
	for _, num := range ids {
		strArray = append(strArray, strconv.Itoa(num))
	}

	err = writer.Write(strArray)
	if err != nil {
		env.Logger.Err(err).Msg("")
	}

	log.Println("CSV file created successfully")
}
*/
