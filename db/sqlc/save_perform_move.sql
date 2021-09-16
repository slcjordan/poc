-- Make a move.
-- name: SavePerformMove :many

WITH last_move AS (
  SELECT COALESCE(max(move_number), 0) last_move_number
  FROM history WHERE game_id = @game_id
), inserted_moves AS (
  INSERT INTO move (
    old_pile_num,
    old_pile_index,
    old_pile_position,
    new_pile_num,
    new_pile_index,
    new_pile_position
  )
  SELECT UNNEST(@old_pile_nums) AS old_pile_num,
    UNNEST(@old_pile_indexes) AS old_pile_index,
    UNNEST(@old_pile_positions) AS old_pile_position,
    UNNEST(@new_pile_nums) AS new_pile_num,
    UNNEST(@new_pile_indexes) AS new_pile_index,
    UNNEST(@new_pile_positions) AS new_pile_position
  RETURNING id AS move_id
)
INSERT INTO history (game_id, move_id, move_number)
SELECT @game_id, last_move_number, move_id
FROM last_move, inserted_moves
RETURNING game_id, move_id;
