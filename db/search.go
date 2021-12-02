package db

import (
	"context"
	"errors"

	"github.com/slcjordan/poc"
	"github.com/slcjordan/poc/db/sqlc"
	"github.com/slcjordan/poc/logger"
)

// Search uses search parameters to find multiple results.
type Search struct {
	Pool Pool
}

// CallListGames has a default ordering by game id.
func (s *Search) CallListGames(ctx context.Context, list poc.ListGames) (poc.ListGames, error) {
	conn, err := s.Pool.Acquire(ctx)
	if err != nil {
		logger.Infof(ctx, "could not acquire connection: %s", err)
		return list, poc.Error{Actual: errors.New("db unavailable"), Category: poc.UnavailableError}
	}
	defer conn.Release()
	rows, err := sqlc.New(conn).SearchGame(ctx, sqlc.SearchGameParams{
		Limit:  list.Cursor.Limit,
		Offset: list.Cursor.Offset,
	})
	if err != nil {
		logger.Errorf(ctx, "could not list games: %s", err)
		return list, poc.Error{Actual: errors.New("could not list games"), Category: poc.UnknownError}
	}
	list.Games = make([]poc.SavedGameSummary, len(rows))
	for i, row := range rows {
		list.Games[i].Score = row.Score
		list.Games[i].GameID = row.ID
	}
	return list, nil
}
