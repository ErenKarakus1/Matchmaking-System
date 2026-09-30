package service

import (
	"context"
	"errors"

	"github.com/ErenKarakus1/Matchmaking-System/rating-service/internal/models"
	"github.com/ErenKarakus1/Matchmaking-System/rating-service/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreatePlayerRating(ctx context.Context, pool *pgxpool.Pool, playerID uuid.UUID) (models.Rating, error) {
	rating, err := repository.InsertRating(ctx, pool, playerID)
	if err != nil {
		if errors.Is(err, repository.ErrRatingAlreadyExists) {
			return models.Rating{}, err
		}
		return models.Rating{}, errors.New("internal server error")
	}
	return rating, nil
}
