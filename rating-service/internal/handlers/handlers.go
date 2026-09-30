package handlers

import (
	"errors"
	"net/http"

	"github.com/ErenKarakus1/Matchmaking-System/rating-service/internal/repository"
	"github.com/ErenKarakus1/Matchmaking-System/rating-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreatePlayerRatingHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		playerID := ctx.Param("player_id")
		if playerID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "player id is required"})
			return
		}
		parsedPlayerID, err := uuid.Parse(playerID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid player id"})
			return
		}
		rating, err := service.CreatePlayerRating(ctx.Request.Context(), pool, parsedPlayerID)
		if err != nil {
			if errors.Is(err, repository.ErrRatingAlreadyExists) {
				ctx.JSON(http.StatusConflict, gin.H{"error": "rating already exists"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		ctx.JSON(http.StatusCreated, rating)
	}
}

func GetPlayerRatingHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		playerID := ctx.Param("player_id")
		if playerID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "player id is required"})
			return
		}
		parsedPlayerID, err := uuid.Parse(playerID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid player id"})
			return
		}
		rating, err := service.GetPlayerRatingByPlayerID(ctx.Request.Context(), pool, parsedPlayerID)
		if err != nil {
			if errors.Is(err, repository.ErrRatingNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "rating not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		ctx.JSON(http.StatusOK, rating)
	}
}
