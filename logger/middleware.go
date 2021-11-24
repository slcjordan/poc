package logger

import (
	"context"
	"encoding/json"

	"github.com/slcjordan/poc"
)

// A ByteCaller processes bytes.
type ByteCaller interface {
	CallBytes(context.Context, []byte) ([]byte, error)
}

type middlewareKey string

// MiddlewareKey is used for logger context.
const MiddlewareKey = middlewareKey("key")

// Middleware injects context values from commands.
type Middleware struct {
	NextStartGame poc.StartGameCaller
	NextBytes     ByteCaller
}

// CallBytes sets up logging context based on start game values.
func (m Middleware) CallBytes(ctx context.Context, input []byte) ([]byte, error) {
	var val Values
	prev, ok := ctx.Value(MiddlewareKey).(Values)
	if ok {
		val = prev
	}
	val.Bytes = input
	return m.NextBytes.CallBytes(context.WithValue(ctx, MiddlewareKey, val), input)
}

// CallStartGame sets up logging context based on start game values.
func (m Middleware) CallStartGame(ctx context.Context, input poc.StartGame) (poc.StartGame, error) {
	var val Values
	prev, ok := ctx.Value(MiddlewareKey).(Values)
	if ok {
		val = prev
	}
	val.StartGame = input
	return m.NextStartGame.CallStartGame(context.WithValue(ctx, MiddlewareKey, val), input)
}

// BytesMiddleware sets up sMiddlewareKeymiddleware.
func BytesMiddleware(next ByteCaller) Middleware {
	return Middleware{
		NextBytes: next,
	}
}

// StartGameMiddleware sets up start-game middleware.
func StartGameMiddleware(next poc.StartGameCaller) Middleware {
	return Middleware{
		NextStartGame: next,
	}
}

// Values holds important logging context information.
type Values struct {
	Bytes       []byte
	StartGame   poc.StartGame
	PerformMove poc.PerformMove
	ListGames   poc.ListGames
	Message     string
}

// MarshalJSON formats context for single-level logging values and masks PII.
func (v Values) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message                             string `json:"message"`
		NumBytes                            int    `json:"num_input_bytes,omitempty"`
		StartGameVariantMaxTimesThroughDeck int32  `json:"start_game_variant_max_times_through_deck,omitempty"`
		PerformMoveGameID                   int64  `json:"perform_move_game_id,omitempty"`
		PerformMoveNumCardsToMove           int    `json:"perform_move_num_cards_to_move,omitempty"`
		ListGamesOffset                     int32  `json:"list_games_offset,omitempty"`
		ListGamesLimit                      int32  `json:"list_games_limit,omitempty"`
	}{
		Message:                             v.Message,
		NumBytes:                            len(v.Bytes),
		StartGameVariantMaxTimesThroughDeck: v.StartGame.Input.MaxTimesThroughDeck,
		PerformMoveGameID:                   v.PerformMove.Input.GameID,
		PerformMoveNumCardsToMove:           len(v.PerformMove.Input.Move),
		ListGamesOffset:                     v.ListGames.Input.Offset,
		ListGamesLimit:                      v.ListGames.Input.Limit,
	})
}
