package coach

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	insertCoach       = `INSERT INTO coaches (id, name) VALUES ($1, $2)`
	insertCoachCareer = `INSERT INTO coach_careers (coach, team, start, "end") VALUES ($1, $2, $3, $4)`
)

type CoachRepo struct {
	Pool *pgxpool.Pool
}

// Coach represents the coaches table
type CoachRow struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (cm *CoachRepo) Insert(c *CoachRow) (int64, error) {
	row, err := cm.Pool.Exec(
		context.Background(),
		insertCoach,
		c.ID, c.Name,
	)
	return row.RowsAffected(), err
}

// CoachCareer represents the coach_careers table
type CoachCareerRow struct {
	CoachID int        `json:"coach_id"`
	TeamID  int        `json:"team_id"`
	Start   *time.Time `json:"start"`
	End     *time.Time `json:"end"`
}

func (cm *CoachRepo) InsertCareer(c *CoachCareerRow) (int64, error) {
	row, err := cm.Pool.Exec(
		context.Background(),
		insertCoachCareer,
		c.CoachID, c.TeamID, c.Start, c.End,
	)

	return row.RowsAffected(), err
}
