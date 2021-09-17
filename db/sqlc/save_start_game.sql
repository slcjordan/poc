-- Start a game.
-- name: SaveStartGame :one
WITH inserted_game AS (
  INSERT INTO game (score, max_times_through_deck)
  VALUES (@score, @max_times_through_deck)
  RETURNING id, score, max_times_through_deck
)
INSERT INTO pile_card (
  pile_num,
  pile_index,
  suit,
  index,
  position,
  game_id
)
SELECT UNNEST(@pile_nums::smallint[]) AS pile_num,
  UNNEST(@pile_indexes::smallint[]) AS pile_index,
  UNNEST(@suits::smallint[]) AS suit,
  UNNEST(@indexes::smallint[]) AS index,
  UNNEST(@positions::integer[]) AS position,
  inserted_game.id AS game_id
FROM inserted_game
RETURNING id AS game_id;
