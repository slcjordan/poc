package logger

import (
	"context"
)

// A Logger accepts a format string and arguments.
type Logger interface {
	Printf(context.Context, string, ...interface{})
}

var infos []Logger
var errors []Logger

// RegisterInfo may be called by a logging package's init function.
func RegisterInfo(l Logger) {
	infos = append(infos, l)
}

// RegisterError may be called by a logging package's init function.
func RegisterError(l Logger) {
	errors = append(errors, l)
}

// Infof may be used to log error conditions that are the fault of external applications.
func Infof(ctx context.Context, format string, a ...interface{}) {
	for _, l := range infos {
		l.Printf(ctx, format+"\n", a...)
	}
}

// Errorf must be used to log error conditions that may be caused by errors in this application.
func Errorf(ctx context.Context, format string, a ...interface{}) {
	for _, l := range errors {
		l.Printf(ctx, format+"\n", a...)
	}
}
