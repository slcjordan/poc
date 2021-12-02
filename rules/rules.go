package rules

import (
	"context"
	"errors"
	"math/rand"
	"sort"

	"github.com/slcjordan/poc"
)

// ErrInvalidMove means the user tried a bad move.
var ErrInvalidMove = errors.New("invalid move")

func nextMoves(
	stock []poc.PositionedCard,
	talon []poc.PositionedCard,
	tableau [][]poc.PositionedCard,
	foundation [][]poc.PositionedCard) [][]poc.Move {

	var result [][]poc.Move
	var suits [4]poc.Index // top cards on the foundation
	for i, pile := range foundation {
		if len(pile) < 1 {
			continue
		}
		suits[i] = pile[len(pile)-1].Card.Index
	}

	// indexes of top tableau cards as move destinations.
	// e.g. if we can place a part of a tableau pile onto any table of type red
	// jack, we will have already saved possible destination piles in red[10]
	// (Aces start at index 0).
	var red [13][]int
	var black [13][]int

	for pileNum := range tableau {
		length := len(tableau[pileNum])
		if length < 1 {
			continue
		}
		card := tableau[pileNum][length-1]
		index := int(card.Card.Index - poc.Ace)
		if card.Position&poc.FaceUp == 0 { // flip over top face-down card on a tableau pile
			result = append(result, []poc.Move{{
				OldPileNum:      pileNum + 2,
				OldPileIndex:    length - 1,
				OldPilePosition: card.Position,
				NewPileNum:      pileNum + 2,
				NewPileIndex:    length - 1,
				NewPilePosition: card.Position | poc.FaceUp,
			}})
			continue
		}
		suit := card.Card.Suit
		switch suit { // index the top card as a destination pile
		case poc.Hearts, poc.Diamonds:
			red[index] = append(red[index], pileNum)
		case poc.Spades, poc.Clubs:
			black[index] = append(black[index], pileNum)
		}
		if suits[card.Card.Suit-1] == card.Card.Index-1 { // move top face-up card onto a foundation pile
			result = append(result, []poc.Move{{
				OldPileNum:      int(card.Card.Suit + 7),
				OldPileIndex:    length - 1,
				OldPilePosition: card.Position,
				NewPileNum:      int(card.Card.Suit + 7),
				NewPileIndex:    len(foundation[card.Card.Suit-1]),
				NewPilePosition: card.Position,
			}})
		}
	}

	for pileNum := range tableau {
		for i, card := range tableau[pileNum] {
			if card.Position&poc.FaceUp == 0 {
				continue
			}
			suit := card.Card.Suit
			index := int(card.Card.Index - poc.Ace - 1)
			if index < 1 { // Aces can't be placed on top of other piles
				continue
			}
			var dest []int
			switch suit {
			case poc.Hearts, poc.Diamonds:
				dest = black[index]
			case poc.Spades, poc.Clubs:
				dest = red[index]
			}

			for _, curr := range dest { // move part of tableau pile onto another tableau pile
				var currMove []poc.Move
				for idx := i; idx < len(tableau[pileNum]); idx++ {
					currMove = append(currMove, poc.Move{
						OldPileNum:      pileNum + 2,
						OldPileIndex:    idx,
						OldPilePosition: card.Position,
						NewPileNum:      curr + 2,
						NewPileIndex:    len(tableau[curr]) + (idx - i),
						NewPilePosition: card.Position | poc.FaceUp,
					})
				}
				result = append(result, currMove)
			}
		}
	}

	if len(stock) > 0 { // draw a card from the stock.
		card := stock[len(stock)-1]
		result = append(result, []poc.Move{{
			OldPileNum:      0,
			OldPileIndex:    len(stock) - 1,
			OldPilePosition: card.Position,
			NewPileNum:      1,
			NewPileIndex:    len(talon),
			NewPilePosition: card.Position | poc.FaceUp,
		}})
	} else { // return talon to the stock
		length := len(talon)
		var currMove []poc.Move

		for i := 0; i < length; i++ {
			currMove = append(currMove, poc.Move{
				OldPileNum:      1,
				OldPileIndex:    length - 1 - i,
				OldPilePosition: talon[len(talon)-1].Position,
				NewPileNum:      0,
				NewPileIndex:    i,
				NewPilePosition: talon[len(talon)-1].Position & ^poc.FaceUp,
			})
		}
		result = append(result, currMove)
	}
	return result
}

type sortable []poc.Move

func (s *sortable) Len() int {
	return len(*s)
}

func (s *sortable) Less(i int, j int) bool {
	a := []poc.Move(*s)[i]
	b := []poc.Move(*s)[j]

	if a.OldPileNum < b.OldPileNum {
		return true
	}
	if a.OldPileNum > b.OldPileNum {
		return false
	}
	if a.OldPileIndex < b.OldPileIndex {
		return true
	}
	if a.OldPileIndex > b.OldPileIndex {
		return false
	}
	if a.OldPilePosition < b.OldPilePosition {
		return true
	}
	if a.OldPilePosition > b.OldPilePosition {
		return false
	}
	if a.NewPileNum < b.NewPileNum {
		return true
	}
	if a.NewPileNum > b.NewPileNum {
		return false
	}
	if a.NewPileIndex < b.NewPileIndex {
		return true
	}
	if a.NewPileIndex > b.NewPileIndex {
		return false
	}
	if a.NewPilePosition < b.NewPilePosition {
		return true
	}
	if a.NewPilePosition > b.NewPilePosition {
		return false
	}
	return true
}

