package poc

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

type Variant struct {
	MaxTimesThroughDeck int32
}
