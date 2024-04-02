package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	insertTeam         = `INSERT INTO teams (id, name, country, code)VALUES ($1, $2, $3, $4)`
	selectTeam         = `SELECT * FROM teams WHERE id=$1`
	insertLeagueSeason = `INSERT INTO "seasons" ("league", "season", "team") VALUES ($1, $2, $3)`
	selectTeamsSeason  = `SELECT "team" FROM "seasons" WHERE "league"=$1 AND "season"=$2`
	selectTeamsByIds   = `SELECT * FROM "teams" WHERE id = ANY($1)`
)

type TeamModel struct {
	Pool *pgxpool.Pool
}

type TeamRow struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Code    string `json:"code"`
}

func (tm *TeamModel) Insert(t *TeamRow) (int64, error) {

	row, err := tm.Pool.Exec(
		context.Background(),
		insertTeam,
		t.Id,
		t.Name,
		t.Country,
		t.Code,
	)
	return row.RowsAffected(), err
}

func (tm *TeamModel) Select(id int) (*TeamRow, error) {

	t := &TeamRow{}
	err := tm.Pool.QueryRow(context.Background(), selectTeam, id).Scan(&t.Id, &t.Name, &t.Country, &t.Code)
	return t, err
}

// List teams of a league of one season
type TeamSeasonRow struct {
	League int `json:"league"`
	Season int `json:"season"`
	Team   int `json:"team"`
}

func (tm *TeamModel) InsertTeamSeason(r *TeamSeasonRow) (int64, error) {

	row, err := tm.Pool.Exec(
		context.Background(),
		insertLeagueSeason,
		r.League, r.Season, r.Team,
	)
	return row.RowsAffected(), err
}

type TeamIds struct {
	Team int
}

// Select all teams of that played in a certain league in a certain season.
func (tm *TeamModel) SelectTeamsSeason(league, season int) (*[]*TeamIds, error) {
	rows, err := tm.Pool.Query(context.Background(), selectTeamsSeason, league, season)

	if err != nil {
		return nil, err
	}

	ts, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[TeamIds])
	return &ts, err

}

// Select teams of a list of ids.
func (tm *TeamModel) SelectTeamsByIds(ids *[]int) (*[]*TeamRow, error) {
	rows, err := tm.Pool.Query(context.Background(), selectTeamsByIds, ids)

	if err != nil {
		return nil, err
	}

	ts, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[TeamRow])
	return &ts, err

}
