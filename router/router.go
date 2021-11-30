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

//go:generate go run github.com/golang/mock/mockgen -package=mocks -destination=../test/mocks/router.go -source=router.go

// A ByteCaller processes request bodies and returns response bodies.
type ByteCaller interface {
	CallBytes(context.Context, []byte) ([]byte, error)
}

// V1Handlers members must not be nil.
type V1Handlers struct {
	StartGame   ByteCaller
	PerformMove ByteCaller
	ListGames   ByteCaller
}

// New sets up routes with passed middleware.
func New(
	v1 V1Handlers,
	middlewares ...func(http.Handler) http.Handler,
) chi.Router {
	router := chi.NewRouter().With(middlewares...)
	router.Post("/v1/game/start", handlerFunc(v1.StartGame))
	router.Post(fmt.Sprintf("/v1/game/{%s}/move", gameIDKey), handlerFunc(v1.PerformMove))
	router.Get("/v1/game/list", handlerFunc(v1.ListGames))
	return router
}

func handlerFunc(f ByteCaller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeHeader := func(err error) {
			var status int
			var catErr poc.Error
			if errors.As(err, &catErr) {
				status = map[poc.ErrorCategory]int{
					poc.SemanticError:      http.StatusUnprocessableEntity,
					poc.MalformedError:     http.StatusBadRequest,
					poc.UnavailableError:   http.StatusServiceUnavailable,
					poc.UnimplementedError: http.StatusNotImplemented,
					poc.NotFoundError:      http.StatusNotFound,
					poc.UnknownError:       http.StatusInternalServerError,
				}[catErr.Category]
				if status == 0 {
					status = http.StatusInternalServerError
				}
			} else {
				status = http.StatusInternalServerError
			}
			w.WriteHeader(status)
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
	if err != nil {
		return move, poc.Error{Actual: err, Category: poc.MalformedError}
	}
	return move, nil
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
