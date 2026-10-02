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

var ErrAlreadyQueued = errors.New("player is already queued")

func generateTicket(playerID uuid.UUID) models.Ticket {
	return models.Ticket{
		TicketID:  uuid.New(),
		PlayerID:  playerID,
		Status:    TicketStatusQueued,
		CreatedAt: time.Now(),
	}
}

func CreateTicket(ctx context.Context, client *redis.Client, playerID uuid.UUID) (models.Ticket, error) {
	playerTicketKey := "player_ticket:" + playerID.String()
	_, err := client.Get(ctx, playerTicketKey).Result()
	if err == nil {
		return models.Ticket{}, ErrAlreadyQueued
	}
	if err != nil && !errors.Is(err, redis.Nil) {
		return models.Ticket{}, errors.New("internal server error")
	}

	pipe := client.TxPipeline()

	ticket := generateTicket(playerID)
	ticketBytes, err := json.Marshal(ticket)
	if err != nil {
		return models.Ticket{}, errors.New("internal server error")
	}
	key := "ticket:" + ticket.TicketID.String()
	pipe.Set(ctx, key, ticketBytes, 0)
	pipe.Set(ctx, playerTicketKey, ticket.TicketID.String(), 0)
	pipe.ZAdd(ctx, "queue", redis.Z{
		Score:  float64(ticket.CreatedAt.Unix()),
		Member: ticket.TicketID.String(),
	})

	_, err = pipe.Exec(ctx)
	if err != nil {
		return models.Ticket{}, errors.New("internal server error")
	}

	return ticket, nil
}
