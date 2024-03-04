package rest

import (
	"net/http"
	"strconv"

	"github.com/bernhardson/prefoot/data-fetch/pkg/model"
	"github.com/bernhardson/prefoot/server/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var pool *pgxpool.Pool

type Data struct {
	PS []*model.Player `json:"player_statistics"`
	PD *[]model.Player `json:"player_details"`
}

type Response struct {
	Data    Data   `json:"data"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

var Pool *pgxpool.Pool

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
	resp, err := service.GetPlayersByTeamId(id, Pool)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	c.IndentedJSON(http.StatusOK, resp)

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
	ss, err := service.GetStandings(l, s)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	c.IndentedJSON(http.StatusOK, ss)
} 
