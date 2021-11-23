package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/slcjordan/poc/boot"
	"github.com/slcjordan/poc/config"
	_ "github.com/slcjordan/poc/config/env" // fetch config from environment variables.
	"github.com/slcjordan/poc/test/logger"
)

// TestAPI checks that the API is ready for running.
func TestAPI(t *testing.T) {
	config.MustParse()
	logger.RegisterVerbose(t)

	pool := boot.PGXConnect(config.DB.ConnString)
	server := httptest.NewServer(boot.APIServer(pool))
	defer server.Close()
	client := server.Client()

	t.Run("check start game endpoints", startGame(client, server.URL))
}

func startGame(client *http.Client, url string) func(*testing.T) {
	return func(t *testing.T) {
		response, err := client.Post(url+"/v1/game/start", "application/json", strings.NewReader("{}"))
		if err != nil {
			t.Fatalf("expected no error but got: %s", err)
		}
		if response.StatusCode != 200 {
			t.Fatalf("expected no error code but got: %d", response.StatusCode)
		}
	}
}
