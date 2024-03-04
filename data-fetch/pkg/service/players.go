package service

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"github.com/bernhardson/prefoot/data-fetch/pkg/database"
	"github.com/bernhardson/prefoot/data-fetch/pkg/fetch"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

func FetchAndInsertPlayers(pool *pgxpool.Pool, league int, season int, logger zerolog.Logger) error {

	ps, pg, err := fetch.GetPlayers(league, season, 1)

	if err != nil {
		return err
	}
	failedP, failedS := database.InsertPlayers(ps, pool, logger)

	writeToCSV(failedP, "player.csv", logger)
	writeToCSV(failedS, "player_statistics.csv", logger)

	for i := 2; i < pg.Total; i++ {
		ps, pg, err = fetch.GetPlayers(league, season, i)

		if err != nil {
			return err
		}
		database.InsertPlayers(ps, pool, logger)
	}
	return nil
}

func writeToCSV(ids []int, filename string, logger zerolog.Logger) {
	// Create a new CSV file
	file, err := os.Create(filename)
	if err != nil {
		logger.Err(err).Msg("")
	}
	defer file.Close()

	// Create a new CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Convert integers to strings and write to CSV file
	var strArray []string
	for _, num := range ids {
		strArray = append(strArray, strconv.Itoa(num))
	}

	err = writer.Write(strArray)
	if err != nil {
		logger.Err(err).Msg("")
	}

	log.Println("CSV file created successfully")
}
