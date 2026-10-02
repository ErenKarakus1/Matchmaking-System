package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/ErenKarakus1/Matchmaking-System/matchmaking-service/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	TicketStatusQueued    = "queued"
	TicketStatusMatched   = "matched"
	TicketStatusCancelled = "cancelled"
)

func generateTicket(playerID uuid.UUID) models.Ticket {
	return models.Ticket{
		TicketID:  uuid.New(),
		PlayerID:  playerID,
		Status:    TicketStatusQueued,
		CreatedAt: time.Now(),
	}
}

func CreateTicket(ctx context.Context, client *redis.Client, playerID uuid.UUID) (models.Ticket, error) {
	ticket := generateTicket(playerID)
	ticketBytes, err := json.Marshal(ticket)
	if err != nil {
		return models.Ticket{}, errors.New("internal server error")
	}
	key := "ticket:" + ticket.TicketID.String()
	if err := client.Set(ctx, key, ticketBytes, 0).Err(); err != nil {
		return models.Ticket{}, errors.New("internal server error")
	}
	return ticket, nil
}
