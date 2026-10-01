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

type Match struct {
	ID                 uuid.UUID `json:"id"`
	WinnerID           uuid.UUID `json:"winner_id"`
	LoserID            uuid.UUID `json:"loser_id"`
	WinnerRatingBefore int       `json:"winner_rating_before"`
	LoserRatingBefore  int       `json:"loser_rating_before"`
	WinnerRatingAfter  int       `json:"winner_rating_after"`
	LoserRatingAfter   int       `json:"loser_rating_after"`
	CreatedAt          time.Time `json:"created_at"`
}

type CreateMatchRequest struct {
	WinnerID uuid.UUID `json:"winner_id"`
	LoserID  uuid.UUID `json:"loser_id"`
}

type LeaderboardEntry struct {
	Rank        int       `json:"rank"`
	PlayerID    uuid.UUID `json:"player_id"`
	Rating      int       `json:"rating"`
	GamesPlayed int       `json:"games_played"`
	Wins        int       `json:"wins"`
	Losses      int       `json:"losses"`
}
