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

// DB should only be parsed when DB.ShoulParse is true.
var DB struct {
	ShouldParse bool
	ConnString  string
}

// HTTP should only be parsed when HTTP.ShoulParse is true.
var HTTP struct {
	ShouldParse   bool
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
