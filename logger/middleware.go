package logger

import (
	"context"
	"encoding/json"

	"github.com/slcjordan/poc"
)

// KeyValue holds custom values.
type KeyValue struct {
	Key   string
	Value interface{}
}

// WithValues adds extra info to a logger call.
type WithValues map[string]interface{}

// Infof may be used to log error conditions that are the fault of external applications.
func (w WithValues) Infof(ctx context.Context, format string, a ...interface{}) {
	for _, l := range infos {
		l.Printf(w.Context(ctx), format+"\n", a...)
	}
}

// Errorf must be used to log error conditions that may be caused by errors in this application.
func (w WithValues) Errorf(ctx context.Context, format string, a ...interface{}) {
	for _, l := range errors {
		l.Printf(w.Context(ctx), format+"\n", a...)
	}
}

// Create logger context with the given values.
func (w WithValues) Context(ctx context.Context) context.Context {
	var curr Values
	prev, ok := ctx.Value(MiddlewareKey).(Values)
	if ok {
		curr = prev
	}
	for key, val := range w {
		curr.Extra = append(curr.Extra, KeyValue{Key: key, Value: val})
	}
	return context.WithValue(ctx, MiddlewareKey, curr)
}

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

func (m Middleware) StartGameUse(next poc.StartGameCaller) poc.StartGameCaller {
	m.NextStartGame = next
	return m
}

// BytesMiddleware sets up sMiddlewareKeymiddleware.
func BytesMiddleware(next ByteCaller) Middleware {
	return Middleware{
		NextBytes: next,
	}
}

// Values holds important logging context information.
type Values struct {
	Bytes       []byte
	StartGame   poc.StartGame
	PerformMove poc.PerformMove
	ListGames   poc.ListGames
	Message     string
	Extra       []KeyValue
}

// MarshalJSON formats context for single-level logging values and masks PII.
func (v Values) MarshalJSON() ([]byte, error) {
	var extra map[string]interface{}
	if v.Extra != nil {
		extra = make(map[string]interface{})
		for _, item := range v.Extra {
			extra[item.Key] = item.Value
		}
	}
	return json.Marshal(struct {
		Message                             string                 `json:"message"`
		NumBytes                            int                    `json:"num_input_bytes,omitempty"`
		StartGameVariantMaxTimesThroughDeck int32                  `json:"start_game_variant_max_times_through_deck,omitempty"`
		PerformMoveGameID                   int64                  `json:"perform_move_game_id,omitempty"`
		PerformMoveNumCardsToMove           int                    `json:"perform_move_num_cards_to_move,omitempty"`
		ListGamesOffset                     int32                  `json:"list_games_offset,omitempty"`
		ListGamesLimit                      int32                  `json:"list_games_limit,omitempty"`
		Extra                               map[string]interface{} `json:"extra,omitempty"`
	}{
		Message:                             v.Message,
		NumBytes:                            len(v.Bytes),
		StartGameVariantMaxTimesThroughDeck: v.StartGame.Variant.MaxTimesThroughDeck,
		PerformMoveGameID:                   v.PerformMove.SavedGameDetail.GameID,
		PerformMoveNumCardsToMove:           len(v.PerformMove.Next),
		ListGamesOffset:                     v.ListGames.Cursor.Offset,
		ListGamesLimit:                      v.ListGames.Cursor.Limit,
		Extra:                               extra,
	})
}
