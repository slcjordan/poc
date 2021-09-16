-- Search returns a list of games.
-- name: SearchGame :many

SELECT id, score
FROM game
ORDER BY id LIMIT $1 OFFSET $2;
