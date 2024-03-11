package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	insertTeam = `INSERT INTO teams (id, name, country)VALUES ($1, $2, $3)`
)

type TeamModel struct {
	Pool *pgxpool.Pool
}

type TeamRow struct {
	Id      int
	Name    string
	Country string
}

func (tm *TeamModel) Insert(t *TeamRow) (int64, error) {

	row, err := tm.Pool.Exec(
		context.Background(),
		insertTeam,
		t.Id,
		t.Name,
		t.Country,
	)
	return row.RowsAffected(), err
}
