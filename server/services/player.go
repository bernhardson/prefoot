package services

import (
	pkg "github.com/bernhardson/prefoot-init/pkg"
	db "github.com/bernhardson/prefoot-server/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPlayersByTeamId(teamId int, pool *pgxpool.Pool) ([]*db.Player, *[]pkg.Player, error) {
	ids, ps, err := pkg.GetPlayerIdsByTeamId(teamId)
	if err != nil {
		return nil, nil, err
	}
	psdb, err := db.SelectPlayerStatistics(ids, pool)

	if err != nil {
		return nil, nil, err
	}
	return psdb, ps, nil
}
