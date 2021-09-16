package pipeline

import (
	"context"

	"github.com/slcjordan/poc"
)

type PerformMove []poc.PerformMoveCaller

func (ppipe PerformMove) CallPerformMove(ctx context.Context, p poc.PerformMove) (poc.PerformMove, error) {

	for _, step := range ppipe {
		p, err := step.CallPerformMove(ctx, p)
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

type StartGame []poc.StartGameCaller

func (spipe StartGame) CallStartGame(ctx context.Context, s poc.StartGame) (poc.StartGame, error) {
	for _, step := range spipe {
		s, err := step.CallStartGame(ctx, s)
		if err != nil {
			return s, err
		}
	}
	return s, nil
}

type ListGames []poc.ListGamesCaller

func (lpipe ListGames) CallListGames(ctx context.Context, s poc.ListGames) (poc.ListGames, error) {
	for _, step := range lpipe {
		s, err := step.CallListGames(ctx, s)
		if err != nil {
			return s, err
		}
	}
	return s, nil
}
