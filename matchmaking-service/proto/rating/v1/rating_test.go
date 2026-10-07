package ratingv1

import "testing"

func TestGetPlayerRatingRequestGetters(t *testing.T) {
	req := &GetPlayerRatingRequest{PlayerId: "player-1"}

	if req.GetPlayerId() != "player-1" {
		t.Fatalf("GetPlayerId() = %q, want %q", req.GetPlayerId(), "player-1")
	}

	var nilReq *GetPlayerRatingRequest
	if nilReq.GetPlayerId() != "" {
		t.Fatalf("nil GetPlayerId() = %q, want empty string", nilReq.GetPlayerId())
	}
}

func TestGetPlayerRatingResponseGetters(t *testing.T) {
	resp := &GetPlayerRatingResponse{
		PlayerId:    "player-1",
		Rating:      1510,
		GamesPlayed: 4,
		Wins:        3,
		Losses:      1,
	}

	if resp.GetPlayerId() != "player-1" {
		t.Fatalf("GetPlayerId() = %q, want %q", resp.GetPlayerId(), "player-1")
	}
	if resp.GetRating() != 1510 {
		t.Fatalf("GetRating() = %d, want 1510", resp.GetRating())
	}
	if resp.GetGamesPlayed() != 4 {
		t.Fatalf("GetGamesPlayed() = %d, want 4", resp.GetGamesPlayed())
	}
	if resp.GetWins() != 3 {
		t.Fatalf("GetWins() = %d, want 3", resp.GetWins())
	}
	if resp.GetLosses() != 1 {
		t.Fatalf("GetLosses() = %d, want 1", resp.GetLosses())
	}

	var nilResp *GetPlayerRatingResponse
	if nilResp.GetPlayerId() != "" {
		t.Fatalf("nil GetPlayerId() = %q, want empty string", nilResp.GetPlayerId())
	}
	if nilResp.GetRating() != 0 {
		t.Fatalf("nil GetRating() = %d, want 0", nilResp.GetRating())
	}
	if nilResp.GetGamesPlayed() != 0 {
		t.Fatalf("nil GetGamesPlayed() = %d, want 0", nilResp.GetGamesPlayed())
	}
	if nilResp.GetWins() != 0 {
		t.Fatalf("nil GetWins() = %d, want 0", nilResp.GetWins())
	}
	if nilResp.GetLosses() != 0 {
		t.Fatalf("nil GetLosses() = %d, want 0", nilResp.GetLosses())
	}
}
