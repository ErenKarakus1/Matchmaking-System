package repository

import (
	"context"
	"errors"

	"github.com/ErenKarakus1/Matchmaking-System/player-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUsernameAlreadyUsed = errors.New("username is already used")
var ErrPlayerNotFound = errors.New("player not found")

const insertPlayerSQL = `
	INSERT INTO players (
		id,
		username
	)
	VALUES ($1,$2)
	RETURNING
		id,
		username,
		created_at
`

const getPlayerByIDSQL = `
	SELECT
		id,
		username,
		created_at
	FROM players
	WHERE id=$1
`

const getPlayerByUsernameSQL = `
	SELECT
		id,
		username,
		created_at
	FROM players
	WHERE username=$1
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

func GetPlayerByID(ctx context.Context, pool *pgxpool.Pool, playerID uuid.UUID) (models.Player, error) {
	var player models.Player
	err := pool.QueryRow(
		ctx,
		getPlayerByIDSQL,
		playerID,
	).Scan(
		&player.ID,
		&player.Username,
		&player.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Player{}, ErrPlayerNotFound
		}
		return models.Player{}, errors.New("internal server error")
	}
	return player, nil
}

func GetPlayerByUsername(ctx context.Context, pool *pgxpool.Pool, username string) (models.Player, error) {
	var player models.Player
	err := pool.QueryRow(
		ctx,
		getPlayerByUsernameSQL,
		username,
	).Scan(
		&player.ID,
		&player.Username,
		&player.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Player{}, ErrPlayerNotFound
		}
		return models.Player{}, errors.New("internal server error")
	}
	return player, nil
}
