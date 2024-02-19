package rest

import (
	"net/http"
	"strconv"

	"github.com/bernhardson/prefoot-init/pkg"
	"github.com/bernhardson/prefoot-server/internal/database"
	"github.com/bernhardson/prefoot-server/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var Pool *pgxpool.Pool

type Data struct {
	PS []*database.Player `json:"player_statistics"`
	PD *[]pkg.Player      `json:"player_details"`
}

type Response struct {
	Data    Data   `json:"data"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func StartServer() {
	router := gin.Default()
	router.GET("/players/:teamId", getPlayersByTeamId)
	router.GET("leagues/:league/:season", getLeagueStanding)
	router.Run("localhost:8080")
}

func getPlayersByTeamId(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("teamId"))
	if err != nil {
		log.Err(err).Msg("")
	}
	p, pl, err := services.GetPlayersByTeamId(id, Pool)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	response := Response{
		Data: Data{
			PS: p,
			PD: pl},
		Success: true,
		Message: "Data fetched successfully",
	}
	c.IndentedJSON(http.StatusOK, response)
}

func getLeagueStanding(c *gin.Context) {
	l, err := strconv.Atoi(c.Param("league"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	s, err := strconv.Atoi(c.Param("season"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	ss, err := services.GetStandings(l, s)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	c.IndentedJSON(http.StatusOK, ss)
}
