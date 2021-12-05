package poc

//go:generate build/assert -filename=model.go

// Suit is a the pip part of the card.
//go:generate stringer -type=Suit
type Suit uint8

// Standard deck suits.
const (
	_ Suit = iota
	Hearts
	Clubs
	Diamonds
	Spades
	Joker
)

// Index is the value part of the card.
//go:generate stringer -type=Index
type Index uint8

// Standard deck indices.
const (
	_ Index = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Juggler
	Fool
)

// Card is a face card.
type Card struct {
	Suit  Suit
	Index Index
}

// Position is a way of positioning a card.
type Position uint64

// Possible ways of positioning a card.
const (
	FaceUp Position = 1 << iota
)

// PositionedCard is a card that has a position
type PositionedCard struct {
	Position Position
	Card     Card
}

// Board is the current state of the board.
type Board struct {
	Piles [13][]PositionedCard
	Score int32
}

// Variant is a rules variant.
type Variant struct {
	MaxTimesThroughDeck int32
}

// Move is a transformation of the board.
type Move struct {
	OldPileNum      int
	OldPileIndex    int
	OldPilePosition Position
	NewPileNum      int
	NewPileIndex    int
	NewPilePosition Position
}

/* Service objects */

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
