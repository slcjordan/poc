package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/slcjordan/poc"
	"github.com/slcjordan/poc/db/sqlc"
	"github.com/slcjordan/poc/logger"
)

// PgxPoolIface is a pool of postgres connections.
type PgxPoolIface interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginFunc(ctx context.Context, f func(pgx.Tx) error) (err error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Close()
}

// Save commands create or update data.
type Save struct {
	Pool PgxPoolIface
}

// CallStartGame saves start.Result as a new game.
func (s *Save) CallStartGame(ctx context.Context, start poc.StartGame) (poc.StartGame, error) {
	conn, err := s.Pool.Acquire(ctx)
	if err != nil {
		logger.Infof(ctx, "could not acquire connection: %s", err)
		return start, poc.Error{Actual: errors.New("db unavailable"), Category: poc.UnavailableError}
	}
	var pileNums []int16
	var pileIndexes []int16
	var suits []int16
	var indexes []int16
	var positions []int32

	for pileNum, curr := range start.Result.Board.Piles {
		for pileIndex, card := range curr {
			pileNums = append(pileNums, int16(pileNum))
			pileIndexes = append(pileIndexes, int16(pileIndex))
			suits = append(suits, int16(card.Card.Suit))
			indexes = append(indexes, int16(card.Card.Index))
			positions = append(positions, int32(card.Position))
		}
	}
	logger.Infof(ctx, "pileNums: %v, pileIndexes: %v suits: %v, indexes: %v, positions: %v", pileNums, pileIndexes, suits, indexes, positions)
	gameID, err := sqlc.New(conn).SaveStartGame(ctx, sqlc.SaveStartGameParams{
		Score:               start.Result.Board.Score,
		PileNums:            pileNums,
		PileIndexes:         pileIndexes,
		Suits:               suits,
		Indexes:             indexes,
		Positions:           positions,
		MaxTimesThroughDeck: start.Result.Variant.MaxTimesThroughDeck,
	})
	if err != nil {
		logger.Errorf(ctx, "could not save game: %s", err)
		return start, poc.Error{Actual: errors.New("could not save game"), Category: poc.UnknownError}
	}
	start.Result.GameID = gameID
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
