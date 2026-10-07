package main

import (
	"log"
	"net"

	"github.com/ErenKarakus1/Matchmaking-System/rating-service/internal/config"
	"github.com/ErenKarakus1/Matchmaking-System/rating-service/internal/db"
	ratinggrpc "github.com/ErenKarakus1/Matchmaking-System/rating-service/internal/grpc"
	"github.com/ErenKarakus1/Matchmaking-System/rating-service/internal/handlers"
	ratingv1 "github.com/ErenKarakus1/Matchmaking-System/rating-service/proto/rating/v1"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()
	pool, err := db.NewPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to Postgres")

	grpcServer := grpc.NewServer()
	ratingv1.RegisterRatingServiceServer(grpcServer, ratinggrpc.NewRatingServer(pool))

	listener, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Println("Rating gRPC server is running on :9091")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	router := gin.Default()

	router.POST("/ratings/players/:player_id", handlers.CreatePlayerRatingHandler(pool))
	router.GET("/ratings/players/:player_id", handlers.GetPlayerRatingHandler(pool))
	router.POST("/ratings/matches/:match_id/result", handlers.CreateMatchHandler(pool))
	router.GET("/ratings/matches/:match_id", handlers.GetMatchHandler(pool))
	router.GET("/ratings/leaderboard", handlers.GetLeaderboardHandler(pool))
	router.GET("/ratings/players/:player_id/matches", handlers.GetPlayerMatchesHandler(pool))

	if err := router.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
