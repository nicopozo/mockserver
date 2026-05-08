package httputils

import (
	"net/http"
)

// CORS is a simple middleware that adds CORS headers to the response.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-tracking-id")

		if request.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusNoContent)

			return
		}

		next.ServeHTTP(writer, request)
	})
}
