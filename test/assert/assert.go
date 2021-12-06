package assert

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/slcjordan/poc"
	"github.com/slcjordan/poc/logger"
)

type Assertion struct {
	PerformMove PerformMove
	StartGame   StartGame
	ListGames   ListGames
	Error       Error

	noError bool
}

func New() *Assertion {
	var assertion Assertion
	assertion.ListGames = newListGames(&assertion)
	assertion.PerformMove = newPerformMove(&assertion)
	assertion.StartGame = newStartGame(&assertion)
	assertion.Error = newError(&assertion)
	return &assertion
}

func (a *Assertion) CheckListGames(t *testing.T, desc string, val poc.ListGames) {
	a.ListGames.CheckListGames(t, desc+"ListGames", val)
}

func (a *Assertion) CheckPerformMove(t *testing.T, desc string, val poc.PerformMove) {
	a.PerformMove.CheckPerformMove(t, desc+"PerformMove", val)
}

func (a *Assertion) CheckStartGame(t *testing.T, desc string, val poc.StartGame) {
	a.StartGame.CheckStartGame(t, desc+"StartGame", val)
}

func (a *Assertion) CheckError(t *testing.T, desc string, val error) {
	if a.noError {
		t.Run(desc+"no error", func(t *testing.T) {
			if val != nil {
				t.Errorf("expected nil error, but got: %s", val)
			}
		})
		return
	}
	var perr poc.Error
	if !errors.As(val, &perr) {
		logger.Infof(context.Background(), "assert package has wrapped %T error as poc.Error", val)
		perr = poc.Error{Actual: val}
	}
	a.Error.CheckError(t, desc+"error", perr)
}

func (a *Assertion) NoError() *Assertion {
	a.noError = true
	return a
}

type Equals int64

func (e Equals) CheckUint8(t *testing.T, desc string, val uint8) {
	e.CheckInt64(t, desc, int64(val))
}

func (e Equals) CheckUint16(t *testing.T, desc string, val uint16) {
	e.CheckInt64(t, desc, int64(val))
}

func (e Equals) CheckUint32(t *testing.T, desc string, val uint32) {
	e.CheckInt64(t, desc, int64(val))
}

func (e Equals) CheckUint64(t *testing.T, desc string, val uint64) {
	e.CheckInt64(t, desc, int64(val))
}

func (e Equals) CheckInt(t *testing.T, desc string, val int) {
	e.CheckInt64(t, desc, int64(val))
}

func (e Equals) CheckInt8(t *testing.T, desc string, val int8) {
	e.CheckInt64(t, desc, int64(val))
}

func (e Equals) CheckInt16(t *testing.T, desc string, val int16) {
	e.CheckInt64(t, desc, int64(val))
}

func (e Equals) CheckInt32(t *testing.T, desc string, val int32) {
	e.CheckInt64(t, desc, int64(val))
}

func (e Equals) CheckInt64(t *testing.T, desc string, val int64) {
	t.Run(fmt.Sprintf("%s equals %d", desc, e), func(t *testing.T) {
		expected := int64(e)
		if val != expected {
			t.Errorf("expected %d but got %d", expected, val)
		}
	})
}
