package boot

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/slcjordan/poc/config"
	"github.com/slcjordan/poc/db"
	"github.com/slcjordan/poc/encoding/json"
	"github.com/slcjordan/poc/handler"
	"github.com/slcjordan/poc/logger"
	"github.com/slcjordan/poc/pipeline"
	"github.com/slcjordan/poc/router"
	"github.com/slcjordan/poc/rules"
)

func mustConnect(connString string) *pgxpool.Pool {
	conn, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		logger.Fatalf(context.Background(), "could not connect: %s", err)
	}
	logger.Infof(
		context.Background(),
		"connected to %#v database on host %#v and port %d",
		conn.Config().ConnConfig.Config.Database,
		conn.Config().ConnConfig.Config.Host,
		conn.Config().ConnConfig.Config.Port,
	)
	return conn
}

func APIServer() http.Handler {
	mux := router.NewMux()
	v1HydrateParams := router.V1HydrateURLAndQueryParams{OffsetKey: "offset", LimitKey: "limit"}
	pool := mustConnect(config.DB.ConnString)
	save := &db.Save{Pool: pool}
	search := &db.Search{Pool: pool}
	lookup := &db.Lookup{Pool: pool}

	mux.RegisterRoutesV1(router.V1Handlers{
		StartGame: handler.StartGame{
			Encoding: json.V1{},
			Command: pipeline.StartGame{
				rules.Shuffle{Source: rand.NewSource(time.Now().UnixNano())},
				save,
				rules.NextMove{},
			},
		},
		PerformMove: handler.PerformMove{
			Encoding: json.V1{},
			Command: pipeline.PerformMove{
				v1HydrateParams,
				lookup,
				rules.Validate{},
				save,
				rules.NextMove{},
			},
		},
		ListGames: handler.ListGames{
			Encoding: json.V1{},
			Command: pipeline.ListGames{
				v1HydrateParams,
				search,
			},
		},
	},
		middleware.Logger,
	)
	return mux
}

func MustServe() {
	config.DB.ShouldParse = true
	config.HTTP.ShouldParse = true
	config.MustParse()

	logger.Infof(context.Background(), "Listening at %#v", config.HTTP.ListenAddress)
	err := http.ListenAndServe(config.HTTP.ListenAddress, APIServer())
	if err != nil {
		logger.Fatalf(context.Background(), "could not server: %s", err)
	}
}
