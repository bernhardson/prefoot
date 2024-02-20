package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"

	"github.com/bernhardson/prefoot/data-fetch/pkg/comm"
	"github.com/bernhardson/prefoot/data-fetch/pkg/database"
	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
)

const (
	fixturesURL = "https://api-football-v1.p.rapidapi.com/v3/fixtures?league=%d&season=%d"

	fixtureDetailURL = "https://api-football-v1.p.rapidapi.com/v3/fixtures?id=%d"
)

func FetchAndInsert(fixtures *[]model.Match, pool *pgxpool.Pool, season int) error {
	for _, f := range *fixtures {
		//insert league
		fr, err := GetFixtureDetail(f.Fixture.ID)
		if err != nil {
			return err
		}
		for _, fd := range fr.FixtureDetail {
			//insert fixture
			database.InsertFixture(pool, f, fd, season)
			if err != nil {
				log.Err(err).Msg(fmt.Sprintf("fixture_%d", f.Fixture.ID))
			}
			for i, l := range fd.Lineups {
				ts := getTeamStatistics(i, &fd)
				err = database.InsertTeams(pool, &l, &f.Fixture, ts)
				if err != nil {
					log.Err(err).Msg(fmt.Sprintf("teamstatistics#team_%d#fixture_%d", l.Team.ID, f.Fixture.ID))
				}
				err := database.InsertFormation(pool, &l, &f.Fixture, ts)
				if err != nil {
					log.Err(err).Msg(fmt.Sprintf("fixture_%d", f.Fixture.ID))
				}
			}
			for _, teamStats := range fd.Players {
				// insert player statistics
				for _, player := range teamStats.Players {
					ps := player.Statistics[0]
					if ps.Games.Rating == "" {
						ps.Games.Rating = "0"
					}
					if ps.Passes.Accuracy == "" {
						ps.Passes.Accuracy = "0"
					}
					database.InsertPlayerStatistic(pool, &player, &ps, &f.Fixture, season)
					if err != nil {
						log.Err(err).Msg(fmt.Sprintf("playerstatistics#player_%d#fixture_%d", player.Player.ID, f.Fixture.ID))
					}
				}

			}
		}
	}
	return nil
}

func GetMatches(league int, season int) (*[]model.Match, error) {
	data, err := comm.GetHttpBody(fixturesURL, league, season)
	if err != nil {
		return nil, err
	}

	matches := &[]model.Match{}
	err = mapstructure.Decode(data["response"], matches)
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func GetFixtureDetail(id int) (*model.FixtureResponse, error) {
	data, err := comm.GetHttpBody(fixtureDetailURL, id)
	if err != nil {
		return nil, err
	}

	fd := &model.FixtureResponse{}
	err = mapstructure.Decode(data, fd)
	if err != nil {
		return nil, err
	}
	return fd, nil

}

// convert team statistics since they come as mixed type json
func getTeamStatistics(home int, fd *model.FixtureDetail) *model.TeamStatistics {
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
				log.Err(err).Msg("")
				sd.ExpectedGoals = 0
			} else {
				sd.ExpectedGoals = fl
			}
		}
	}
	return &sd
}
