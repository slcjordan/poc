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

// ListGames uses the same context for every command, but uses the
// listGame output from the previous command as input to the next command.
type ListGames []poc.ListGamesCaller

// CallListGames exits early at the first command that returns an error.
func (lpipe ListGames) CallListGames(ctx context.Context, s poc.ListGames) (poc.ListGames, error) {
	var err error

	for _, step := range lpipe {
		s, err = step.CallListGames(ctx, s)
		if err != nil {
			return s, err
		}
	}
	return s, nil
}
