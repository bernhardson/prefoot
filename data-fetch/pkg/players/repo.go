package players

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	insertPlayerStatistics                   = `INSERT INTO player_statistics (player, fixture, team, league, season, minutes, position, rating, captain, substitute, shots_total, shots_on, goals_scored, goals_assisted, passes_total, passes_key, accuracy, tackles, block, interceptions, duels_total, duels_won, dribbles_total,dribbles_won, yellow, red, penalty_won, penalty_committed, penalty_scored, penalty_missed, penalty_saved, saves)	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,$18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32)`
	insertPlayer                             = `INSERT INTO players (id, team, season, firstname, lastname, birthplace, birthcountry, birthdate) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	insertPlayerStatisticsSeason             = `INSERT INTO player_statistics_season ("player","season", "team", "minutes", "position", "rating","captain", "games", "lineups", "shots_total", "shots_on", "goals_scored", "goals_assisted", "passes_total", "passes_key","accuracy", "tackles", "block", "interceptions", "duels_total", "duels_won", "dribbles_total", "dribbles_won", "yellow", "red", "penalty_won","penalty_committed", "penalty_scored", "penalty_missed", "penalty_saved", "saves") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31)`
	selectPlayer                             = `SELECT * FROM players WHERE id=$1`
	selectPlayerStats                        = `SELECT p.id AS player_id, p.team, p.season, p.firstname, p.lastname, p.birthplace, p.birthcountry, p.birthdate, ps.fixture, ps.minutes, ps.position, ps.rating, ps.captain, ps.substitute, ps.shots_total, ps.shots_on, ps.goals_scored, ps.goals_assisted, ps.passes_total, ps.passes_key, ps.accuracy, ps.tackles, ps.block, ps.interceptions, ps.duels_total, ps.duels_won, ps.dribbles_total, ps.dribbles_won, ps.yellow, ps.red, ps.penalty_won, ps.penalty_committed, ps.penalty_scored, ps.penalty_missed, ps.penalty_saved, ps.saves FROM players p JOIN player_statistics ps ON p.team = ps.team WHERE p.team = $1`
	selectPlayersByTeam                      = `SELECT id, team, season, firstname, lastname, birthplace, birthcountry, birthdate FROM players WHERE team = $1;`
	selectPlayerIdsByTeam                    = `SELECT id FROM players WHERE season = $1 AND team = $2;`
	selectPlayersByTeamLeagueSeason          = `SELECT id, team, season, firstname, lastname, birthplace, birthcountry, birthdate FROM players WHERE season=$1 AND team = $2;`
	selectPlayerStatistics                   = `SELECT * FROM player_statistics WHERE player= ANY($1) AND fixture=ANY($2) AND team=$3`
	selectKeyPlayerStatsByFixturesAndPlayers = `SELECT p.id AS player_id, p.firstname, p.lastname, ps.team,SUM(ps.goals_scored) AS total_goals_scored, SUM(ps.goals_assisted) AS total_goals_assisted, AVG(ps.duels_total) AS avg_duels_total, AVG(ps.duels_won) AS avg_duels_won, AVG(ps.passes_key) AS avg_key_passes, AVG(ps.rating) AS avg_rating FROM players p JOIN player_statistics ps ON p.id = ps.player where ps.player= ANY($1) and ps.fixture = ANY($2) GROUP BY ps.team, p.id, p.firstname, p.lastname;`
)

type PlayerRow struct {
	Id           int    `json:"playerID"`
	Team         int    `json:"team"`
	Season       int    `json:"season"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	BirthPlace   string `json:"birthPlace"`
	BirthCountry string `json:"birthCountry"`
	BirthDate    string `json:"birthDate"`
}

type Repo struct {
	Pool *pgxpool.Pool
}

func (pm *Repo) Insert(p *PlayerRow) (int64, error) {
	row, err := pm.Pool.Exec(
		context.Background(),
		insertPlayer,
		p.Id, p.Team, p.Season, p.FirstName,
		p.LastName, p.BirthPlace,
		p.BirthCountry, p.BirthDate)
	return row.RowsAffected(), err
}

func (pm *Repo) Select(id int) (*PlayerRow, error) {

	p := &PlayerRow{}
	err := pm.Pool.QueryRow(context.Background(), selectPlayer, id).Scan(&p.Id, &p.Team, &p.Season, &p.FirstName, &p.LastName, &p.BirthPlace, &p.BirthCountry, &p.BirthDate)
	return p, err
}

