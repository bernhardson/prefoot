package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	insertPlayerStatistics       = `INSERT INTO player_statistics (player, fixture, team, season, minutes, position, rating, captain, substitute, shots_total, shots_on, goals_scored, goals_assisted, passes_total, passes_key, accuracy, tackles, block, interceptions, duels_total, duels_won, dribbles_total,dribbles_won, yellow, red, penalty_won, penalty_committed, penalty_scored, penalty_missed, penalty_saved, saves)	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,$18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31)`
	insertPlayer                 = `INSERT INTO players (id, team, season, firstname, lastname, birthplace, birthcountry, birthdate) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	insertPlayerStatisticsSeason = `INSERT INTO player_statistics_season ("player","season", "team", "minutes", "position", "rating","captain", "games", "lineups", "shots_total", "shots_on", "goals_scored", "goals_assisted", "passes_total", "passes_key","accuracy", "tackles", "block", "interceptions", "duels_total", "duels_won", "dribbles_total", "dribbles_won", "yellow", "red", "penalty_won","penalty_committed", "penalty_scored", "penalty_missed", "penalty_saved", "saves") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31)`
)

type PlayerModel struct {
	Pool *pgxpool.Pool
}

type PlayerRow struct {
	Id           int    `json:"playerID"`
	TeamID       int    `json:"teamID"`
	Season       int    `json:"season"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	BirthPlace   string `json:"birthPlace"`
	BirthCountry string `json:"birthCountry"`
	BirthDate    string `json:"birthDate"`
}

func (pm *PlayerModel) Insert(p *PlayerRow) (int64, error) {
	row, err := pm.Pool.Exec(
		context.Background(),
		insertPlayer,
		p.Id, p.TeamID, p.Season, p.FirstName,
		p.LastName, p.BirthPlace,
		p.BirthCountry, p.BirthDate)
	return row.RowsAffected(), err
}

type PlayerStatsRow struct {
	PlayerID         int     `json:"player"`
	FixtureID        int     `json:"fixture"`
	TeamID           int     `json:"team"`
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

func (pm *PlayerModel) InsertStats(ps *PlayerStatsRow) (int64, error) {

	row, err := pm.Pool.Exec(
		context.Background(),
		insertPlayerStatistics,
		ps.PlayerID, ps.FixtureID, ps.TeamID, ps.Season, ps.Minutes, ps.Position, ps.Rating,
		ps.Captain, ps.Substitute, ps.ShotsTotal, ps.ShotsOn, ps.GoalsScored,
		ps.GoalsAssisted, ps.PassesTotal, ps.PassesKey, ps.Accuracy, ps.Tackles,
		ps.Block, ps.Interceptions, ps.DuelsTotal, ps.DuelsWon,
		ps.DribblesTotal, ps.DribblesWon, ps.Yellow, ps.Red, ps.PenaltyWon,
		ps.PenaltyCommitted, ps.PenaltyScored, ps.PenaltyMissed, ps.PenaltySaved, ps.Saves,
	)
	return row.RowsAffected(), err
}

type PlayerStatisticsRow struct {
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

func (pm *PlayerModel) InsertSeasonStats(s *PlayerStatisticsRow) (int64, error) {

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

func (pm *PlayerModel) SelectPlayersAndStatisticsByTeamId(id int) (*[]*PlayersJoinOnPlayerStatsRow, error) {

	// Perform a JOIN query to retrieve data from both tables
	rows, err := pm.Pool.Query(
		context.Background(),
		`SELECT
        p.id AS player_id,
        p.team,
        p.season,
        p.firstname,
        p.lastname,
        p.birthplace,
        p.birthcountry,
        p.birthdate,
        ps.fixture,
        ps.minutes,
        ps.position,
        ps.rating,
        ps.captain,
        ps.substitute,
        ps.shots_total,
        ps.shots_on,
        ps.goals_scored,
        ps.goals_assisted,
        ps.passes_total,
        ps.passes_key,
        ps.accuracy,
        ps.tackles,
        ps.block,
        ps.interceptions,
        ps.duels_total,
        ps.duels_won,
        ps.dribbles_total,
        ps.dribbles_won,
        ps.yellow,
        ps.red,
        ps.penalty_won,
        ps.penalty_committed,
        ps.penalty_scored,
        ps.penalty_missed,
        ps.penalty_saved,
        ps.saves
    FROM
        players p
    JOIN
        player_statistics ps ON p.team = ps.team
	WHERE p.team = $1;`, id)

	if err != nil {
		return nil, err
	}

	pls, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[PlayersJoinOnPlayerStatsRow])
	return &pls, err
}

type Player struct {
	Player           int
	ID               int
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

func (pm *PlayerModel) SelectPlayersByTeamId(id int) ([]*Player, error) {

	rows, err := pm.Pool.Query(context.Background(), "SELECT * FROM players WHERE team = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Player])
	if err != nil {
		return nil, err
	}
	return players, nil
}
