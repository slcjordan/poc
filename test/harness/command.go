// Code generated by cmd/harness; DO NOT EDIT.

package harness

import (
	"context"
	"testing"

	"github.com/slcjordan/poc"
)

type ErrorChecker interface {
	CheckError(*testing.T, string, error)
}

type StartGameChecker interface {
	ErrorChecker
	CheckStartGame(*testing.T, string, poc.StartGame)
}

type StartGame []struct {
	Desc    string
	Input   poc.StartGame
	Command poc.StartGameCaller
	Result  StartGameChecker
}

func (h StartGame) Run(t *testing.T) {
	for _, testCase := range h {
		t.Run(testCase.Desc, func(t *testing.T) {
			result, err := testCase.Command.CallStartGame(context.Background(), testCase.Input)
			if testCase.Result != nil {
				testCase.Result.CheckError(t, "", err)
				testCase.Result.CheckStartGame(t, "", result)
			}
		})
	}
}

type PerformMoveChecker interface {
	ErrorChecker
	CheckPerformMove(*testing.T, string, poc.PerformMove)
}

type PerformMove []struct {
	Desc    string
	Input   poc.PerformMove
	Command poc.PerformMoveCaller
	Result  PerformMoveChecker
}

func (h PerformMove) Run(t *testing.T) {
	for _, testCase := range h {
		t.Run(testCase.Desc, func(t *testing.T) {
			result, err := testCase.Command.CallPerformMove(context.Background(), testCase.Input)
			if testCase.Result != nil {
				testCase.Result.CheckError(t, "", err)
				testCase.Result.CheckPerformMove(t, "", result)
			}
		})
	}
}

type ListGamesChecker interface {
	ErrorChecker
	CheckListGames(*testing.T, string, poc.ListGames)
}

type ListGames []struct {
	Desc    string
	Input   poc.ListGames
	Command poc.ListGamesCaller
	Result  ListGamesChecker
}

func (h ListGames) Run(t *testing.T) {
	for _, testCase := range h {
		t.Run(testCase.Desc, func(t *testing.T) {
			result, err := testCase.Command.CallListGames(context.Background(), testCase.Input)
			if testCase.Result != nil {
				testCase.Result.CheckError(t, "", err)
				testCase.Result.CheckListGames(t, "", result)
			}
		})
	}
}
