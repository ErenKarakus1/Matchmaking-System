package handlers

import (
	"errors"
	"net/http"

	"github.com/ErenKarakus1/Matchmaking-System/matchmaking-service/internal/models"
	"github.com/ErenKarakus1/Matchmaking-System/matchmaking-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
			if errors.Is(err, service.ErrAlreadyQueued) {
				ctx.JSON(http.StatusConflict, gin.H{"error": "player is already queued"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		ctx.JSON(http.StatusCreated, ticket)
	}
}

func GetTicketHandler(client *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ticketID := ctx.Param("ticket_id")
		if ticketID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "ticket id is required"})
			return
		}

		parsedTicketID, err := uuid.Parse(ticketID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
			return
		}

		ticket, err := service.GetTicket(ctx.Request.Context(), client, parsedTicketID)
		if err != nil {
			if errors.Is(err, service.ErrTicketNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		ctx.JSON(http.StatusOK, ticket)
	}
}

func DeleteTicketHandler(client *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ticketID := ctx.Param("ticket_id")
		if ticketID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "ticket id is required"})
			return
		}

		parsedTicketID, err := uuid.Parse(ticketID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
			return
		}

		err = service.DeleteTicket(ctx.Request.Context(), client, parsedTicketID)
		if err != nil {
			if errors.Is(err, service.ErrTicketNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
				return
			}
			if errors.Is(err, service.ErrTicketNotQueued) {
				ctx.JSON(http.StatusConflict, gin.H{"error": "ticket is not queued"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		ctx.Status(http.StatusNoContent)
	}
}
