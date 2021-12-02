package rules_test

import (
	"math/rand"
	"testing"

	"github.com/slcjordan/poc/pipeline"
	"github.com/slcjordan/poc/rules"
	"github.com/slcjordan/poc/test"
	"github.com/slcjordan/poc/test/logger"
)

func TestNextMoves(t *testing.T) {
	logger.RegisterVerbose(t)
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
	}.Run(t)
}
