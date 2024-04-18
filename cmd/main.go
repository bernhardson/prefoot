package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/bernhardson/prefoot/pkg/coach"
	"github.com/bernhardson/prefoot/pkg/fixture"
	"github.com/bernhardson/prefoot/pkg/leagues"
	"github.com/bernhardson/prefoot/pkg/players"
	"github.com/bernhardson/prefoot/pkg/team"
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
		player: &players.PlayerModel{Logger: &logger,
			Repo: &players.Repo{Pool: pool}},
		fixture: &fixture.FixtureModel{
			Logger: &logger,
			Repo: &fixture.FixtureRepo{
				Pool: pool,
			},
		},
		league: &leagues.LeaguesModel{
			Logger: &logger,
			Repo: &leagues.LeagueRepo{
				Pool: pool,
			},
		},
	}

	addr := "localhost:8080"
	srv := &http.Server{
		Addr:    addr,
		Handler: app.routes(),
	}
	srv.ListenAndServe()
}
