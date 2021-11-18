package main

import (
	"github.com/slcjordan/poc/boot"
	_ "github.com/slcjordan/poc/config/env"
	_ "github.com/slcjordan/poc/logger/stdlib"
)

func main() {
	boot.MustServeFromConfig()
}
