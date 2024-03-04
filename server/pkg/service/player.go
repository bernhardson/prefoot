package service

import (
	"github.com/bernhardson/prefoot/server/pkg/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPlayersByTeamId(teamId int, pool *pgxpool.Pool) (*[]*database.Player, error) {

	res, err := database.SelectPlayersAndStatisticsByTeamId(teamId, pool)

	if err != nil {
		return nil, err
	}
	return res, nil
}
