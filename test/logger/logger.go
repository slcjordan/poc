package logger

import (
	"context"
	"testing"

	"github.com/slcjordan/poc/logger"
)

type verbose struct {
	T *testing.T
}

func (v verbose) Printf(ctx context.Context, format string, a ...interface{}) {
	v.T.Logf(format, a...)
}

// RegisterVerbose logs using testing.T if -test.v set.
func RegisterVerbose(t *testing.T) {
	logger.RegisterInfo(verbose{t})
	logger.RegisterError(verbose{t})
}
