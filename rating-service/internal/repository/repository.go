package repository

import (
	"context"
	"errors"
	"math"

	"github.com/ErenKarakus1/Matchmaking-System/rating-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRatingAlreadyExists = errors.New("rating already exists")
var ErrRatingNotFound = errors.New("rating not found")
var ErrMatchAlreadyExists = errors.New("match already exists")
var ErrMatchNotFound = errors.New("match not found")

const insertRatingSQL = `
	INSERT INTO ratings (
		player_id
	)
	VALUES ($1)
	RETURNING
		player_id,
		rating,
		games_played,
		wins,
		losses,
		created_at,
		updated_at
`

const getRatingByPlayerIDSQL = `
	SELECT
		player_id,
		rating,
		games_played,
		wins,
		losses,
		created_at,
		updated_at
	FROM ratings
	WHERE player_id=$1
`

const getRatingForUpdateSQL = `
	SELECT
		player_id,
		rating,
		games_played,
		wins,
		losses,
		created_at,
		updated_at
	FROM ratings
	WHERE player_id=$1
	FOR UPDATE
`

const updateRatingSQL = `
	UPDATE ratings
	SET 
		rating=$1,
		games_played=games_played+1,
		wins=$2,
		losses=$3,
		updated_at=NOW()
	WHERE player_id=$4

`

const insertMatchSQL = `
	INSERT INTO matches (
		id,
		winner_id,
		loser_id,
		winner_rating_before,
		loser_rating_before,
		winner_rating_after,
		loser_rating_after
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	RETURNING
		id,
		winner_id,
		loser_id,
		winner_rating_before,
		loser_rating_before,
		winner_rating_after,
		loser_rating_after,
		created_at
`

const getMatchSQL = `
	SELECT
		id,
		winner_id,
		loser_id,
		winner_rating_before,
		loser_rating_before,
		winner_rating_after,
		loser_rating_after,
		created_at
	FROM matches
	WHERE id=$1
`

const kFactor = 32

func expectedScore(playerRating, opponentRating int) float64 {
	return 1 / (1 + math.Pow(10, float64(opponentRating-playerRating)/400))
}

func calculateNewRating(playerRating, opponentRating int, actualScore float64) int {
	expected := expectedScore(playerRating, opponentRating)
	return playerRating + int(math.Round(kFactor*(actualScore-expected)))
}

func calculateRating(winnerRating, loserRating int) (int, int) {
	newWinnerRating := calculateNewRating(winnerRating, loserRating, 1)
	newLoserRating := calculateNewRating(loserRating, winnerRating, 0)
	return newWinnerRating, newLoserRating
}

func lockRatingForUpdate(ctx context.Context, tx pgx.Tx, playerID uuid.UUID) (models.Rating, error) {
	var rating models.Rating
	err := tx.QueryRow(
		ctx,
		getRatingForUpdateSQL,
		playerID,
	).Scan(
		&rating.PlayerID,
		&rating.Rating,
		&rating.GamesPlayed,
		&rating.Wins,
		&rating.Losses,
		&rating.CreatedAt,
		&rating.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Rating{}, ErrRatingNotFound
		}
		return models.Rating{}, errors.New("internal server error")
	}
	return rating, nil
}

func getRatingsForUpdate(ctx context.Context, tx pgx.Tx, winnerID uuid.UUID, loserID uuid.UUID) (models.Rating, models.Rating, error) {
	firstID := winnerID
	secondID := loserID
	changed := false
	if secondID.String() < firstID.String() {
		firstID, secondID = secondID, firstID
		changed = true
	}

	firstRating, err := lockRatingForUpdate(ctx, tx, firstID)
	if err != nil {
		if errors.Is(err, ErrRatingNotFound) {
			return models.Rating{}, models.Rating{}, err
		}
		return models.Rating{}, models.Rating{}, errors.New("internal server error")
	}

	secondRating, err := lockRatingForUpdate(ctx, tx, secondID)
	if err != nil {
		if errors.Is(err, ErrRatingNotFound) {
			return models.Rating{}, models.Rating{}, err
		}
		return models.Rating{}, models.Rating{}, errors.New("internal server error")
	}

	if changed {
		return secondRating, firstRating, nil
	}
	return firstRating, secondRating, nil
}

