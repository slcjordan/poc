package router

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	chi "github.com/go-chi/chi/v5"
	"github.com/slcjordan/poc"
)

type key string

const query = key("query")

const (
	gameIDKey = "gameID"
)

// A ByteCaller processes request bodies and returns response bodies.
type ByteCaller interface {
	CallBytes(context.Context, []byte) ([]byte, error)
}

// Mux routes mux requests.
type Mux struct {
	router *chi.Mux
}

// V1Handlers members must not be nil.
type V1Handlers struct {
	StartGame   ByteCaller
	PerformMove ByteCaller
	ListGames   ByteCaller
}

// NewMux initializes a mux handler.
func NewMux() *Mux {
	return &Mux{
		router: chi.NewRouter(),
	}
}

// Use sets up mux middleware.
func (mux *Mux) Use(
	middlewares ...func(http.Handler) http.Handler,
) {
	mux.router.Use(middlewares...)
}

// RegisterRoutesV1 sets up v1 routes with passed middleware.
func (mux *Mux) RegisterRoutesV1(
	handlers V1Handlers,
	middlewares ...func(http.Handler) http.Handler,
) {
	router := mux.router.With(middlewares...)
	router.Post("/v1/game/start", handlerFunc(handlers.StartGame))
	router.Post(fmt.Sprintf("/mux/game/{%s}/move", gameIDKey), handlerFunc(handlers.PerformMove))
	router.Get("/v1/game/list", handlerFunc(handlers.ListGames))
}

func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.router.ServeHTTP(w, r)
}

func handlerFunc(f ByteCaller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeHeader := func(err error) {
			var catErr poc.Error
			if errors.As(err, &catErr) {
				switch catErr.Category {
				case poc.SemanticError:
					w.WriteHeader(http.StatusUnprocessableEntity)
				case poc.MalformedError:
					w.WriteHeader(http.StatusBadRequest)
				case poc.UnavailableError:
					w.WriteHeader(http.StatusServiceUnavailable)
				case poc.UnimplementedError:
					w.WriteHeader(http.StatusNotImplemented)
				case poc.NotFoundError:
					w.WriteHeader(http.StatusNotFound)
				case poc.UnknownError:
					w.WriteHeader(http.StatusInternalServerError)
				default:
					w.WriteHeader(http.StatusInternalServerError)
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeHeader(err)
			w.Write([]byte(err.Error()))
			return
		}
		ctx := context.WithValue(r.Context(), query, r.URL.Query())
		result, err := f.CallBytes(ctx, body)
		if err != nil {
			writeHeader(err)
			w.Write([]byte(err.Error()))
		}
		io.Copy(w, bytes.NewReader(result))
	}
}

// V1HydrateURLAndQueryParams adds url path and query params from context to data.
type V1HydrateURLAndQueryParams struct {
	OffsetKey string
	LimitKey  string
}

// CallPerformMove adds move.Input.GameID url path param.
func (params V1HydrateURLAndQueryParams) CallPerformMove(ctx context.Context, move poc.PerformMove) (poc.PerformMove, error) {
	gameID := chi.URLParamFromCtx(ctx, gameIDKey)
	var err error
	move.Input.GameID, err = strconv.ParseInt(gameID, 10, 64)
	return move, poc.Error{Actual: err, Category: poc.MalformedError}
}

// CallListGames adds list.Input.Limit and list.Input.Offset url query param.
func (params V1HydrateURLAndQueryParams) CallListGames(ctx context.Context, cursor poc.ListGames) (poc.ListGames, error) {
	values := ctx.Value(query).(url.Values)
	offset, err := strconv.ParseInt(values.Get(params.OffsetKey), 10, 32)
	if err != nil {
		return cursor, poc.Error{Actual: fmt.Errorf("parsing url param %#v: %w", params.OffsetKey, err), Category: poc.MalformedError}
	}
	limit, err := strconv.ParseInt(values.Get(params.LimitKey), 10, 32)
	if err != nil {
		return cursor, poc.Error{Actual: fmt.Errorf("parsing url param %#v: %w", params.LimitKey, err), Category: poc.MalformedError}
	}
	cursor.Input.Offset = int32(offset)
	cursor.Input.Limit = int32(limit)
	return cursor, nil
}
