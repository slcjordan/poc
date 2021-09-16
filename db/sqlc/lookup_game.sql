-- Lookup a game.
-- name: LookupGameDetail :one

SELECT game.id, score, max_times_through_deck, array_agg(piles) piles, array_agg(hist_moves) history
FROM game JOIN LATERAL (
  SELECT array_agg((suit, index, position) ORDER BY pile_index) piles
  FROM pile_card WHERE game_id = game.id
  GROUP BY pile_num ORDER BY pile_num
) p ON TRUE JOIN LATERAL(
  select array_agg(
    (old_pile_num, old_pile_index, old_pile_position, new_pile_num, new_pile_index, new_pile_position)
    ORDER BY old_pile_num, old_pile_index) hist_moves
  FROM history
  JOIN move ON move.id = history.move_id
  WHERE history.game_id = game.id
  GROUP BY move_number
) m ON TRUE WHERE game.id = @game_id GROUP BY game.id;
