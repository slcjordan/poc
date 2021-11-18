package boot

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	chi "github.com/go-chi/chi/v5"
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

// PGXConnect connects to db using pgx library
func PGXConnect(connString string) *pgxpool.Pool {
	conn, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		logger.Errorf(context.Background(), "could not connect: %s", err)
		panic(err)
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

// MustServe serves the api server and fatally exits on error.
func MustServe(s *http.Server) {
	logger.Infof(context.Background(), "Listening at %#v", s.Addr)
	err := s.ListenAndServe()
	if err != nil {
		logger.Errorf(context.Background(), "while serving the api server: %s", err)
		panic(err)
	}
}

// MustServeFromConfig parses config and serves
func MustServeFromConfig() {
	config.DB.ShouldParse = true
	config.HTTP.ShouldParse = true
	config.MustParse()
	pool := PGXConnect(config.DB.ConnString)

	MustServe(&http.Server{
		Addr:    config.HTTP.ListenAddress,
		Handler: APIServer(pool),
	})
}

// APIServer connects to database; sets up routes and handlers.
func APIServer(pool db.PgxPoolIface) chi.Router {
	v1HydrateParams := router.V1HydrateURLAndQueryParams{OffsetKey: "offset", LimitKey: "limit"}
	save := &db.Save{Pool: pool}
	search := &db.Search{Pool: pool}
	lookup := &db.Lookup{Pool: pool}

	return router.New(router.V1Handlers{
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
}
