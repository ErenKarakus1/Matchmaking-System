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
	TicketStatusQueued  = "queued"
	TicketStatusMatched = "matched"
)

type RatingProvider interface {
	GetPlayerRating(ctx context.Context, playerID uuid.UUID) (int, error)
}

var ErrAlreadyQueued = errors.New("player is already queued")
var ErrTicketNotFound = errors.New("ticket not found")
var ErrTicketNotQueued = errors.New("ticket is not queued")
var ErrNotEnoughPlayers = errors.New("not enough players")
var ErrMatchNotFound = errors.New("match not found")

func generateTicket(playerID uuid.UUID) models.Ticket {
	return models.Ticket{
		TicketID:  uuid.New(),
		PlayerID:  playerID,
		Status:    TicketStatusQueued,
		CreatedAt: time.Now().UTC(),
	}
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func selectTicketsForMatch(ctx context.Context, client *redis.Client, ratingProvider RatingProvider) ([]models.Ticket, error) {
	ticketIDs, err := client.ZRange(ctx, "queue", 0, 9).Result()
	if err != nil {
		return []models.Ticket{}, errors.New("internal server error")
	}
	if len(ticketIDs) < 2 {
		return []models.Ticket{}, ErrNotEnoughPlayers
	}
	var tickets []models.Ticket
	for _, ticketID := range ticketIDs {
		key := "ticket:" + ticketID
		ticketBytes, err := client.Get(ctx, key).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}
			return []models.Ticket{}, errors.New("internal server error")
		}
		var ticket models.Ticket
		if err := json.Unmarshal([]byte(ticketBytes), &ticket); err != nil {
			return []models.Ticket{}, errors.New("internal server error")
		}
		tickets = append(tickets, ticket)
	}
	if len(tickets) < 2 {
		return []models.Ticket{}, ErrNotEnoughPlayers
	}

	bestA := 0
	bestB := 1
	bestDiff := -1

	for i := 0; i < len(tickets); i++ {
		ratingI, err := ratingProvider.GetPlayerRating(ctx, tickets[i].PlayerID)
		if err != nil {
			return []models.Ticket{}, errors.New("internal server error")
		}

		for j := i + 1; j < len(tickets); j++ {
			ratingJ, err := ratingProvider.GetPlayerRating(ctx, tickets[j].PlayerID)
			if err != nil {
				return []models.Ticket{}, errors.New("internal server error")
			}

			diff := abs(ratingI - ratingJ)
			if bestDiff == -1 || diff < bestDiff {
				bestA = i
				bestB = j
				bestDiff = diff
			}
		}
	}

	return []models.Ticket{tickets[bestA], tickets[bestB]}, nil
}

func CreateTicket(ctx context.Context, client *redis.Client, playerID uuid.UUID) (models.Ticket, error) {
	playerTicketKey := "player_ticket:" + playerID.String()
	_, err := client.Get(ctx, playerTicketKey).Result()
	if err == nil {
		return models.Ticket{}, ErrAlreadyQueued
	}
	if !errors.Is(err, redis.Nil) {
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

func GetTicket(ctx context.Context, client *redis.Client, ticketID uuid.UUID) (models.Ticket, error) {
	ticketKey := "ticket:" + ticketID.String()
	ticketBytes, err := client.Get(ctx, ticketKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return models.Ticket{}, ErrTicketNotFound
		}
		return models.Ticket{}, errors.New("internal server error")
	}
	var ticket models.Ticket
	err = json.Unmarshal([]byte(ticketBytes), &ticket)
	if err != nil {
		return models.Ticket{}, errors.New("internal server error")
	}

	return ticket, nil
}

func DeleteTicket(ctx context.Context, client *redis.Client, ticketID uuid.UUID) error {
	ticketKey := "ticket:" + ticketID.String()
	ticketBytes, err := client.Get(ctx, ticketKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrTicketNotFound
		}
		return errors.New("internal server error")
	}
	var ticket models.Ticket
	err = json.Unmarshal([]byte(ticketBytes), &ticket)
	if err != nil {
		return errors.New("internal server error")
	}

	if ticket.Status != TicketStatusQueued {
		return ErrTicketNotQueued
	}

	pipe := client.TxPipeline()

	pipe.Del(ctx, ticketKey)
	playerTicketKey := "player_ticket:" + ticket.PlayerID.String()
	pipe.Del(ctx, playerTicketKey)
	pipe.ZRem(ctx, "queue", ticketID.String())

	_, err = pipe.Exec(ctx)
	if err != nil {
		return errors.New("internal server error")
	}

	return nil
}

