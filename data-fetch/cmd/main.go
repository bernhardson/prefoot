package main

import (
	"github.com/bernhardson/prefoot-init/pkg"
	"github.com/jackc/pgx"
	"github.com/rs/zerolog/log"
)

func main() {
	config := pgx.ConnConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "prefoot",
		User:     "ulf",
		Password: "123",
	}
	// Connect to the PostgreSQL database
	conn, err := pgx.Connect(config)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	season := 2023
	league := 78

	pkg.InsertLeagues(conn)

	teams, err := pkg.GetTeams(league, season)
	if err != nil {
		log.Err(err).Msg("")
	}
	pkg.InsertTeams(teams, conn)

	players, err := pkg.GetPlayers(league, season)
	if err != nil {
		log.Err(err).Msg("")
	}
	pkg.InsertPlayers(players, conn)

	pkg.InsertCoaches(teams, conn)

	matches, err := pkg.GetMatches(league, season)
	if err != nil {
		log.Err(err).Msg("")
	}

	pkg.InsertFixtures(matches, conn, season)

}
