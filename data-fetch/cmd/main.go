package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bernhardson/prefoot/data-fetch/pkg/database"
	"github.com/bernhardson/prefoot/data-fetch/pkg/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type application struct {
	logger *zerolog.Logger
	repo   *database.Repository
}

func main() {

	connConfig, err := pgxpool.ParseConfig("postgres://peterson:123@localhost/prefoot")
	if err != nil {
		log.Err(err).Msg("")
	}

	pool, err := pgxpool.New(context.Background(), connConfig.ConnString())
	if err != nil {
		log.Err(err).Msg("")
	}

	defer pool.Close()

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
		Level(zerolog.DebugLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	repo := database.Repository{
		Teams:    &database.TeamModel{Pool: pool},
		Venues:   &database.VenueModel{Pool: pool},
		Players:  &database.PlayerModel{Pool: pool},
		Fixture:  &database.FixtureModel{Pool: pool},
		League:   &database.LeagueModel{Pool: pool},
		Coach:    &database.CoachModel{Pool: pool},
		Standing: &database.StandingModel{Pool: pool},
		Logger:   logger,
	}

	app := &application{
		logger: &logger,
		repo:   &repo,
	}

	season := 2023
	league := 78
	fs, err := service.FetchAndInsertLeagues(&repo)
	logger.Err(err).Msg(fmt.Sprintf("insert leagues: failed=%v", *fs))

	ts, err := service.FetchAndInsertTeams(&repo, league, season)
	logger.Err(err).Msg(fmt.Sprintf("insert teams: league=%d#season=%d", league, season))

	fp, fs, err := service.FetchAndInsertPlayers(&repo, league, season)
	logger.Err(err).Msg(fmt.Sprintf("insert players: league:%d#season=%d", league, season))
	logger.Err(err).Msg(fmt.Sprintf("insert players: failedP=%v # failedS=%v", *fp, *fs))

	err = service.FetchAndInsertFixtures(&repo, league, season)
	logger.Err(err).Msg(fmt.Sprintf("insert fixtures: league:%d#season=%d", league, season))

	fc, fcc, err := service.FetchAndInsertCoaches(&repo, ts)
	logger.Err(err).Msg(fmt.Sprintf("insert coaches: league:%d#season=%d#failedC=%v#failedCC=%v", league, season, *fc, *fcc))

	addr := "localhost:8080"
	srv := &http.Server{
		Addr:    addr,
		Handler: app.routes(),
	}
	srv.ListenAndServe()

}
