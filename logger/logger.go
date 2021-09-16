package logger

import (
	"context"
	"os"
)

type Logger interface {
	Printf(context.Context, string, ...interface{})
}

var infos []Logger
var errors []Logger

func RegisterInfo(l Logger) {
	infos = append(infos, l)
}

func RegisterError(l Logger) {
	errors = append(errors, l)
}

func Infof(ctx context.Context, format string, a ...interface{}) {
	for _, l := range infos {
		l.Printf(ctx, format+"\n", a...)
	}
}

func Errorf(ctx context.Context, format string, a ...interface{}) {
	for _, l := range errors {
		l.Printf(ctx, format+"\n", a...)
	}
}

func Fatalf(ctx context.Context, format string, a ...interface{}) {
	for _, l := range errors {
		l.Printf(ctx, format+"\n", a...)
	}
	os.Exit(1)
}
