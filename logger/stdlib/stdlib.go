package stdlib

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/slcjordan/poc/logger"
)

func init() {
	flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.LUTC
	logger.RegisterInfo(withDepth{
		Logger: log.New(os.Stdout, " [INFO] ", flags),
		Depth:  3,
	})
	logger.RegisterError(withDepth{
		Logger: log.New(os.Stderr, " [ERROR] ", flags),
		Depth:  3,
	})
}

//go:generate stringer -type=outputFormat
type outputFormat uint8

// supported formats
const (
	Log outputFormat = iota
	JSON
)

// Format to be used for logging. Not safe for concurrent read/writes.
var Format outputFormat

type withDepth struct {
	Logger *log.Logger
	Depth  int
}

func (w withDepth) Printf(ctx context.Context, format string, a ...interface{}) {
	switch Format {
	case JSON:
		var val logger.Values
		prev, ok := ctx.Value(logger.MiddlewareKey).(logger.Values)
		if ok {
			val = prev
		}
		val.Message = fmt.Sprintf(format, a...)
		encoder := json.NewEncoder(w.Logger.Writer())
		err := encoder.Encode(val)
		if err != nil {
			panic(err)
		}
	default:
		err := w.Logger.Output(w.Depth, fmt.Sprintf(format, a...))
		if err != nil {
			panic(err)
		}
	}
}
