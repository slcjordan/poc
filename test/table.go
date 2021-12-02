package test

import (
	"context"
	"errors"
	"math"
	"testing"

	"github.com/slcjordan/poc"
)

// StartGame is a table-driven test
type StartGame []struct {
	Desc    string
	Input   poc.StartGame
	Command poc.StartGameCaller
	Error   ErrorAssertion
	Result  StartGameAssertion
}

// Run executes tests.
func (s StartGame) Run(t *testing.T) {
	for _, testCase := range s {
		t.Run(testCase.Desc, func(t *testing.T) {
			result, err := testCase.Command.CallStartGame(context.Background(), testCase.Input)
			if testCase.Error != nil {
				testCase.Error.AssertError(t, err)
			}
			if testCase.Result != nil {
				testCase.Result.AssertStartGame(t, result)
			}
		})
	}
}

// StartGameAssertion checks a StartGame result.
type StartGameAssertion interface {
	AssertStartGame(*testing.T, poc.StartGame)
}

// ErrorAssertion checks an error result.
type ErrorAssertion interface {
	AssertError(*testing.T, error)
}

// Category checks the category of an error.
type Category struct {
	Expected poc.ErrorCategory
}

// AssertError will only pass if not nil and of correct category.
func (c Category) AssertError(t *testing.T, err error) {
	t.Run("error category", func(t *testing.T) {
		var categorized poc.Error
		if !errors.As(err, &categorized) {
			t.Fatalf("Error of type %T does not wrap nor is of type poc.Error: %s", err, err)
		}
		actual := categorized.Category
		if actual != c.Expected {
			t.Fatalf("Expected error of category %s but got %s", c.Expected, actual)
		}
	})
}

// IsNil fails if error is not nil.
type IsNil struct{}

// AssertError will only pass if not nil and of correct category.
func (i IsNil) AssertError(t *testing.T, err error) {
	t.Run("error is nil", func(t *testing.T) {
		if err != nil {
			t.Fatalf("Error of type %T is not nil: %s", err, err)
		}
	})
}

// Assert makes assertions from game attributes.
type Assert struct {
	Bytes             BytesAssertion
	GameID            GameIDAssertion
	Board             BoardAssertion
	History           HistoryAssertion
	PossibleNextMoves PossibleNextMovesAssertion
	Variant           VariantAssertion
}

// AssertStartGame applies all data assertions to games fields.
func (a Assert) AssertStartGame(t *testing.T, s poc.StartGame) {
	if a.GameID != nil {
		t.Run("gameID", func(t *testing.T) { a.GameID.AssertGameID(t, s.Result.GameID) })
	}
	if a.Board != nil {
		t.Run("Board", func(t *testing.T) { a.Board.AssertBoard(t, s.Result.Board) })
	}
	if a.History != nil {
		t.Run("History", func(t *testing.T) { a.History.AssertHistory(t, s.Result.History) })
	}
	if a.PossibleNextMoves != nil {
		t.Run("PossibleNextMoves", func(t *testing.T) { a.PossibleNextMoves.AssertPossibleNextMoves(t, s.Result.PossibleNextMoves) })
	}
	if a.Variant != nil {
		t.Run("Variant", func(t *testing.T) { a.Variant.AssertVariant(t, s.Input) })
	}
}

// Array performs array assertions.
type Array struct {
	Length []IntAssertion
}

// Length checks the length of slices for conditions
func Length(testCases ...IntAssertion) Array {
	return Array{
		Length: testCases,
	}
}

// AssertBytes checks properties of byte arrays.
func (a Array) AssertBytes(t *testing.T, val []byte) {
	t.Run("length", func(t *testing.T) {
		for _, testCase := range a.Length {
			testCase.AssertInt(t, len(val))
		}
	})
}

// AssertPossibleNextMoves checks properties of possible next moves.
func (a Array) AssertPossibleNextMoves(t *testing.T, val [][]poc.Move) {
	t.Run("length", func(t *testing.T) {
		for _, testCase := range a.Length {
			testCase.AssertInt(t, len(val))
		}
	})
}

// GTE checks that values are greater than or equal to the value.
func GTE(low int64) Between {
	return Between{
		Low:  low,
		High: math.MaxInt64,
	}
}

// GT checks that values are strictly greater than the value.
func GT(low int64) Between {
	return GTE(low + 1)
}

// LTE checks that values are less than or equal to the value.
func LTE(high int64) Between {
	return Between{
		Low:  math.MinInt64,
		High: high,
	}
}

// LT checks that values are strictly less than the value.
func LT(high int64) Between {
	return LTE(high - 1)
}

// Eq checks that values are equal to the value
func Eq(expected int64) Between {
	return Between{
		Low:  expected,
		High: expected,
	}
}

// Between checks that values fall between low and high inclusive.
type Between struct {
	Low  int64
	High int64
}

// AssertInt checks int types.
func (b Between) AssertInt(t *testing.T, val int) {
	b.AssertInt64(t, int64(val))
}

// AssertInt32 checks int32 types.
func (b Between) AssertInt32(t *testing.T, val int32) {
	b.AssertInt64(t, int64(val))
}

// AssertInt64 checks int64 types.
func (b Between) AssertInt64(t *testing.T, val int64) {
	if val > b.High {
		t.Errorf("actual (%d) should not have violated upper bound %d", val, b.High)
	}
	if val < b.Low {
		t.Errorf("actual (%d) should not have violated lower bound %d", val, b.Low)
	}
}

// IntAssertion asserts properties of byte arrays.
type IntAssertion interface {
	AssertInt(*testing.T, int)
}

// BytesAssertion asserts properties of byte arrays.
type BytesAssertion interface {
	AssertBytes(*testing.T, []byte)
}

// GameIDAssertion asserts properties of game ids.
type GameIDAssertion interface {
	AssertGameID(*testing.T, int64)
}

// BoardAssertion asserts properties of boards.
type BoardAssertion interface {
	AssertBoard(*testing.T, poc.Board)
}

// HistoryAssertion asserts properties of histories.
type HistoryAssertion interface {
	AssertHistory(*testing.T, poc.History)
}

// PossibleNextMovesAssertion asserts properties of possible next movess.
type PossibleNextMovesAssertion interface {
	AssertPossibleNextMoves(*testing.T, [][]poc.Move)
}

// VariantAssertion asserts properties of variants.
type VariantAssertion interface {
	AssertVariant(*testing.T, poc.Variant)
}

// ScoreAssertion asserts properties of scores.
type ScoreAssertion interface {
	AssertScore(*testing.T, int32)
}
