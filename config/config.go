package config

import (
	"context"

	"github.com/slcjordan/poc/logger"
)

// A Parser should only write to config state when ParseConfig is called.
// /boot package may set ShouldParse values and call a Parse function
// Otherwise, all access to config state should be read-only.
type Parser interface {
	ParseConfig() error
}

var parsers []Parser

// Register should be called by a parse package's init function.
func Register(p Parser) {
	parsers = append(parsers, p)
}

// DB options.
var DB struct {
	ConnString string
}

// HTTP options.
var HTTP struct {
	ListenAddress string
}

// MustParse calls ParseConfig for all registered parsers and panics if there
// is an error.
func MustParse() {
	for _, p := range parsers {
		err := p.ParseConfig()
		if err != nil {
			logger.Errorf(context.Background(), "whlie parsing config: %s", err)
			panic(err)
		}
	}
}
