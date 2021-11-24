package rules_test

import (
	"math/rand"
	"testing"

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
	}.Run(t)
}
