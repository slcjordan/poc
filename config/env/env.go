package env

import (
	"os"

	"github.com/slcjordan/poc/config"
)

func init() {
	config.Register(parser{})
}

type parser struct{}

func (p parser) ParseConfig() error {
	if config.DB.ShouldParse {
		config.DB.ConnString = os.Getenv("DB_CONN_STRING")
	}
	if config.HTTP.ShouldParse {
		config.HTTP.ListenAddress = os.Getenv("HTTP_LISTEN_ADDRESS")
	}
	return nil
}
