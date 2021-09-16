-- Start a game.
-- name: SaveStartGame :one

INSERT INTO game (score, max_times_through_deck)
VALUES ($1, $2)
RETURNING id, score, max_times_through_deck;
