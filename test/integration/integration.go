package integration

import (
	"testing"

	"github.com/slcjordan/poc/config"
	_ "github.com/slcjordan/poc/config/env" // fetch config from environment variables.
)

// TestAPI checks that the API is ready for running.
func TestAPI(t *testing.T) {
	config.MustParse()
}
