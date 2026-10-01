package main

import (
	"log"

	"github.com/ErenKarakus1/Matchmaking-System/rating-service/internal/config"
	"github.com/ErenKarakus1/Matchmaking-System/rating-service/internal/db"
	"github.com/ErenKarakus1/Matchmaking-System/rating-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	pool, err := db.NewPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to Postgres")

	router := gin.Default()

	router.POST("/ratings/players/:player_id", handlers.CreatePlayerRatingHandler(pool))
	router.GET("/ratings/players/:player_id", handlers.GetPlayerRatingHandler(pool))
	router.POST("/ratings/matches/:match_id/result", handlers.CreateMatchHandler(pool))
	router.GET("/ratings/matches/:match_id", handlers.GetMatchHandler(pool))

	if err := router.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
