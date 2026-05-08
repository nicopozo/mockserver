package mockscontext

import (
	"context"
	"net/http"

	"github.com/nicopozo/mockserver/internal/utils/log"
)

type loggerKey struct{}

func New(request *http.Request) context.Context {
	trackingID := request.Header.Get("x-tracking-id")

	if len(trackingID) == 0 {
		return context.WithValue(request.Context(), loggerKey{}, log.DefaultLogger())
	}

	return context.WithValue(request.Context(), loggerKey{}, log.NewLogger(trackingID))
}

func Logger(ctx context.Context) log.ILogger {
	logger, ok := ctx.Value(loggerKey{}).(log.ILogger)
	if !ok {
		return log.DefaultLogger()
	}

	return logger
}

func Background() context.Context {
	return context.WithValue(context.Background(), loggerKey{}, log.DefaultLogger())
}
