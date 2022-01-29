package router_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/slcjordan/poc"
	"github.com/slcjordan/poc/router"
	"github.com/slcjordan/poc/test/mocks"
)

func TestRouter(t *testing.T) {
	// run tests
	for _, client := range []struct {
		Method string
		Path   string
	}{
		{http.MethodPost, "/v1/game/start"},
		{http.MethodPost, "/v1/game/2021/move"},
		{http.MethodGet, "/v1/game/list"},
	} {
		for _, testCase := range []struct {
			Error poc.ErrorCategory
			Code  int
		}{
			{poc.SemanticError, http.StatusUnprocessableEntity},
			{poc.MalformedError, http.StatusBadRequest},
			{poc.UnavailableError, http.StatusServiceUnavailable},
			{poc.UnimplementedError, http.StatusNotImplemented},
			{poc.NotFoundError, http.StatusNotFound},
			{poc.UnknownError, http.StatusInternalServerError},
		} {
			t.Run(fmt.Sprintf("%s %v %s %d", client.Method, client.Path, testCase.Error, testCase.Code), checkSaveErrorCode(client.Method, client.Path, testCase.Error, testCase.Code))
		}
	}
}

func checkSaveErrorCode(method string, path string, category poc.ErrorCategory, expected int) func(*testing.T) {
	return func(t *testing.T) {
		// setup dependencies
		ctrl := gomock.NewController(t)
		command := mocks.NewMockByteCaller(ctrl)
		mux := router.New(router.V1Handlers{
			PostGameStart:    command,
			PostGameByIDMove: command,
			GetGameList:      command,
		})
		command.
			EXPECT().
			CallBytes(gomock.Any(), gomock.Any()).
			Return(
				nil,
				poc.Error{
					Actual:   errors.New("actual"),
					Category: category,
				})
		r := httptest.NewRequest(method, path, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != expected {
			t.Fatalf("expected %d for %s but got %d", expected, category, w.Code)
		}
	}
}
