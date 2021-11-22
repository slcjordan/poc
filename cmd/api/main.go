package main

import (
	"github.com/slcjordan/poc/boot"
	_ "github.com/slcjordan/poc/config/env"
	"github.com/slcjordan/poc/logger/stdlib"
)

func main() {
	stdlib.JSON = true
	boot.MustServeFromConfig()
}
