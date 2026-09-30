package repository

import (
	"context"
	"errors"

	"github.com/ErenKarakus1/Matchmaking-System/player-service/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUsernameAlreadyUsed = errors.New("username is already used")

const insertPlayerSQL = `
	INSERT INTO players (
		id,
		username
	)
	VALUES ($1,$2)
	RETURNING
		id,
		username,
		created_at,
		updated_at
`

func InsertPlayer(ctx context.Context, pool *pgxpool.Pool, player models.Player) (models.Player, error) {
	var createdPlayer models.Player
	err := pool.QueryRow(
		ctx,
		insertPlayerSQL,
		player.ID,
		player.Username,
	).Scan(
		&createdPlayer.ID,
		&createdPlayer.Username,
		&createdPlayer.CreatedAt,
		&createdPlayer.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.Player{}, ErrUsernameAlreadyUsed
		}
		return models.Player{}, errors.New("internal server error")
	}
	return createdPlayer, nil
}
