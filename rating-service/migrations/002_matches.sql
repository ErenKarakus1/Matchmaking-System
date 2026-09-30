CREATE TABLE IF NOT EXISTS matches (
    id UUID PRIMARY KEY,
    winner_id UUID NOT NULL,
    loser_id UUID NOT NULL,
    winner_rating_before INT NOT NULL,
    loser_rating_before INT NOT NULL,
    winner_rating_after INT NOT NULL,
    loser_rating_after INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (winner_id <> loser_id)
);