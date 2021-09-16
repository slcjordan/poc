package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/slcjordan/poc"
	"github.com/slcjordan/poc/db/sqlc"
	"github.com/slcjordan/poc/logger"
)

// Save commands create or update data.
type Save struct {
	Pool *pgxpool.Pool
}

// CallStartGame saves start.Result as a new game.
func (s *Save) CallStartGame(ctx context.Context, start poc.StartGame) (poc.StartGame, error) {
	conn, err := s.Pool.Acquire(ctx)
	if err != nil {
		logger.Infof(ctx, "could not acquire connection: %s", err)
		return start, poc.Error{Actual: errors.New("db unavailable"), Category: poc.UnavailableError}
	}
	row, err := sqlc.New(conn).SaveStartGame(ctx, sqlc.SaveStartGameParams{
		Score:               start.Result.Board.Score,
		MaxTimesThroughDeck: start.Result.Variant.MaxTimesThroughDeck,
	})
	if err != nil {
		logger.Errorf(ctx, "could not save game: %s", err)
		return start, poc.Error{Actual: errors.New("could not save game"), Category: poc.UnknownError}
	}
	start.Result.GameID = row.ID
	start.Result.Board.Score = row.Score
	start.Result.Variant.MaxTimesThroughDeck = row.MaxTimesThroughDeck
	return start, nil
}

// CallPerformMove updates start.Result as an existing game and expects move.Input.GameID to be set.
func (s *Save) CallPerformMove(ctx context.Context, move poc.PerformMove) (poc.PerformMove, error) {
	conn, err := s.Pool.Acquire(ctx)
	if err != nil {
		logger.Infof(ctx, "could not acquire connection: %s", err)
		return move, poc.Error{Actual: errors.New("db unavailable"), Category: poc.UnavailableError}
	}
	oldPileNums := make([]int, len(move.Input.Move))
	oldPileIndexes := make([]int, len(move.Input.Move))
	oldPilePositions := make([]uint64, len(move.Input.Move))
	newPileNums := make([]int, len(move.Input.Move))
	newPileIndexes := make([]int, len(move.Input.Move))
	newPilePositions := make([]uint64, len(move.Input.Move))

	for i, curr := range move.Input.Move {
		oldPileNums[i] = curr.OldPileNum
		oldPileIndexes[i] = curr.OldPileIndex
		oldPilePositions[i] = uint64(curr.OldPilePosition)
		newPileNums[i] = curr.NewPileNum
		newPileIndexes[i] = curr.NewPileIndex
		newPilePositions[i] = uint64(curr.NewPilePosition)
	}
	_, err = sqlc.New(conn).SavePerformMove(ctx, sqlc.SavePerformMoveParams{
		GameID:           move.Input.GameID,
		OldPileNums:      oldPileNums,
		OldPileIndexes:   oldPileIndexes,
		OldPilePositions: oldPilePositions,
		NewPileNums:      newPileNums,
		NewPileIndexes:   newPileIndexes,
		NewPilePositions: newPilePositions,
	})
	if err != nil {
		logger.Errorf(ctx, "could not save game %d: %s", move.Input.GameID, err)
		return move, poc.Error{Actual: errors.New("could not save game"), Category: poc.UnknownError}
	}
	return move, nil
}
