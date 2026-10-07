package ratingclient

import (
	"context"
	"errors"
	"testing"

	ratingv1 "github.com/ErenKarakus1/Matchmaking-System/matchmaking-service/proto/rating/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type fakeConn struct {
	method string
	req    *ratingv1.GetPlayerRatingRequest
	resp   *ratingv1.GetPlayerRatingResponse
	err    error
}

func (f *fakeConn) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	f.method = method

	req, ok := args.(*ratingv1.GetPlayerRatingRequest)
	if !ok {
		return errors.New("unexpected request type")
	}
	f.req = req

	if f.err != nil {
		return f.err
	}

	resp, ok := reply.(*ratingv1.GetPlayerRatingResponse)
	if !ok {
		return errors.New("unexpected response type")
	}
	*resp = *f.resp
	return nil
}

func (f *fakeConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, errors.New("streaming is not supported")
}

func TestGetPlayerRating(t *testing.T) {
	playerID := uuid.New()
	conn := &fakeConn{
		resp: &ratingv1.GetPlayerRatingResponse{
			PlayerId: playerID.String(),
			Rating:   1532,
		},
	}
	client := New(conn)

	rating, err := client.GetPlayerRating(context.Background(), playerID)
	if err != nil {
		t.Fatalf("GetPlayerRating returned error: %v", err)
	}

	if rating != 1532 {
		t.Fatalf("rating = %d, want 1532", rating)
	}
	if conn.method != ratingv1.RatingService_GetPlayerRating_FullMethodName {
		t.Fatalf("method = %q, want %q", conn.method, ratingv1.RatingService_GetPlayerRating_FullMethodName)
	}
	if conn.req == nil {
		t.Fatal("expected request to be sent")
	}
	if conn.req.PlayerId != playerID.String() {
		t.Fatalf("PlayerId = %q, want %q", conn.req.PlayerId, playerID.String())
	}
}

func TestGetPlayerRatingReturnsError(t *testing.T) {
	wantErr := errors.New("rating service unavailable")
	client := New(&fakeConn{err: wantErr})

	_, err := client.GetPlayerRating(context.Background(), uuid.New())
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}
