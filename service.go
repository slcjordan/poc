package poc

import "context"

// Move is a transformation of the board.
type Move struct {
	OldPileNum      int
	OldPileIndex    int
	OldPilePosition Position
	NewPileNum      int
	NewPileIndex    int
	NewPilePosition Position
}

// History is a record of moves.
type History [][]Move

// SavedGameSummary is a saved game with summary of the game state.
type SavedGameSummary struct {
	GameID int64
	Score  int32
}

// SavedGameDetail is a saved game with detail of the game state.
type SavedGameDetail struct {
	GameID            int64
	Board             Board
	History           History
	PossibleNextMoves [][]Move
	Variant           Variant
}

type StartGameCaller interface {
	CallStartGame(context.Context, StartGame) (StartGame, error)
}

type PerformMoveCaller interface {
	CallPerformMove(context.Context, PerformMove) (PerformMove, error)
}

type ListGamesCaller interface {
	CallListGames(context.Context, ListGames) (ListGames, error)
}

// StartGame starts a game.
type StartGame struct {
	Input  Variant
	Result SavedGameDetail
}

// PerformMove executes a move on a game.
type PerformMove struct {
	Input struct {
		GameID int64
		Move   []Move
	}
	Result SavedGameDetail
}

// ListGames lists running games.
type ListGames struct {
	Input struct {
		Offset int32
		Limit  int32
	}
	Result []SavedGameSummary
}
