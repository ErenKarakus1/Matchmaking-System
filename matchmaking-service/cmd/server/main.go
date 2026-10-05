package main

import (
	"log"

	"github.com/ErenKarakus1/Matchmaking-System/matchmaking-service/internal/config"
	"github.com/ErenKarakus1/Matchmaking-System/matchmaking-service/internal/handlers"
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
