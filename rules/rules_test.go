package rules_test

import (
	"math/rand"
	"testing"

	"github.com/slcjordan/poc/pipeline"
	"github.com/slcjordan/poc/rules"
	"github.com/slcjordan/poc/test/assert"
	"github.com/slcjordan/poc/test/harness"
	"github.com/slcjordan/poc/test/logger"
)

func TestNextMoves(t *testing.T) {
	logger.RegisterVerbose(t)
	/*
		test.StartGame{
			{
				Desc:    "next move sanity check",
				Command: rules.NextMove{},
				Error:   test.IsNil{},
			},
			{
				Desc:    "shuffle sanity check",
				Command: rules.Shuffle{rand.NewSource(0)},
				Error:   test.IsNil{},
			},
			{
				Desc: "next move can be found from shuffled game without error",
				Command: pipeline.StartGame{
					rules.Shuffle{rand.NewSource(0)},
					rules.NextMove{},
				},
				Result: test.Assert{
					PossibleNextMoves: test.Length(test.Eq(8)),
				},
			},
			{
				Desc: "cannot move more than max times through deck",
				Command: pipeline.StartGame{
					rules.Shuffle{rand.NewSource(0)},
					rules.NextMove{},
				},
				Input: poc.StartGame{
					Variant: poc.Variant{
						MaxTimesThroughDeck: 1,
					},
					SavedGameDetail: poc.SavedGameDetail{
						History: [][]poc.Move{
							{
								{
									NewPileNum: 0, // talon has already been returned to the stock once
								},
							},
						},
					},
				},
				Result: test.Assert{
					PossibleNextMoves: test.Length(test.Eq(0)),
				},
			},
		}.Run(t)
	*/
	harness.StartGame{
		{
			Desc: "Given 0 random seed",
			Command: pipeline.StartGame{
				rules.Shuffle{rand.NewSource(0)},
				rules.NextMove{},
			},
			Result: assert.New().StartGame.SavedGameDetail.PossibleNextMoves.Length(assert.Equals(8)),
		},
	}.Run(t)
}
