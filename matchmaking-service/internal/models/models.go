package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateTicketRequest struct {
	PlayerID uuid.UUID `json:"player_id"`
}

type Ticket struct {
	TicketID  uuid.UUID `json:"ticket_id"`
	PlayerID  uuid.UUID `json:"player_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
