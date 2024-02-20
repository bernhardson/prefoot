package services

import (
	"github.com/bernhardson/prefoot/server/pkg/service"
	"github.com/bernhardson/prefoot/server/pkg/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPlayersByTeamId(teamId int, pool *pgxpool.Pool) ([]*database.Player, *[]pkg.Player, error) {
	ids, ps, err := GetPlayerIdsByTeamId(teamId)
	if err != nil {
		return nil, nil, err
	}
	psdb, err := database.SelectPlayerStatistics(ids, pool)

	if err != nil {
		return nil, nil, err
	}
	return psdb, ps, nil
}
