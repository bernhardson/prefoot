package team

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	insertVenue = `INSERT INTO venues (id, name, city) VALUES ($1, $2, $3)`
)

type VenueModel struct {
	Pool *pgxpool.Pool
}

type VenueRow struct {
	Id   int
	Name string
	City string
}

func (tm *VenueModel) Insert(v *VenueRow) (int64, error) {
	//insert venue
	row, err := tm.Pool.Exec(
		context.Background(),
		insertVenue,
		v.Id,
		v.Name,
		v.City,
	)

	return row.RowsAffected(), err
}
