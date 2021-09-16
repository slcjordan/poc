package stdlib

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/slcjordan/poc/logger"
)

type withDepth struct {
	Logger *log.Logger
	Depth  int
}

func (w withDepth) Printf(ctx context.Context, format string, a ...interface{}) {
	w.Logger.Output(w.Depth, fmt.Sprintf(format, a...))
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
