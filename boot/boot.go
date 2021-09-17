package boot

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/slcjordan/poc"
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

// APIServer connects to database; sets up routes and handlers.
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
				Debug{Name: "decoding"},
				rules.Shuffle{Source: rand.NewSource(time.Now().UnixNano())},
				Debug{Name: "shuffle"},
				save,
				Debug{Name: "save"},
				rules.NextMove{},
				Debug{Name: "next move"},
			},
		},
		PerformMove: handler.PerformMove{
			Encoding: json.V1{},
			Command: pipeline.PerformMove{
				Debug{Name: "decoding"},
				v1HydrateParams,
				Debug{Name: "hydrate"},
				lookup,
				Debug{Name: "lookup"},
				rules.Validate{},
				Debug{Name: "validate"},
				save,
				Debug{Name: "save"},
				rules.NextMove{},
				Debug{Name: "next move"},
			},
		},
		ListGames: handler.ListGames{
			Encoding: json.V1{},
			Command: pipeline.ListGames{
				Debug{Name: "decoding"},
				v1HydrateParams,
				Debug{Name: "hydrate"},
				search,
				Debug{Name: "search"},
			},
		},
	},
		middleware.Logger,
	)
	return mux
}

// MustServe serves the api server and fatally exits on error.
func MustServe() {
	config.DB.ShouldParse = true
	config.HTTP.ShouldParse = true
	config.MustParse()

	logger.Infof(context.Background(), "Listening at %#v", config.HTTP.ListenAddress)
	err := http.ListenAndServe(config.HTTP.ListenAddress, APIServer())
	if err != nil {
		logger.Errorf(context.Background(), "while serving the api server: %s", err)
		panic(err)
	}
}

// Debug debugs the app.
type Debug struct {
	Name    string
	Enabled bool
}

// CallStartGame debugs a start game command.
func (d Debug) CallStartGame(ctx context.Context, game poc.StartGame) (poc.StartGame, error) {
	if d.Enabled {
		logger.Infof(ctx, "after %#v the value is %v", d.Name, game)
	}
	return game, nil
}

// CallPerformMove debugs a perform move command.
func (d Debug) CallPerformMove(ctx context.Context, move poc.PerformMove) (poc.PerformMove, error) {
	if d.Enabled {
		logger.Infof(ctx, "after %#v the value is %v", d.Name, move)
	}
	return move, nil
}

// CallListGames debugs a list game command.
func (d Debug) CallListGames(ctx context.Context, list poc.ListGames) (poc.ListGames, error) {
	if d.Enabled {
		logger.Infof(ctx, "after %#v the value is %v", d.Name, list)
	}
	return list, nil
}
