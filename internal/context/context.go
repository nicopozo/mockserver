package mockscontext

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/nicopozo/mockserver/internal/utils/log"
)

type loggerKey struct{}

func New(c *gin.Context) context.Context {
	trackingID := c.GetHeader("x-tracking-id")
	if len(trackingID) == 0 {
		return context.WithValue(c.Request.Context(), loggerKey{}, log.DefaultLogger())
	}

	return context.WithValue(c.Request.Context(), loggerKey{}, log.NewLogger(trackingID))
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
