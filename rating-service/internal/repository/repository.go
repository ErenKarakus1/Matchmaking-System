package repository

import (
	"context"
	"errors"

	"github.com/ErenKarakus1/Matchmaking-System/rating-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRatingAlreadyExists = errors.New("rating already exists")

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