func InsertRating(ctx context.Context, pool *pgxpool.Pool, playerID uuid.UUID) (models.Rating, error) {
	var rating models.Rating
	err := pool.QueryRow(
		ctx,
		insertRatingSQL,
		playerID,
	).Scan(
		&rating.PlayerID,
		&rating.Rating,
		&rating.GamesPlayed,
		&rating.Wins,
		&rating.Losses,
		&rating.CreatedAt,
		&rating.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.Rating{}, ErrRatingAlreadyExists
		}
		return models.Rating{}, errors.New("internal server error")
	}
	return rating, nil
}

func GetRatingByPlayerID(ctx context.Context, pool *pgxpool.Pool, playerID uuid.UUID) (models.Rating, error) {
	var rating models.Rating
	err := pool.QueryRow(
		ctx,
		getRatingByPlayerIDSQL,
		playerID,
	).Scan(
		&rating.PlayerID,
		&rating.Rating,
		&rating.GamesPlayed,
		&rating.Wins,
		&rating.Losses,
		&rating.CreatedAt,
		&rating.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Rating{}, ErrRatingNotFound
		}
		return models.Rating{}, errors.New("internal server error")
	}
	return rating, nil
}

func SubmitMatchResult(ctx context.Context, pool *pgxpool.Pool, matchID uuid.UUID, winnerID uuid.UUID, loserID uuid.UUID) (models.Match, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return models.Match{}, errors.New("internal server error")
	}
	defer tx.Rollback(ctx)

	winnerRating, loserRating, err := getRatingsForUpdate(ctx, tx, winnerID, loserID)
	if err != nil {
		if errors.Is(err, ErrRatingNotFound) {
			return models.Match{}, err
		}
		return models.Match{}, errors.New("internal server error")
	}

	winnerAfter, loserAfter := calculateRating(winnerRating.Rating, loserRating.Rating)

	_, err = tx.Exec(
		ctx,
		updateRatingSQL,
		winnerAfter,
		winnerRating.Wins+1,
		winnerRating.Losses,
		winnerID,
	)
	if err != nil {
		return models.Match{}, errors.New("internal server error")
	}

	_, err = tx.Exec(
		ctx,
		updateRatingSQL,
		loserAfter,
		loserRating.Wins,
		loserRating.Losses+1,
		loserID,
	)
	if err != nil {
		return models.Match{}, errors.New("internal server error")
	}

	var match models.Match
	err = tx.QueryRow(
		ctx,
		insertMatchSQL,
		matchID,
		winnerID,
		loserID,
		winnerRating.Rating,
		loserRating.Rating,
		winnerAfter,
		loserAfter,
	).Scan(
		&match.ID,
		&match.WinnerID,
		&match.LoserID,
		&match.WinnerRatingBefore,
		&match.LoserRatingBefore,
		&match.WinnerRatingAfter,
		&match.LoserRatingAfter,
		&match.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.Match{}, ErrMatchAlreadyExists
		}
		return models.Match{}, errors.New("internal server error")
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Match{}, errors.New("internal server error")
	}

	return match, nil
}

func GetMatch(ctx context.Context, pool *pgxpool.Pool, matchID uuid.UUID) (models.Match, error) {
	var match models.Match
	err := pool.QueryRow(
		ctx,
		getMatchSQL,
		matchID,
	).Scan(
		&match.ID,
		&match.WinnerID,
		&match.LoserID,
		&match.WinnerRatingBefore,
		&match.LoserRatingBefore,
		&match.WinnerRatingAfter,
		&match.LoserRatingAfter,
		&match.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Match{}, ErrMatchNotFound
		}
		return models.Match{}, errors.New("internal server error")
	}
	return match, nil
}
