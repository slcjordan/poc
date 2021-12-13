package pipeline

import (
	"context"
	"fmt"

	"github.com/slcjordan/poc"
)

// PerformMove uses the same context for every command, but uses the
// performMove output from the previous command as input to the next command.
type PerformMove []poc.PerformMoveCaller

// CallPerformMove exits early at the first command that returns an error.
func (ppipe PerformMove) CallPerformMove(ctx context.Context, p poc.PerformMove) (poc.PerformMove, error) {
	var err error

	for i, step := range ppipe {
		fmt.Printf("step %d\n", i)
		p, err = step.CallPerformMove(ctx, p)
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

type PerformMoveMiddleware interface {
	PerformMoveUse(poc.PerformMoveCaller) poc.PerformMoveCaller
}

// Use middleware to wrap each command.
func (ppipe PerformMove) UseEach(middleware ...PerformMoveMiddleware) PerformMove {
	result := make([]poc.PerformMoveCaller, len(ppipe))
	for _, step := range ppipe {
		for _, mw := range middleware {
			result = append(result, mw.PerformMoveUse(step))
		}
	}
	return result
}

// StartGame uses the same context for every command, but uses the
// startGame output from the previous command as input to the next command.
type StartGame []poc.StartGameCaller

// CallStartGame exits early at the first command that returns an error.
func (spipe StartGame) CallStartGame(ctx context.Context, s poc.StartGame) (poc.StartGame, error) {
	var err error

	for _, step := range spipe {
		s, err = step.CallStartGame(ctx, s)
		if err != nil {
			return s, err
		}
	}
	return s, nil
}

type StartGameMiddleware interface {
	StartGameUse(poc.StartGameCaller) poc.StartGameCaller
}

// Use middleware to wrap each command.
func (spipe StartGame) UseEach(middleware ...StartGameMiddleware) StartGame {
	result := make([]poc.StartGameCaller, len(spipe))
	for _, step := range spipe {
		for _, mw := range middleware {
			result = append(result, mw.StartGameUse(step))
		}
	}
	return result
}

// ListGames uses the same context for every command, but uses the
// listGame output from the previous command as input to the next command.
type ListGames []poc.ListGamesCaller

// CallListGames exits early at the first command that returns an error.
func (lpipe ListGames) CallListGames(ctx context.Context, l poc.ListGames) (poc.ListGames, error) {
	var err error

	for _, step := range lpipe {
		l, err = step.CallListGames(ctx, l)
		if err != nil {
			return l, err
		}
	}
	return l, nil
}

type ListGamesMiddleware interface {
	ListGamesUse(poc.ListGamesCaller) poc.ListGamesCaller
}

// Use middleware to wrap each command.
func (lpipe ListGames) UseEach(middleware ...ListGamesMiddleware) ListGames {
	result := make([]poc.ListGamesCaller, len(lpipe))
	for _, step := range lpipe {
		for _, mw := range middleware {
			result = append(result, mw.ListGamesUse(step))
		}
	}
	return result
}

