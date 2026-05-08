package httputils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nicopozo/mockserver/internal/model"
)

// WriteJSON writes a JSON response with the given status code and data.
func WriteJSON(writer http.ResponseWriter, status int, data any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	if err := json.NewEncoder(writer).Encode(data); err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err) //nolint:forbidigo
	}
}

// WriteError writes a consistent error response using the model.Error structure.
func WriteError(writer http.ResponseWriter, causeCode int64, message string, args ...interface{}) {
	errResult := model.NewError(causeCode, message, args...)

	WriteJSON(writer, errResult.Status, errResult)
}
