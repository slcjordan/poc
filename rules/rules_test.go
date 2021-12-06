package rules_test

import (
	"math/rand"
	"testing"

	"github.com/slcjordan/poc"
	"github.com/slcjordan/poc/pipeline"
	"github.com/slcjordan/poc/rules"
	"github.com/slcjordan/poc/test/assert"
	"github.com/slcjordan/poc/test/harness"
	"github.com/slcjordan/poc/test/logger"
)

func TestNextMoves(t *testing.T) {
	logger.RegisterVerbose(t)
	harness.StartGame{
		{
			Desc:    "Sanity check",
			Command: rules.NextMove{},
			Result:  assert.New().NoError(),
		},
		{
			Desc: "Given 0 random seed",
			Command: pipeline.StartGame{
				rules.Shuffle{rand.NewSource(0)},
				rules.NextMove{},
			},
			Result: assert.New().StartGame.SavedGameDetail.PossibleNextMoves.Length(assert.Equals(8)),
		},
		{
			Desc: "Max times variant is respected",
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
			Result: assert.New().NoError().
				StartGame.SavedGameDetail.PossibleNextMoves.Length(assert.Equals(0)),
		},
	}.Run(t)
}
