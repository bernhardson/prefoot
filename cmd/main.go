package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/bernhardson/prefoot/internal/models"
	"github.com/bernhardson/prefoot/pkg/coach"
	"github.com/bernhardson/prefoot/pkg/fixture"
	"github.com/bernhardson/prefoot/pkg/leagues"
	"github.com/bernhardson/prefoot/pkg/players"
	"github.com/bernhardson/prefoot/pkg/result"
	"github.com/bernhardson/prefoot/pkg/rounds"
	"github.com/bernhardson/prefoot/pkg/team"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type application struct {
	logger         *zerolog.Logger
	player         *players.PlayerModel
	fixture        *fixture.FixtureModel
	league         *leagues.LeaguesModel
	team           *team.TeamModel
	coach          *coach.CoachModel
	sessionManager *scs.SessionManager
	users          *models.UserModel
}

func main() {

	//Initialize postgres connection
	connConfig, err := pgxpool.ParseConfig("postgres://peterson:123@localhost/prefoot")
	if err != nil {
		log.Err(err).Msg("")
	}

	pool, err := pgxpool.New(context.Background(), connConfig.ConnString())
	if err != nil {
		log.Err(err).Msg("")
	}

	defer pool.Close()

	// Initialize a new session manager and configure it to use pgxstore as the session store.
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(pool)
	sessionManager.Lifetime = 12 * time.Hour

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
		Level(zerolog.DebugLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	playerRepo := &players.Repo{
		Pool: pool,
	}

	app := &application{
		logger:         &logger,
		sessionManager: sessionManager,
		player: &players.PlayerModel{
			Logger: &logger,
			Repo:   playerRepo,
		},
		fixture: &fixture.FixtureModel{
			Logger:     &logger,
			PlayerRepo: playerRepo,
			Repo: &fixture.FixtureRepo{
				Pool: pool,
			},
			RoundRepo: &rounds.Repo{
				Pool: pool,
			},

			ResultRepo: &result.ResultRepo{
				Pool: pool,
			},
		},
		league: &leagues.LeaguesModel{
			Logger: &logger,
			Repo: &leagues.LeagueRepo{
				Pool: pool,
			},
		},
		team: &team.TeamModel{
			Logger: &logger,
			TeamRepo: &team.TeamRepository{
				Pool: pool,
			},
		},
		coach: &coach.CoachModel{
			Logger: &logger,
			Repo: &coach.CoachRepo{
				Pool: pool,
			},
		},
		users: &models.UserModel{
			Pool: pool,
		},
	}

	addr := "localhost:8080"
	tlsConfig := &tls.Config{CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256}}
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = srv.ListenAndServeTLS("/Users/peterson/git/prefoot/tls/cert.pem", "/Users/peterson/git/prefoot/tls/key.pem")
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

}
