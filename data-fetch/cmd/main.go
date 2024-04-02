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
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	repo := database.Repository{
		Teams:   &database.TeamModel{Pool: pool},
		Venues:  &database.VenueModel{Pool: pool},
		Players: &database.PlayerModel{Pool: pool},
		Fixture: &database.FixtureModel{Pool: pool},
		League:  &database.LeagueModel{Pool: pool},
		Coach:   &database.CoachModel{Pool: pool},
		Result:  &database.ResultModel{Pool: pool},
		Logger:  logger,
	}

	app := &application{
		logger: &logger,
		repo:   &repo,
	}

	init := false

	update := false

	if update {

	}

	if init {
		// leagues := []int{2, 39, 45, 48, 3, 142, 78, 79, 81, 143, 71, 137}
		leagues := []int{71, 137}
		go initDB(leagues, &logger, &repo)
	}
	addr := "localhost:8080"
	srv := &http.Server{
		Addr:    addr,
		Handler: app.routes(),
	}
	srv.ListenAndServe()

}

func initDB(leagues []int, logger *zerolog.Logger, repo *database.Repository) {
	for _, l := range leagues {
		lresp, fs, err := service.FetchAndInsertLeague(repo, l)
		logger.Err(err).Msg(fmt.Sprintf("insert leagues: failed=%v", *fs))

		for _, s := range lresp.Response[0].Seasons {
			// for _, year := range seasons {

			year := s.Year

			logger.Info().Msg(fmt.Sprintf("inserting league=%d#season=%d", l, year))
			ts, err := service.FetchAndInsertTeams(repo, l, year)
			logger.Err(err).Msg(fmt.Sprintf("insert teams: league=%d#season=%d", l, year))

			fp, fs, err := service.FetchAndInsertPlayers(repo, l, year)
			logger.Err(err).Msg(fmt.Sprintf("insert players: league:%d#season=%d", l, year))
			logger.Err(err).Msg(fmt.Sprintf("insert players: failedP=%v # failedS=%v", *fp, *fs))

			err = service.FetchAndInsertFixtures(repo, l, year)
			logger.Err(err).Msg(fmt.Sprintf("insert fixtures: league:%d#season=%d", l, year))

			fc, fcc, err := service.FetchAndInsertCoaches(repo, ts)
			logger.Err(err).Msg(fmt.Sprintf("insert coaches: league:%d#season=%d#failedC=%v#failedCC=%v", l, year, *fc, *fcc))
		}
	}
}
