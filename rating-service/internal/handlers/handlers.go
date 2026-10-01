package handlers

import (
	"errors"
	"net/http"

	"github.com/ErenKarakus1/Matchmaking-System/rating-service/internal/models"
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

func CreateMatchHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req models.CreateMatchRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		if req.WinnerID == req.LoserID {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "winner id and loser id cannot be same"})
			return
		}
		matchID := ctx.Param("match_id")
		if matchID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "match id is required"})
			return
		}
		parsedMatchID, err := uuid.Parse(matchID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
			return
		}
		match, err := service.CreateMatch(ctx.Request.Context(), pool, parsedMatchID, req.WinnerID, req.LoserID)
		if err != nil {
			if errors.Is(err, repository.ErrRatingNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "winner or loser id was not found"})
				return
			}
			if errors.Is(err, repository.ErrMatchAlreadyExists) {
				ctx.JSON(http.StatusConflict, gin.H{"error": "match already exists"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		ctx.JSON(http.StatusCreated, match)
	}
}

func GetMatchHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		matchID := ctx.Param("match_id")
		if matchID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "match id is required"})
			return
		}
		parsedMatchID, err := uuid.Parse(matchID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
			return
		}
		match, err := service.GetMatch(ctx.Request.Context(), pool, parsedMatchID)
		if err != nil {
			if errors.Is(err, repository.ErrMatchNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		ctx.JSON(http.StatusOK, match)
	}
}
