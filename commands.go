package poc

import "context"

//go:generate go run github.com/golang/mock/mockgen -package=mocks -destination=test/mocks/commands.go -source=commands.go

// StartGame starts a game.
type StartGame struct {
	Variant         Variant
	SavedGameDetail SavedGameDetail
}

// StartGameCaller is a start game command.
type StartGameCaller interface {
	CallStartGame(context.Context, StartGame) (StartGame, error)
}

// PerformMove executes a move on a game.
type PerformMove struct {
	Next            []Move
	SavedGameDetail SavedGameDetail
}

// PerformMoveCaller is a perform move command.
type PerformMoveCaller interface {
	CallPerformMove(context.Context, PerformMove) (PerformMove, error)
}

// ListGames lists running games.
type ListGames struct {
	Cursor struct {
		Offset int32
		Limit  int32
	}
	Games []SavedGameSummary
}

// ListGamesCaller is a list game command.
type ListGamesCaller interface {
	CallListGames(context.Context, ListGames) (ListGames, error)
}
