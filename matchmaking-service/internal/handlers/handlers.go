package handlers

import (
	"net/http"

	"github.com/ErenKarakus1/Matchmaking-System/matchmaking-service/internal/models"
	"github.com/ErenKarakus1/Matchmaking-System/matchmaking-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func CreateTicketHandler(client *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req models.CreateTicketRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		ticket, err := service.CreateTicket(ctx.Request.Context(), client, req.PlayerID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		ctx.JSON(http.StatusCreated, ticket)
	}
}