func (s *sortable) Compare(other *sortable) int {
	if len(*s) < len(*other) {
		return -1
	}
	if len(*s) > len(*other) {
		return 1
	}
	if len(*s) == 0 {
		return 0
	}
	a := []poc.Move(*s)[0]
	b := []poc.Move(*other)[0]
	if a.OldPileNum < b.OldPileNum {
		return -1
	}
	if a.OldPileNum > b.OldPileNum {
		return 1
	}
	if a.OldPileIndex < b.OldPileIndex {
		return -1
	}
	if a.OldPileIndex > b.OldPileIndex {
		return 1
	}
	if a.OldPilePosition < b.OldPilePosition {
		return -1
	}
	if a.OldPilePosition > b.OldPilePosition {
		return 1
	}
	if a.NewPileNum < b.NewPileNum {
		return -1
	}
	if a.NewPileNum > b.NewPileNum {
		return 1
	}
	if a.NewPileIndex < b.NewPileIndex {
		return -1
	}
	if a.NewPileIndex > b.NewPileIndex {
		return 1
	}
	if a.NewPilePosition < b.NewPilePosition {
		return -1
	}
	if a.NewPilePosition > b.NewPilePosition {
		return 1
	}
	remainderA := sortable((*s)[1:])
	remainderB := sortable((*other)[1:])
	return (&remainderA).Compare(&remainderB)
}

func (s *sortable) Swap(i int, j int) {
	[]poc.Move(*s)[i], []poc.Move(*s)[j] = []poc.Move(*s)[j], []poc.Move(*s)[i]
}

// Validate validates command input.
type Validate struct{}

// CallPerformMove checks that the move can actually be performed.
func (v Validate) CallPerformMove(ctx context.Context, move poc.PerformMove) (poc.PerformMove, error) {
	// get list of possible moves
	possible := nextMoves(
		move.Result.Board.Piles[0],
		move.Result.Board.Piles[1],
		move.Result.Board.Piles[2:9],
		move.Result.Board.Piles[9:],
	)
	// sort everything
	for i := 0; i < len(possible); i++ {
		curr := sortable(possible[i])
		sort.Sort(&curr)
	}
	input := sortable(move.Input.Move)
	sort.Sort(&input)
	sort.Slice(possible, func(i int, j int) bool {
		a := sortable(possible[i])
		b := sortable(possible[j])
		return (&a).Compare(&b) <= 0
	})
	// perform binary search on possible moves to find the input in the list of valid moves.
	result := sort.Search(len(possible), func(i int) bool {
		a := sortable(possible[i])
		return (&a).Compare(&input) >= 0
	})
	if result >= len(possible) {
		return move, poc.Error{Actual: ErrInvalidMove, Category: poc.SemanticError}
	}
	target := sortable(possible[result])
	if (&target).Compare(&input) != 0 {
		return move, poc.Error{Actual: ErrInvalidMove, Category: poc.SemanticError}
	}
	return move, nil
}

// NextMove hydrates next move in result with the next reachable moves.
type NextMove struct{}

// CallStartGame moves.
func (n NextMove) CallStartGame(ctx context.Context, game poc.StartGame) (poc.StartGame, error) {
	var timesThroughDeck int32
outer:
	for _, currMove := range game.Result.History {
		for _, currCard := range currMove {
			if currCard.NewPileNum == 0 {
				timesThroughDeck++
				continue outer
			}
		}
	}
	if timesThroughDeck >= game.Input.MaxTimesThroughDeck && game.Input.MaxTimesThroughDeck > 0 {
		game.Result.PossibleNextMoves = nil
		return game, nil
	}
	game.Result.PossibleNextMoves = nextMoves(
		game.Result.Board.Piles[0],
		game.Result.Board.Piles[1],
		game.Result.Board.Piles[2:9],
		game.Result.Board.Piles[9:],
	)
	return game, nil
}

// CallPerformMove moves.
func (n NextMove) CallPerformMove(ctx context.Context, move poc.PerformMove) (poc.PerformMove, error) {
	move.Result.PossibleNextMoves = nextMoves(
		move.Result.Board.Piles[0],
		move.Result.Board.Piles[1],
		move.Result.Board.Piles[2:9],
		move.Result.Board.Piles[9:],
	)
	return move, nil
}

// Shuffle shuffles the deck.
type Shuffle struct{ Source rand.Source }

// CallStartGame shuffles a new deck and deals it to the result.Board piles.
func (s Shuffle) CallStartGame(ctx context.Context, game poc.StartGame) (poc.StartGame, error) {
	cards := make([]poc.PositionedCard, 52)

	for i := range cards {
		cards[i].Card.Suit = poc.Suit((i % 4) + 1)
		cards[i].Card.Index = poc.Index((i / 4) + 1)
	}
	swap := func(i int, j int) { cards[i], cards[j] = cards[j], cards[i] }
	rand.New(s.Source).Shuffle(len(cards), swap)

	game.Result.Board.Piles[8] = cards[21:28]
	game.Result.Board.Piles[7] = cards[15:21]
	game.Result.Board.Piles[6] = cards[10:15]
	game.Result.Board.Piles[5] = cards[6:10]
	game.Result.Board.Piles[4] = cards[3:6]
	game.Result.Board.Piles[3] = cards[1:3]
	game.Result.Board.Piles[2] = cards[:1]
	game.Result.Board.Piles[0] = cards[28:]
	return game, nil
}
