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

type Match struct {
	MatchID   uuid.UUID `json:"match_id"`
	PlayerAID uuid.UUID `json:"player_a_id"`
	PlayerBID uuid.UUID `json:"player_b_id"`
	TicketAID uuid.UUID `json:"ticket_a_id"`
	TicketBID uuid.UUID `json:"ticket_b_id"`
	CreatedAt time.Time `json:"created_at"`
}
