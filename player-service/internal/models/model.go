package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Player struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type CreatePlayerRequest struct {
	Username string `json:"username"`
}

func (r *CreatePlayerRequest) Normalize() {
	r.Username = strings.ToLower(strings.TrimSpace(r.Username))
}
