// Code generated by cmd/assert; DO NOT EDIT.

package assert

import (
	"fmt"
	"testing"

	"github.com/slcjordan/poc"
)

type ListGames struct {
	assertion            *Assertion
	cursorLimitCheckers  []Int32Checker
	cursorOffsetCheckers []Int32Checker

	Games SavedGameSummaryArray1D
}

func newListGames(assertion *Assertion) ListGames {
	return ListGames{
		assertion: assertion,
		Games:     newSavedGameSummaryArray1D(assertion),
	}
}

func (parent *ListGames) CursorLimit(checkers ...Int32Checker) *Assertion {
	parent.cursorLimitCheckers = checkers
	return parent.assertion
}

func (parent *ListGames) CursorOffset(checkers ...Int32Checker) *Assertion {
	parent.cursorOffsetCheckers = checkers
	return parent.assertion
}

func (parent *ListGames) CheckListGames(t *testing.T, desc string, val poc.ListGames) {
	for _, checker := range parent.cursorLimitCheckers {
		checker.CheckInt32(t, desc+".Cursor.Limit", val.Cursor.Limit)
	}
	for _, checker := range parent.cursorOffsetCheckers {
		checker.CheckInt32(t, desc+".Cursor.Offset", val.Cursor.Offset)
	}
	parent.Games.CheckSavedGameSummaryArray1D(t, desc+".Games", val.Games)
}

type PerformMove struct {
	assertion *Assertion

	Next            MoveArray1D
	SavedGameDetail SavedGameDetail
}

func newPerformMove(assertion *Assertion) PerformMove {
	return PerformMove{
		assertion:       assertion,
		Next:            newMoveArray1D(assertion),
		SavedGameDetail: newSavedGameDetail(assertion),
	}
}

func (parent *PerformMove) CheckPerformMove(t *testing.T, desc string, val poc.PerformMove) {
	parent.Next.CheckMoveArray1D(t, desc+".Next", val.Next)
	parent.SavedGameDetail.CheckSavedGameDetail(t, desc+".SavedGameDetail", val.SavedGameDetail)
}

type SavedGameDetail struct {
	assertion      *Assertion
	gameIDCheckers []Int64Checker

	Board             Board
	History           MoveArray2D
	PossibleNextMoves MoveArray2D
	Variant           Variant
}

func newSavedGameDetail(assertion *Assertion) SavedGameDetail {
	return SavedGameDetail{
		assertion:         assertion,
		Board:             newBoard(assertion),
		History:           newMoveArray2D(assertion),
		PossibleNextMoves: newMoveArray2D(assertion),
		Variant:           newVariant(assertion),
	}
}

func (parent *SavedGameDetail) GameID(checkers ...Int64Checker) *Assertion {
	parent.gameIDCheckers = checkers
	return parent.assertion
}

func (parent *SavedGameDetail) CheckSavedGameDetail(t *testing.T, desc string, val poc.SavedGameDetail) {
	for _, checker := range parent.gameIDCheckers {
		checker.CheckInt64(t, desc+".GameID", val.GameID)
	}
	parent.Board.CheckBoard(t, desc+".Board", val.Board)
	parent.History.CheckMoveArray2D(t, desc+".History", val.History)
	parent.PossibleNextMoves.CheckMoveArray2D(t, desc+".PossibleNextMoves", val.PossibleNextMoves)
	parent.Variant.CheckVariant(t, desc+".Variant", val.Variant)
}

type SavedGameSummary struct {
	assertion      *Assertion
	gameIDCheckers []Int64Checker
	scoreCheckers  []Int32Checker
}

func newSavedGameSummary(assertion *Assertion) SavedGameSummary {
	return SavedGameSummary{
		assertion: assertion,
	}
}

func (parent *SavedGameSummary) GameID(checkers ...Int64Checker) *Assertion {
	parent.gameIDCheckers = checkers
	return parent.assertion
}

func (parent *SavedGameSummary) Score(checkers ...Int32Checker) *Assertion {
	parent.scoreCheckers = checkers
	return parent.assertion
}

func (parent *SavedGameSummary) CheckSavedGameSummary(t *testing.T, desc string, val poc.SavedGameSummary) {
	for _, checker := range parent.gameIDCheckers {
		checker.CheckInt64(t, desc+".GameID", val.GameID)
	}
	for _, checker := range parent.scoreCheckers {
		checker.CheckInt32(t, desc+".Score", val.Score)
	}
}

type StartGame struct {
	assertion *Assertion

	SavedGameDetail SavedGameDetail
	Variant         Variant
}

