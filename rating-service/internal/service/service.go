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

func GetPlayerRatingByPlayerID(ctx context.Context, pool *pgxpool.Pool, playerID uuid.UUID) (models.Rating, error) {
	rating, err := repository.GetRatingByPlayerID(ctx, pool, playerID)
	if err != nil {
		if errors.Is(err, repository.ErrRatingNotFound) {
			return models.Rating{}, err
		}
		return models.Rating{}, errors.New("internal server error")
	}
	return rating, nil
}

func CreateMatch(ctx context.Context, pool *pgxpool.Pool, matchID uuid.UUID, winnerID uuid.UUID, loserID uuid.UUID) (models.Match, error) {
	match, err := repository.SubmitMatchResult(ctx, pool, matchID, winnerID, loserID)
	if err != nil {
		if errors.Is(err, repository.ErrRatingNotFound) || errors.Is(err, repository.ErrMatchAlreadyExists) {
			return models.Match{}, err
		}
		return models.Match{}, errors.New("internal server error")
	}
	return match, nil
}

func GetMatch(ctx context.Context, pool *pgxpool.Pool, matchID uuid.UUID) (models.Match, error) {
	match, err := repository.GetMatch(ctx, pool, matchID)
	if err != nil {
		if errors.Is(err, repository.ErrMatchNotFound) {
			return models.Match{}, err
		}
		return models.Match{}, errors.New("internal server error")
	}
	return match, nil
}

func GetLeaderboard(ctx context.Context, pool *pgxpool.Pool) ([]models.LeaderboardEntry, error) {
	leaderboard, err := repository.GetLeaderboard(ctx, pool)
	if err != nil {
		return []models.LeaderboardEntry{}, errors.New("internal server error")
	}
	return leaderboard, nil
}

func GetPlayerMatches(ctx context.Context, pool *pgxpool.Pool, playerID uuid.UUID) ([]models.Match, error) {
	_, err := repository.GetRatingByPlayerID(ctx, pool, playerID)
	if err != nil {
		if errors.Is(err, repository.ErrRatingNotFound) {
			return []models.Match{}, err
		}
		return []models.Match{}, errors.New("internal server error")
	}

	matches, err := repository.GetPlayerMatches(ctx, pool, playerID)
	if err != nil {
		return []models.Match{}, errors.New("internal server error")
	}

	return matches, nil
}
