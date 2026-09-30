package main

import (
	"log"

	"github.com/ErenKarakus1/Matchmaking-System/player-service/internal/config"
	"github.com/ErenKarakus1/Matchmaking-System/player-service/internal/db"
	"github.com/ErenKarakus1/Matchmaking-System/player-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	pool, err := db.NewPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to Postgres!")

	router := gin.Default()

	router.POST("/players", handlers.CreatePlayerHandler(pool))

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
