package grpc

import (
	"context"
	"errors"

	"github.com/ErenKarakus1/Matchmaking-System/rating-service/internal/repository"
	ratingv1 "github.com/ErenKarakus1/Matchmaking-System/rating-service/proto/rating/v1"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RatingServer struct {
	ratingv1.UnimplementedRatingServiceServer
	pool *pgxpool.Pool
}

func NewRatingServer(pool *pgxpool.Pool) *RatingServer {
	return &RatingServer{pool: pool}
}

func (s *RatingServer) GetPlayerRating(ctx context.Context, req *ratingv1.GetPlayerRatingRequest) (*ratingv1.GetPlayerRatingResponse, error) {
	playerID, err := uuid.Parse(req.PlayerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid player id")
	}

	rating, err := repository.GetRatingByPlayerID(ctx, s.pool, playerID)
	if err != nil {
		if errors.Is(err, repository.ErrRatingNotFound) {
			return nil, status.Error(codes.NotFound, "rating not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ratingv1.GetPlayerRatingResponse{
		PlayerId:    rating.PlayerID.String(),
		Rating:      int32(rating.Rating),
		GamesPlayed: int32(rating.GamesPlayed),
		Wins:        int32(rating.Wins),
		Losses:      int32(rating.Losses),
	}, nil
}
