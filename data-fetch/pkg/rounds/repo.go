package rounds

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	insertRound = `INSERT INTO rounds ("league",  "season", "round", "start", "end") VALUES ($1, $2, $3, $4, $5) ON CONFLICT ("league", "season", "round") DO UPDATE SET "start" = EXCLUDED."start", "end" = EXCLUDED."end";`

	selectStartEndFromRounds  = `SELECT "start" FROM rounds WHERE league = $1 AND season = $2 AND round = $3`
	selectRoundsByTimestamp   = `SELECT "round" FROM rounds WHERE "league" = $1 AND season = $2 AND "start" > $3  AND "start" = (SELECT MIN("start") FROM rounds WHERE "league" = $1 AND season = $2 AND "start" > $3) LIMIT 1;`
	selectLatestFinishedRound = `SELECT "round" FROM rounds WHERE "league" = $1 AND season = $2 AND "end" <= $3 ORDER BY ABS("end" - $3) ASC LIMIT 1;`
)

type Repo struct {
	Pool *pgxpool.Pool
}

type RoundRow struct {
	Start  int64 `json:"start"`
	End    int64 `json:"end"`
	Round  int   `json:"round"`
	Season int   `json:"season"`
	League int   `json:"league"`
}

func (rm *Repo) Insert(f *RoundRow) (int64, error) {

	row, err := rm.Pool.Exec(
		context.Background(),
		insertRound,
		f.League, f.Season, f.Round, f.Start, f.End)

	return row.RowsAffected(), err
}

func (rm *Repo) SelectRoundByTimestamp(league, season int, timestamp int64) (*RoundRow, error) {

	row := &RoundRow{Start: timestamp, League: league, Season: season}

	err := rm.Pool.QueryRow(context.Background(), selectRoundsByTimestamp, league, season, timestamp).Scan(&row.Round)
	if err != nil {
		return nil, err
	}
	return row, nil
}

func (rm *Repo) SelectLatestFinishedRound(league, season int, timestamp int64) (*RoundRow, error) {

	row := &RoundRow{Start: timestamp, League: league, Season: season}

	err := rm.Pool.QueryRow(context.Background(), selectLatestFinishedRound, league, season, timestamp).Scan(&row.Round)
	if err != nil {
		return nil, err
	}
	return row, nil
}

func (pm *Repo) SelectTimestampFromRounds(league, season, round int) (int, error) {

	start := -1
	err := pm.Pool.QueryRow(context.Background(), selectStartEndFromRounds, league, season, round).Scan(&start)
	if err != nil {
		return start, err
	}
	return start, nil
}
