package ratingclient

import (
	"context"

	ratingv1 "github.com/ErenKarakus1/Matchmaking-System/matchmaking-service/proto/rating/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Client struct {
	client ratingv1.RatingServiceClient
}

func New(conn grpc.ClientConnInterface) *Client {
	return &Client{
		client: ratingv1.NewRatingServiceClient(conn),
	}
}

func (c *Client) GetPlayerRating(ctx context.Context, playerID uuid.UUID) (int, error) {
	resp, err := c.client.GetPlayerRating(ctx, &ratingv1.GetPlayerRatingRequest{
		PlayerId: playerID.String(),
	})
	if err != nil {
		return 0, err
	}

	return int(resp.Rating), nil
}
