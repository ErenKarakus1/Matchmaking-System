package service

import (
	"context"
	"errors"

	"github.com/ErenKarakus1/Matchmaking-System/player-service/internal/models"
	"github.com/ErenKarakus1/Matchmaking-System/player-service/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func createPlayerWithID(req models.CreatePlayerRequest) models.Player {
	return models.Player{
		ID:       uuid.New(),
		Username: req.Username,
	}
}

func CreatePlayer(ctx context.Context, pool *pgxpool.Pool, req models.CreatePlayerRequest) (models.Player, error) {
	player := createPlayerWithID(req)
	createdPlayer, err := repository.InsertPlayer(ctx, pool, player)
	if err != nil {
		if errors.Is(err, repository.ErrUsernameAlreadyUsed) {
			return models.Player{}, err
		}
		return models.Player{}, errors.New("internal server error")
	}
	return createdPlayer, nil
}

func GetPlayerByID(ctx context.Context, pool *pgxpool.Pool, playerID uuid.UUID) (models.Player, error) {
	player, err := repository.GetPlayerByID(ctx, pool, playerID)
	if err != nil {
		if errors.Is(err, repository.ErrPlayerNotFound) {
			return models.Player{}, err
		}
		return models.Player{}, errors.New("internal server error")
	}
	return player, nil
}

func GetPlayerByUsername(ctx context.Context, pool *pgxpool.Pool, username string) (models.Player, error) {
	player, err := repository.GetPlayerByUsername(ctx, pool, username)
	if err != nil {
		if errors.Is(err, repository.ErrPlayerNotFound) {
			return models.Player{}, err
		}
		return models.Player{}, errors.New("internal server error")
	}
	return player, nil
}
