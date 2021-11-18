package main

import (
	"net/http"

	"github.com/slcjordan/poc/boot"
	"github.com/slcjordan/poc/config"
	_ "github.com/slcjordan/poc/config/env"
	_ "github.com/slcjordan/poc/logger/stdlib"
)

func main() {
	config.DB.ShouldParse = true
	config.HTTP.ShouldParse = true
	config.MustParse()
	pool := boot.PGXConnect(config.DB.ConnString)

	boot.MustServe(&http.Server{
		Addr:    config.HTTP.ListenAddress,
		Handler: boot.APIServer(pool),
	})
}
