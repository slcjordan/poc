package config

import "log"

type Parser interface {
	ParseConfig() error
}

var parsers []Parser

func Register(p Parser) {
	parsers = append(parsers, p)
}

var DB struct {
	ShouldParse bool
	ConnString  string
}

var HTTP struct {
	ShouldParse   bool
	ListenAddress string
}

func MustParse() {
	for _, p := range parsers {
		err := p.ParseConfig()
		if err != nil {
			log.Fatal(err)
		}
	}
}
