package router

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/slcjordan/poc"
	"github.com/slcjordan/poc/test"
)

func TestRouter(t *testing.T) {
	// setup dependencies
	ctrl := gomock.NewController(t)
	command := test.NewMockByteCaller(ctrl)
	mux := NewMux()
	mux.RegisterRoutesV1(V1Handlers{
		StartGame:   command,
		PerformMove: command,
		ListGames:   command,
	})

	// run tests
	t.Run("semantic error on save results in 404", checkSaveErrorCode(command, mux, "/v1/game/list", poc.SemanticError, 422))
}

func checkSaveErrorCode(command *test.MockByteCaller, mux http.Handler, path string, category poc.ErrorCategory, expected int) func(*testing.T) {
	command.
		EXPECT().
		CallBytes(gomock.Any(), gomock.Any()).
		Return(
			nil,
			poc.Error{
				Actual:   errors.New("actual"),
				Category: category,
			})
	return func(t *testing.T) {
		r := httptest.NewRequest("GET", path, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != expected {
			t.Fatalf("expected %d for %s but got %d", expected, category, w.Code)
		}
	}
}