func (pm *Repo) SelectPlayersByTeamId(id int) ([]*PlayerRow, error) {

	rows, err := pm.Pool.Query(context.Background(), "SELECT * FROM players WHERE team = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[PlayerRow])
	if err != nil {
		return nil, err
	}
	return players, nil
}

func (pm *Repo) SelectPlayerIdsBySeasonAndTeamId(season, team int) ([]int, error) {

	rows, err := pm.Pool.Query(context.Background(), selectPlayerIdsByTeam, season, team)
	if err != nil {
		return nil, err
	}

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

type PlayerStatsRow struct {
	Player           int     `json:"player"`
	Fixture          int     `json:"fixture"`
	Team             int     `json:"team"`
	League           int     `json:"league"`
	Season           int     `json:"season"`
	Minutes          int     `json:"minutes"`
	Position         string  `json:"position"`
	Rating           float64 `json:"rating"`
	Captain          bool    `json:"captain"`
	Substitute       bool    `json:"substitute"`
	ShotsTotal       int     `json:"shots_total"`
	ShotsOn          int     `json:"shots_on"`
	GoalsScored      int     `json:"goals_scored"`
	GoalsAssisted    int     `json:"goals_assisted"`
	PassesTotal      int     `json:"passes_total"`
	PassesKey        int     `json:"passes_key"`
	Accuracy         int     `json:"accuracy"`
	Tackles          int     `json:"tackles"`
	Block            int     `json:"block"`
	Interceptions    int     `json:"interceptions"`
	DuelsTotal       int     `json:"duels_total"`
	DuelsWon         int     `json:"duels_won"`
	DribblesTotal    int     `json:"dribbles_total"`
	DribblesWon      int     `json:"dribbles_won"`
	Yellow           int     `json:"yellow"`
	Red              int     `json:"red"`
	PenaltyWon       int     `json:"penalty_won"`
	PenaltyCommitted int     `json:"penalty_committed"`
	PenaltyScored    int     `json:"penalty_scored"`
	PenaltyMissed    int     `json:"penalty_missed"`
	PenaltySaved     int     `json:"penalty_saved"`
	Saves            int     `json:"saves"`
}

func (pm *Repo) InsertStats(ps *PlayerStatsRow) (int64, error) {

	row, err := pm.Pool.Exec(
		context.Background(),
		insertPlayerStatistics,
		ps.Player, ps.Fixture, ps.Team, ps.Season, ps.Minutes, ps.Position, ps.Rating,
		ps.Captain, ps.Substitute, ps.ShotsTotal, ps.ShotsOn, ps.GoalsScored,
		ps.GoalsAssisted, ps.PassesTotal, ps.PassesKey, ps.Accuracy, ps.Tackles,
		ps.Block, ps.Interceptions, ps.DuelsTotal, ps.DuelsWon,
		ps.DribblesTotal, ps.DribblesWon, ps.Yellow, ps.Red, ps.PenaltyWon,
		ps.PenaltyCommitted, ps.PenaltyScored, ps.PenaltyMissed, ps.PenaltySaved, ps.Saves,
	)
	return row.RowsAffected(), err
}

type PlayerSeasonStatsRow struct {
	PlayerID           int     `json:"playerID"`
	Season             int     `json:"season"`
	TeamID             int     `json:"teamID"`
	Minutes            int     `json:"minutes"`
	Position           string  `json:"position"`
	Rating             float64 `json:"rating"`
	Captain            bool    `json:"captain"`
	Appearances        int     `json:"appearances"`
	Lineups            int     `json:"lineups"`
	TotalShots         int     `json:"totalShots"`
	ShotsOnTarget      int     `json:"shotsOnTarget"`
	TotalGoals         int     `json:"totalGoals"`
	Assists            int     `json:"assists"`
	TotalPasses        int     `json:"totalPasses"`
	KeyPasses          int     `json:"keyPasses"`
	PassAccuracy       int     `json:"passAccuracy"`
	TotalTackles       int     `json:"totalTackles"`
	TackleBlocks       int     `json:"tackleBlocks"`
	Interceptions      int     `json:"interceptions"`
	TotalDuels         int     `json:"totalDuels"`
	DuelsWon           int     `json:"duelsWon"`
	DribbleAttempts    int     `json:"dribbleAttempts"`
	DribbleSuccess     int     `json:"dribbleSuccess"`
	YellowCards        int     `json:"yellowCards"`
	RedCards           int     `json:"redCards"`
	PenaltiesWon       int     `json:"penaltiesWon"`
	PenaltiesCommitted int     `json:"penaltiesCommitted"`
	PenaltiesScored    int     `json:"penaltiesScored"`
	PenaltiesMissed    int     `json:"penaltiesMissed"`
	PenaltiesSaved     int     `json:"penaltiesSaved"`
	GoalkeeperSaves    int     `json:"goalkeeperSaves"`
}

func (pm *Repo) InsertSeasonStats(s *PlayerSeasonStatsRow) (int64, error) {

	row, err := pm.Pool.Exec(
		context.Background(),
		insertPlayerStatisticsSeason,
		s.PlayerID, s.Season, s.TeamID, s.Minutes,
		s.Position, s.Rating, s.Captain, s.Appearances,
		s.Lineups, s.TotalShots, s.ShotsOnTarget, s.TotalGoals,
		s.Assists, s.TotalPasses, s.KeyPasses, s.PassAccuracy,
		s.TotalTackles, s.TackleBlocks, s.Interceptions, s.TotalDuels,
		s.DuelsWon, s.DribbleAttempts, s.DribbleSuccess, s.YellowCards, s.RedCards,
		s.PenaltiesWon, s.PenaltiesCommitted, s.PenaltiesScored, s.PenaltiesMissed,
		s.PenaltiesSaved, s.GoalkeeperSaves,
	)
	return row.RowsAffected(), err
}

type PlayersJoinOnPlayerStatsRow struct {
	PlayerID         int
	Team             int
	Season           int
	FirstName        string
	LastName         string
	BirthPlace       string
	BirthCountry     string
	BirthDate        string
	Fixture          int
	Minutes          int
	Position         string
	Rating           float64
	Captain          bool
	Substitute       bool
	ShotsTotal       int
	ShotsOn          int
	GoalsScored      int
	GoalsAssisted    int
	PassesTotal      int
	PassesKey        int
	Accuracy         int
	Tackles          int
	Block            int
	Interceptions    int
	DuelsTotal       int
	DuelsWon         int
	DribblesTotal    int
	DribblesWon      int
	Yellow           int
	Red              int
	PenaltyWon       int
	PenaltyCommitted int
	PenaltyScored    int
	PenaltyMissed    int
	PenaltySaved     int
	Saves            int
}

func (pm *Repo) SelectPlayersAndStatisticsByTeamId(id int) (*[]*PlayersJoinOnPlayerStatsRow, error) {

	rows, err := pm.Pool.Query(
		context.Background(),
		selectPlayerStats, id)
	if err != nil {
		return nil, err
	}

	pls, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[PlayersJoinOnPlayerStatsRow])
	if err != nil {
		return nil, err
	}

	return &pls, err
}

func (pm *Repo) SelectPlayersByTeamLeagueSeason(season, team int) ([]*PlayerRow, error) {

	// Execute the query
	rows, err := pm.Pool.Query(context.Background(), selectPlayersByTeamLeagueSeason, season, team)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set
	var players []*PlayerRow
	players, err = pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[PlayerRow])
	if err != nil {
		return nil, err
	}
	return players, nil
}

type KeyPlayerStats struct {
	PlayerID           int     `json:"player_id"`
	FirstName          string  `json:"firstname"`
	LastName           string  `json:"lastname"`
	Team               int     `json:"team"`
	TotalGoalsScored   int     `json:"total_goals_scored"`
	TotalGoalsAssisted int     `json:"total_goals_assisted"`
	AvgDuelsTotal      float64 `json:"avg_duels_total"`
	AvgDuelsWon        float64 `json:"avg_duels_won"`
	AvgKeyPasses       float64 `json:"avg_key_passes"`
	AvgRating          float64 `json:"avg_rating"`
}

func (pm *Repo) SelectPlayerStatisticsByPlayersFixturesTeam(playerIds []int, fixtureIds *[]int) ([]*KeyPlayerStats, error) {

	rows, err := pm.Pool.Query(context.Background(), selectKeyPlayerStatsByFixturesAndPlayers, playerIds, fixtureIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []*KeyPlayerStats
	stats, err = pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[KeyPlayerStats])
	if err != nil {
		return nil, err
	}
	return stats, nil
}
