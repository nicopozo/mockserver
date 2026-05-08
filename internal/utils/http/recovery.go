package httputils

import (
	"net/http"

	mockscontext "github.com/nicopozo/mockserver/internal/context"
	"github.com/nicopozo/mockserver/internal/model"
)

// Recovery is a middleware that recovers from panics and returns a 500 Internal Server Error.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				ctx := mockscontext.New(request)
				logger := mockscontext.Logger(ctx)

				logger.Error(nil, nil, nil, "Panic recovered: %v", err)

				WriteError(writer, model.InternalError, "Internal Server Error")
			}
		}()

		next.ServeHTTP(writer, request)
	})
}
