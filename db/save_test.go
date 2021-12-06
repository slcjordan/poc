package db_test

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"

	"github.com/slcjordan/poc"
	"github.com/slcjordan/poc/db"
	"github.com/slcjordan/poc/test/assert"
	"github.com/slcjordan/poc/test/harness"
	"github.com/slcjordan/poc/test/logger"
	"github.com/slcjordan/poc/test/mocks"
)

type mockRow struct {
	err error
}

func (m mockRow) Scan(dest ...interface{}) error {
	return m.err
}

func NewSaveTestPool(t *testing.T, err error) *mocks.MockPool {
	ctrl := gomock.NewController(t)
	pool := mocks.NewMockPool(ctrl)
	conn := mocks.NewMockConn(ctrl)
	conn.EXPECT().QueryRow(
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
	).Return(mockRow{err})
	conn.EXPECT().Release()
	pool.
		EXPECT().
		Acquire(gomock.Any()).
		Return(conn, nil)
	return pool
}

func TestSave(t *testing.T) {
	logger.RegisterVerbose(t)
	harness.StartGame{
		{
			Desc: "sanity check",
			Command: &db.Save{
				NewSaveTestPool(t, nil),
			},
			Result: assert.New().NoError(),
		},
		{
			Desc: "unknown query error",
			Command: &db.Save{
				NewSaveTestPool(t, errors.New("check that this error correctly causes the game to not be saved")),
			},
			Result: assert.New().Error.Category.Uint8(assert.Equals(poc.UnknownError)),
		},
	}.Run(t)
}
