package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bernhardson/prefoot/data-fetch/pkg/coach"
	"github.com/bernhardson/prefoot/data-fetch/pkg/fixture"
	"github.com/bernhardson/prefoot/data-fetch/pkg/leagues"
	"github.com/bernhardson/prefoot/data-fetch/pkg/players"
	"github.com/bernhardson/prefoot/data-fetch/pkg/team"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type application struct {
	logger  *zerolog.Logger
	player  *players.PlayerModel
	fixture *fixture.FixtureModel
	league  *leagues.LeaguesModel
	team    *team.TeamModel
	coach   *coach.CoachModel
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

	app := &application{
		logger: &logger,
		player: &players.PlayerModel{Repo: &players.Repo{Pool: pool}},
		fixture: &fixture.FixtureModel{Logger: &logger, Repo: &fixture.FixtureRepo{
			Pool: pool,
		}},
	}

	update := true
	if update {
		logger.Err(app.fixture.UpdateFixture(2, 2023)).Msg("")
		// _, _, err = service.FetchAndInsertPlayers(app.repo, 78, 2018)
		// if err != nil {
		// 	logger.Err(err).Msg("")
		// }
	}

	init := false
	if init {
		leagues := []int{71, 137}
		// leagues := []int{78}
		go app.initDB(leagues)
	}

	addr := "localhost:8080"
	srv := &http.Server{
		Addr:    addr,
		Handler: app.routes(),
	}
	srv.ListenAndServe()
}

func (app *application) initDB(leagues []int) {
	for _, l := range leagues {
		lresp, fs, err := app.league.FetchAndInsertLeague(l)
		app.logger.Err(err).Msg(fmt.Sprintf("insert leagues: failed=%v", *fs))

		for _, s := range lresp.Response[0].Seasons {
			year := s.Year

			app.logger.Info().Msg(fmt.Sprintf("inserting league=%d#season=%d", l, year))
			ts, err := app.team.FetchAndInsertTeams(l, year)
			app.logger.Err(err).Msg(fmt.Sprintf("insert teams: league=%d#season=%d", l, year))

			fp, fs, err := app.player.FetchAndInsertPlayers(l, year)
			app.logger.Err(err).Msg(fmt.Sprintf("insert players: league:%d#season=%d", l, year))
			app.logger.Err(err).Msg(fmt.Sprintf("insert players: failedP=%v # failedS=%v", *fp, *fs))

			err = app.fixture.FetchAndInsertFixtures(l, year)
			app.logger.Err(err).Msg(fmt.Sprintf("insert fixtures: league:%d#season=%d", l, year))

			fc, fcc, err := app.coach.FetchAndInsertCoaches(ts)
			app.logger.Err(err).Msg(fmt.Sprintf("insert coaches: league:%d#season=%d#failedC=%v#failedCC=%v", l, year, *fc, *fcc))
		}
	}
}