func newStartGame(assertion *Assertion) StartGame {
	return StartGame{
		assertion:       assertion,
		SavedGameDetail: newSavedGameDetail(assertion),
		Variant:         newVariant(assertion),
	}
}

func (parent *StartGame) CheckStartGame(t *testing.T, desc string, val poc.StartGame) {
	parent.SavedGameDetail.CheckSavedGameDetail(t, desc+".SavedGameDetail", val.SavedGameDetail)
	parent.Variant.CheckVariant(t, desc+".Variant", val.Variant)
}

type SavedGameSummaryArray1D struct {
	assertion      *Assertion
	lengthCheckers []IntChecker
	nth            map[int]SavedGameSummary

	ForEach SavedGameSummary
}

func newSavedGameSummaryArray1D(assertion *Assertion) SavedGameSummaryArray1D {
	return SavedGameSummaryArray1D{
		assertion: assertion,
		nth:       make(map[int]SavedGameSummary),
		ForEach:   newSavedGameSummary(assertion),
	}
}

func (a *SavedGameSummaryArray1D) Nth(i int) SavedGameSummary {
	prev, ok := a.nth[i]
	if ok {
		return prev
	}
	result := newSavedGameSummary(a.assertion)
	a.nth[i] = result
	return result
}

func (a *SavedGameSummaryArray1D) Length(checkers ...IntChecker) *Assertion {
	a.lengthCheckers = checkers
	return a.assertion
}

func (a *SavedGameSummaryArray1D) CheckSavedGameSummaryArray1D(t *testing.T, desc string, val []poc.SavedGameSummary) {
	for _, checker := range a.lengthCheckers {
		checker.CheckInt(t, desc+".length", len(val))
	}
	for i, checker := range a.nth {
		checker.CheckSavedGameSummary(t, desc+fmt.Sprintf("[%d]", i), val[i])
	}
	for _, curr := range val {
		a.ForEach.CheckSavedGameSummary(t, desc+".ForEach", curr)
	}
}

type MoveArray1D struct {
	assertion      *Assertion
	lengthCheckers []IntChecker
	nth            map[int]Move

	ForEach Move
}

func newMoveArray1D(assertion *Assertion) MoveArray1D {
	return MoveArray1D{
		assertion: assertion,
		nth:       make(map[int]Move),
		ForEach:   newMove(assertion),
	}
}

func (a *MoveArray1D) Nth(i int) Move {
	prev, ok := a.nth[i]
	if ok {
		return prev
	}
	result := newMove(a.assertion)
	a.nth[i] = result
	return result
}

func (a *MoveArray1D) Length(checkers ...IntChecker) *Assertion {
	a.lengthCheckers = checkers
	return a.assertion
}

func (a *MoveArray1D) CheckMoveArray1D(t *testing.T, desc string, val []poc.Move) {
	for _, checker := range a.lengthCheckers {
		checker.CheckInt(t, desc+".length", len(val))
	}
	for i, checker := range a.nth {
		checker.CheckMove(t, desc+fmt.Sprintf("[%d]", i), val[i])
	}
	for _, curr := range val {
		a.ForEach.CheckMove(t, desc+".ForEach", curr)
	}
}

type MoveArray2D struct {
	assertion      *Assertion
	lengthCheckers []IntChecker
	nth            map[int]MoveArray1D

	ForEach MoveArray1D
}

func newMoveArray2D(assertion *Assertion) MoveArray2D {
	return MoveArray2D{
		assertion: assertion,
		nth:       make(map[int]MoveArray1D),
		ForEach:   newMoveArray1D(assertion),
	}
}

func (a *MoveArray2D) Nth(i int) MoveArray1D {
	prev, ok := a.nth[i]
	if ok {
		return prev
	}
	result := newMoveArray1D(a.assertion)
	a.nth[i] = result
	return result
}

func (a *MoveArray2D) Length(checkers ...IntChecker) *Assertion {
	a.lengthCheckers = checkers
	return a.assertion
}

func (a *MoveArray2D) CheckMoveArray2D(t *testing.T, desc string, val [][]poc.Move) {
	for _, checker := range a.lengthCheckers {
		checker.CheckInt(t, desc+".length", len(val))
	}
	for i, checker := range a.nth {
		checker.CheckMoveArray1D(t, desc+fmt.Sprintf("[%d]", i), val[i])
	}
	for _, curr := range val {
		a.ForEach.CheckMoveArray1D(t, desc+".ForEach", curr)
	}
}