func GetQueue(ctx context.Context, client *redis.Client) ([]models.Ticket, error) {
	ticketIDs, err := client.ZRange(ctx, "queue", 0, 49).Result()
	if err != nil {
		return []models.Ticket{}, errors.New("internal server error")
	}
	var queue []models.Ticket
	for _, ticketID := range ticketIDs {
		key := "ticket:" + ticketID
		ticketBytes, err := client.Get(ctx, key).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}
			return []models.Ticket{}, errors.New("internal server error")
		}
		var ticket models.Ticket
		err = json.Unmarshal([]byte(ticketBytes), &ticket)
		if err != nil {
			return []models.Ticket{}, errors.New("internal server error")
		}
		queue = append(queue, ticket)
	}
	return queue, nil
}

func CreateMatch(ctx context.Context, client *redis.Client, ratingProvider RatingProvider) (models.Match, error) {
	tickets, err := selectTicketsForMatch(ctx, client, ratingProvider)
	if err != nil {
		if errors.Is(err, ErrNotEnoughPlayers) {
			return models.Match{}, err
		}
		return models.Match{}, errors.New("internal server error")
	}
	ticketA := tickets[0]
	ticketB := tickets[1]

	match := models.Match{
		MatchID:   uuid.New(),
		PlayerAID: ticketA.PlayerID,
		PlayerBID: ticketB.PlayerID,
		TicketAID: ticketA.TicketID,
		TicketBID: ticketB.TicketID,
		CreatedAt: time.Now().UTC(),
	}
	matchBytes, err := json.Marshal(match)
	if err != nil {
		return models.Match{}, errors.New("internal server error")
	}

	ticketA.Status = TicketStatusMatched
	ticketA.MatchID = &match.MatchID
	ticketB.Status = TicketStatusMatched
	ticketB.MatchID = &match.MatchID

	ticketABytes, err := json.Marshal(ticketA)
	if err != nil {
		return models.Match{}, errors.New("internal server error")
	}
	ticketBBytes, err := json.Marshal(ticketB)
	if err != nil {
		return models.Match{}, errors.New("internal server error")
	}

	pipe := client.TxPipeline()
	pipe.Set(ctx, "ticket:"+ticketA.TicketID.String(), ticketABytes, 0)
	pipe.Set(ctx, "ticket:"+ticketB.TicketID.String(), ticketBBytes, 0)
	pipe.ZRem(ctx, "queue", ticketA.TicketID.String())
	pipe.ZRem(ctx, "queue", ticketB.TicketID.String())
	pipe.Del(ctx, "player_ticket:"+ticketA.PlayerID.String())
	pipe.Del(ctx, "player_ticket:"+ticketB.PlayerID.String())
	pipe.Set(ctx, "match:"+match.MatchID.String(), matchBytes, 0)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return models.Match{}, errors.New("internal server error")
	}
	return match, nil
}

func GetMatch(ctx context.Context, client *redis.Client, matchID uuid.UUID) (models.Match, error) {
	matchBytes, err := client.Get(ctx, "match:"+matchID.String()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return models.Match{}, ErrMatchNotFound
		}
		return models.Match{}, errors.New("internal server error")
	}
	var match models.Match
	if err := json.Unmarshal([]byte(matchBytes), &match); err != nil {
		return models.Match{}, errors.New("internal server error")
	}
	return match, nil
}
