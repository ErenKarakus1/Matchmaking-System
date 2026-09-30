package repository

import (
	"context"
	"errors"

	"github.com/ErenKarakus1/Matchmaking-System/rating-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRatingAlreadyExists = errors.New("rating already exists")
var ErrRatingNotFound = errors.New("rating not found")

const insertRatingSQL = `
	INSERT INTO ratings (
		player_id
	)
	VALUES ($1)
	RETURNING
		player_id,
		rating,
		games_played,
		wins,
		losses,
		created_at,
		updated_at
`

const getRatingByPlayerIDSQL = `
	SELECT
		player_id,
		rating,
		games_played,
		wins,
		losses,
		created_at,
		updated_at
	FROM ratings
	WHERE player_id=$1
`

func InsertRating(ctx context.Context, pool *pgxpool.Pool, playerID uuid.UUID) (models.Rating, error) {
	var rating models.Rating
	err := pool.QueryRow(
		ctx,
		insertRatingSQL,
		playerID,
	).Scan(
		&rating.PlayerID,
		&rating.Rating,
		&rating.GamesPlayed,
		&rating.Wins,
		&rating.Losses,
		&rating.CreatedAt,
		&rating.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.Rating{}, ErrRatingAlreadyExists
		}
		return models.Rating{}, errors.New("internal server error")
	}
	return rating, nil
}

func GetRatingByPlayerID(ctx context.Context, pool *pgxpool.Pool, playerID uuid.UUID) (models.Rating, error) {
	var rating models.Rating
	err := pool.QueryRow(
		ctx,
		getRatingByPlayerIDSQL,
		playerID,
	).Scan(
		&rating.PlayerID,
		&rating.Rating,
		&rating.GamesPlayed,
		&rating.Wins,
		&rating.Losses,
		&rating.CreatedAt,
		&rating.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Rating{}, ErrRatingNotFound
		}
		return models.Rating{}, errors.New("internal server error")
	}
	return rating, nil
}
