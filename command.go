package poc

import "context"

//go:generate go run github.com/golang/mock/mockgen -package=mocks -destination=test/mocks/command.go -source=command.go

// StartGameCaller is a start game command.
type StartGameCaller interface {
	CallStartGame(context.Context, StartGame) (StartGame, error)
}

// PerformMoveCaller is a perform move command.
type PerformMoveCaller interface {
	CallPerformMove(context.Context, PerformMove) (PerformMove, error)
}

// ListGamesCaller is a list game command.
type ListGamesCaller interface {
	CallListGames(context.Context, ListGames) (ListGames, error)
}
