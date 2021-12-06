package poc

//go:generate build/assert -filename=service.go

// SavedGameSummary is a saved game with summary of the game state.
type SavedGameSummary struct {
	GameID int64
	Score  int32
}

// SavedGameDetail is a saved game with detail of the game state.
type SavedGameDetail struct {
	GameID            int64
	Board             Board
	History           [][]Move
	PossibleNextMoves [][]Move
	Variant           Variant
}

// StartGame starts a game.
type StartGame struct {
	Variant         Variant
	SavedGameDetail SavedGameDetail
}

// PerformMove executes a move on a game.
type PerformMove struct {
	Next            []Move
	SavedGameDetail SavedGameDetail
}

// ListGames lists running games.
type ListGames struct {
	Cursor struct {
		Offset int32
		Limit  int32
	}
	Games []SavedGameSummary
}
