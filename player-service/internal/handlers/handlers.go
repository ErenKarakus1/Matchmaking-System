package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ErenKarakus1/Matchmaking-System/player-service/internal/models"
	"github.com/ErenKarakus1/Matchmaking-System/player-service/internal/repository"
	"github.com/ErenKarakus1/Matchmaking-System/player-service/internal/service"
	"github.com/ErenKarakus1/Matchmaking-System/player-service/internal/validation"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreatePlayerHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req models.CreatePlayerRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		req.Normalize()
		if err := validation.ValidateCreatePlayerRequest(req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		player, err := service.CreatePlayer(ctx.Request.Context(), pool, req)
		if err != nil {
			if errors.Is(err, repository.ErrUsernameAlreadyUsed) {
				ctx.JSON(http.StatusConflict, gin.H{"error": "username is already used"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		ctx.JSON(http.StatusCreated, player)
	}
}

func GetPlayerByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		playerID := ctx.Param("id")
		if playerID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "player id is required"})
			return
		}
		parsedPlayerID, err := uuid.Parse(playerID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid player id"})
			return
		}
		player, err := service.GetPlayerByID(ctx.Request.Context(), pool, parsedPlayerID)
		if err != nil {
			if errors.Is(err, repository.ErrPlayerNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		ctx.JSON(http.StatusOK, player)
	}
}

func GetPlayerByUsernameHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.Param("username")
		if username == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
			return
		}
		username = strings.ToLower(strings.TrimSpace(username))
		if err := validation.ValidateUsername(username); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		player, err := service.GetPlayerByUsername(ctx.Request.Context(), pool, username)
		if err != nil {
			if errors.Is(err, repository.ErrPlayerNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		ctx.JSON(http.StatusOK, player)
	}
}
