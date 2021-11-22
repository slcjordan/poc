package stdlib

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/slcjordan/poc/logger"
)

// Use JSON structured logging
var JSON bool

type withDepth struct {
	Logger *log.Logger
	Depth  int
}

func (w withDepth) Printf(ctx context.Context, format string, a ...interface{}) {
	if JSON {
		var val logger.Values
		prev, ok := ctx.Value(logger.ContextKey).(logger.Values)
		if ok {
			val = prev
		}
		val.Message = fmt.Sprintf(format, a...)
		encoder := json.NewEncoder(w.Logger.Writer())
		err := encoder.Encode(val)
		if err != nil {
			panic(err)
		}
	} else {
		err := w.Logger.Output(w.Depth, fmt.Sprintf(format, a...))
		if err != nil {
			panic(err)
		}
	}
}

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
