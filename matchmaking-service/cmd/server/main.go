package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ErenKarakus1/Matchmaking-System/matchmaking-service/internal/config"
	"github.com/ErenKarakus1/Matchmaking-System/matchmaking-service/internal/handlers"
	"github.com/ErenKarakus1/Matchmaking-System/matchmaking-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.LoadConfig()

	options, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(options)
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go startMatchmaker(ctx, client)

	router := gin.Default()

	router.POST("/matchmaking/tickets", handlers.CreateTicketHandler(client))
	router.GET("/matchmaking/tickets/:ticket_id", handlers.GetTicketHandler(client))
	router.DELETE("/matchmaking/tickets/:ticket_id", handlers.DeleteTicketHandler(client))
	router.GET("/matchmaking/queue", handlers.GetQueueHandler(client))
	router.POST("/matchmaking/matches", handlers.CreateMatchHandler(client))
	router.GET("/matchmaking/matches/:match_id", handlers.GetMatchHandler(client))

	if err := router.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}

func startMatchmaker(ctx context.Context, client *redis.Client) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			match, err := service.CreateMatch(ctx, client)
			if err != nil {
				if errors.Is(err, service.ErrNotEnoughPlayers) {
					continue
				}
				log.Println("matchmaker error: ", err)
				continue
			}
			log.Println("created match: ", match.MatchID)
		}
	}
}
