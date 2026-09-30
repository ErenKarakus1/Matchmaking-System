package handlers

import (
	"errors"
	"net/http"

	"github.com/ErenKarakus1/Matchmaking-System/player-service/internal/models"
	"github.com/ErenKarakus1/Matchmaking-System/player-service/internal/repository"
	"github.com/ErenKarakus1/Matchmaking-System/player-service/internal/service"
	"github.com/ErenKarakus1/Matchmaking-System/player-service/internal/validation"
	"github.com/gin-gonic/gin"
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
