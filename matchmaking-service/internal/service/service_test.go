package service

import (
	"context"
	"errors"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type fakeRatingProvider struct {
	ratings map[uuid.UUID]int
	err     error
}

func (f fakeRatingProvider) GetPlayerRating(ctx context.Context, playerID uuid.UUID) (int, error) {
	if f.err != nil {
		return 0, f.err
	}
	return f.ratings[playerID], nil
}

func newTestRedis(t *testing.T) *redis.Client {
	t.Helper()

	server := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	t.Cleanup(func() {
		_ = client.Close()
		server.Close()
	})

	return client
}

func TestAbs(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  int
	}{
		{
			name:  "positive value",
			value: 25,
			want:  25,
		},
		{
			name:  "negative value",
			value: -25,
			want:  25,
		},
		{
			name:  "zero",
			value: 0,
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := abs(tt.value)
			if got != tt.want {
				t.Fatalf("abs(%d) = %d, want %d", tt.value, got, tt.want)
			}
		})
	}
}

func TestGenerateTicket(t *testing.T) {
	playerID := uuid.New()

	ticket := generateTicket(playerID)

	if ticket.TicketID == uuid.Nil {
		t.Fatal("expected generated ticket id")
	}
	if ticket.PlayerID != playerID {
		t.Fatalf("PlayerID = %s, want %s", ticket.PlayerID, playerID)
	}
	if ticket.Status != TicketStatusQueued {
		t.Fatalf("Status = %q, want %q", ticket.Status, TicketStatusQueued)
	}
	if ticket.MatchID != nil {
		t.Fatalf("MatchID = %s, want nil", ticket.MatchID)
	}
	if ticket.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be set")
	}
	if ticket.CreatedAt.Location().String() != "UTC" {
		t.Fatalf("CreatedAt location = %s, want UTC", ticket.CreatedAt.Location())
	}
}

func TestCreateTicketCreatesQueuedTicketAndPreventsDuplicatePlayer(t *testing.T) {
	ctx := context.Background()
	client := newTestRedis(t)
	playerID := uuid.New()

	ticket, err := CreateTicket(ctx, client, playerID)
	if err != nil {
		t.Fatalf("CreateTicket returned error: %v", err)
	}

	if ticket.PlayerID != playerID {
		t.Fatalf("PlayerID = %s, want %s", ticket.PlayerID, playerID)
	}
	if ticket.Status != TicketStatusQueued {
		t.Fatalf("Status = %q, want %q", ticket.Status, TicketStatusQueued)
	}

	storedTicket, err := GetTicket(ctx, client, ticket.TicketID)
	if err != nil {
		t.Fatalf("GetTicket returned error: %v", err)
	}
	if storedTicket.TicketID != ticket.TicketID {
		t.Fatalf("stored ticket id = %s, want %s", storedTicket.TicketID, ticket.TicketID)
	}

	_, err = CreateTicket(ctx, client, playerID)
	if !errors.Is(err, ErrAlreadyQueued) {
		t.Fatalf("duplicate CreateTicket error = %v, want %v", err, ErrAlreadyQueued)
	}
}

func TestGetTicketReturnsNotFound(t *testing.T) {
	client := newTestRedis(t)

	_, err := GetTicket(context.Background(), client, uuid.New())
	if !errors.Is(err, ErrTicketNotFound) {
		t.Fatalf("GetTicket error = %v, want %v", err, ErrTicketNotFound)
	}
}

func TestDeleteTicketRemovesQueuedTicket(t *testing.T) {
	ctx := context.Background()
	client := newTestRedis(t)
	playerID := uuid.New()

	ticket, err := CreateTicket(ctx, client, playerID)
	if err != nil {
		t.Fatalf("CreateTicket returned error: %v", err)
	}

	if err := DeleteTicket(ctx, client, ticket.TicketID); err != nil {
		t.Fatalf("DeleteTicket returned error: %v", err)
	}

	_, err = GetTicket(ctx, client, ticket.TicketID)
	if !errors.Is(err, ErrTicketNotFound) {
		t.Fatalf("GetTicket after delete error = %v, want %v", err, ErrTicketNotFound)
	}

	_, err = CreateTicket(ctx, client, playerID)
	if err != nil {
		t.Fatalf("CreateTicket after delete returned error: %v", err)
	}
}

