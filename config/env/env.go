package env

import (
	"os"

	"github.com/slcjordan/poc/config"
)

func init() {
	config.Register(parser{})
}

type parser struct{}

// ParseConfig uses non-empty environment variables.
func (p parser) ParseConfig() error {
	connString := os.Getenv("DB_CONN_STRING")
	if connString != "" {
		config.DB.ConnString = connString
	}
	listenAddress := os.Getenv("HTTP_LISTEN_ADDRESS")
	if listenAddress != "" {
		config.HTTP.ListenAddress = listenAddress
	}
	return nil
}
