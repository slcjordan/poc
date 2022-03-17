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

type Validator interface {
	Validate() error
}

var parsers [3][]Parser

var validators []Validator

// Order determines parsing order
//go:generate stringer -type=Order
type Order int

// First will be the first to run, meaning the value can be overwritten by later parsers.
const (
	First Order = iota
	Second
	Third
)

// Register should be called by a parse package's init function.
func Register(p Parser, o Order) {
	parsers[o] = append(parsers[o], p)
}

// Validate should be called by a package's init function.
func Validate(v Validator) {
	validators = append(validators, v)
}

// Validate should be called by a package's init function.
func ValidateFunc(f func() error) {
	Validate(ValidatorFunc(f))
}

type ValidatorFunc func() error

func (v ValidatorFunc) Validate() error {
	return v()
}

func Parse() error {
	for _, o := range []Order{First, Second, Third} {
		for _, p := range parsers[o] {
			err := p.ParseConfig()
			if err != nil {
				logger.Errorf(context.Background(), "whlie parsing config: %s", err)
				return err
			}
		}
	}
	for _, v := range validators {
		err := v.Validate()
		if err != nil {
			logger.Errorf(context.Background(), "whlie validating config: %s", err)
			return err
		}
	}
	return nil
}

// MustParse calls ParseConfig for all registered parsers and panics if there
// is an error.
func MustParse() {
	err := Parse()
	if err != nil {
		panic(err)
	}
}
