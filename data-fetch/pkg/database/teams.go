package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	insertTeam = `INSERT INTO teams (id, name, country, code)VALUES ($1, $2, $3, $4)`
	selectTeam = `SELECT * FROM teams WHERE id=$1`
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
