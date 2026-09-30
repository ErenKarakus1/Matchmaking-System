package models

import (
	"time"

	"github.com/google/uuid"
)

type Rating struct {
	PlayerID    uuid.UUID `json:"player_id"`
	Rating      int       `json:"rating"`
	GamesPlayed int       `json:"games_played"`
	Wins        int       `json:"wins"`
	Losses      int       `json:"losses"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
