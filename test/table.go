package test

import (
	"context"
	"errors"
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

// IsCategory checks the category of an error.
type Category struct {
	Expected poc.ErrorCategory
}

// AssertError will only pass if not nil and of correct category.
func (c Category) AssertError(t *testing.T, err error) {
	var categorized poc.Error
	if !errors.As(err, &categorized) {
		t.Fatalf("Error of type %T does not wrap nor is of type poc.Error: %s", err, err)
	}
	actual := categorized.Category
	if actual != c.Expected {
		t.Fatalf("Expected error of category %s but got %s", c.Expected, actual)
	}
}

// IsNil fails if error is not nil.
type IsNil struct{}

// AssertError will only pass if not nil and of correct category.
func (i IsNil) AssertError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Error of type %T is not nil: %s", err, err)
	}
}