func TestCreateMatchSelectsClosestRatedPlayers(t *testing.T) {
	ctx := context.Background()
	client := newTestRedis(t)

	playerA := uuid.New()
	playerB := uuid.New()
	playerC := uuid.New()

	ticketA, err := CreateTicket(ctx, client, playerA)
	if err != nil {
		t.Fatalf("CreateTicket A returned error: %v", err)
	}
	ticketB, err := CreateTicket(ctx, client, playerB)
	if err != nil {
		t.Fatalf("CreateTicket B returned error: %v", err)
	}
	ticketC, err := CreateTicket(ctx, client, playerC)
	if err != nil {
		t.Fatalf("CreateTicket C returned error: %v", err)
	}

	match, err := CreateMatch(ctx, client, fakeRatingProvider{
		ratings: map[uuid.UUID]int{
			playerA: 1000,
			playerB: 1800,
			playerC: 1810,
		},
	})
	if err != nil {
		t.Fatalf("CreateMatch returned error: %v", err)
	}

	if match.PlayerAID != playerB || match.PlayerBID != playerC {
		t.Fatalf("matched players = (%s, %s), want (%s, %s)", match.PlayerAID, match.PlayerBID, playerB, playerC)
	}
	if match.TicketAID != ticketB.TicketID || match.TicketBID != ticketC.TicketID {
		t.Fatalf("matched tickets = (%s, %s), want (%s, %s)", match.TicketAID, match.TicketBID, ticketB.TicketID, ticketC.TicketID)
	}

	matchedB, err := GetTicket(ctx, client, ticketB.TicketID)
	if err != nil {
		t.Fatalf("GetTicket B returned error: %v", err)
	}
	if matchedB.Status != TicketStatusMatched {
		t.Fatalf("ticket B status = %q, want %q", matchedB.Status, TicketStatusMatched)
	}
	if matchedB.MatchID == nil || *matchedB.MatchID != match.MatchID {
		t.Fatalf("ticket B match id = %v, want %s", matchedB.MatchID, match.MatchID)
	}

	queuedTickets, err := GetQueue(ctx, client)
	if err != nil {
		t.Fatalf("GetQueue returned error: %v", err)
	}
	if len(queuedTickets) != 1 {
		t.Fatalf("queued ticket count = %d, want 1", len(queuedTickets))
	}
	if queuedTickets[0].TicketID != ticketA.TicketID {
		t.Fatalf("remaining ticket = %s, want %s", queuedTickets[0].TicketID, ticketA.TicketID)
	}
}

func TestCreateMatchReturnsNotEnoughPlayers(t *testing.T) {
	ctx := context.Background()
	client := newTestRedis(t)

	_, err := CreateMatch(ctx, client, fakeRatingProvider{ratings: map[uuid.UUID]int{}})
	if !errors.Is(err, ErrNotEnoughPlayers) {
		t.Fatalf("CreateMatch error = %v, want %v", err, ErrNotEnoughPlayers)
	}
}

func TestCreateMatchReturnsInternalErrorWhenRatingProviderFails(t *testing.T) {
	ctx := context.Background()
	client := newTestRedis(t)

	if _, err := CreateTicket(ctx, client, uuid.New()); err != nil {
		t.Fatalf("CreateTicket A returned error: %v", err)
	}
	if _, err := CreateTicket(ctx, client, uuid.New()); err != nil {
		t.Fatalf("CreateTicket B returned error: %v", err)
	}

	_, err := CreateMatch(ctx, client, fakeRatingProvider{err: errors.New("rating service down")})
	if err == nil {
		t.Fatal("CreateMatch returned nil error, want error")
	}
}
