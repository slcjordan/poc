package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/slcjordan/poc"
	"github.com/slcjordan/poc/db/sqlc"
	"github.com/slcjordan/poc/logger"
)

type Lookup struct {
	Pool *pgxpool.Pool
}

func (l *Lookup) CallPerformMove(ctx context.Context, move poc.PerformMove) (poc.PerformMove, error) {
	conn, err := l.Pool.Acquire(ctx)
	if err != nil {
		logger.Infof(ctx, "could not acquire connection: %s", err)
		return move, poc.Error{Actual: errors.New("db unavailable"), Category: poc.UnavailableError}
	}
	_, err = sqlc.New(conn).LookupGameDetail(ctx, move.Input.GameID)
	if err != nil {
		logger.Errorf(ctx, "could lookup game %d: %s", move.Input.GameID, err)
		return move, poc.Error{Actual: errors.New("could not find game"), Category: poc.UnknownError}
	}
	return move, nil
}
