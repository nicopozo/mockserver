package main

import (
	"log"
	"net/http"
	"os"

	httputils "github.com/nicopozo/mockserver/internal/utils/http"
)

// main.
func main() {
	mux := http.NewServeMux()

	api := BuildContainer()

	mapRoutes(mux, api)

	// Build handler chain
	var handler http.Handler = mux

	// Apply Recovery middleware
	handler = httputils.Recovery(handler)

	// Conditionally apply CORS middleware
	if os.Getenv("MOCKS_MODE") != "release" {
		handler = httputils.CORS(handler)
	}

	log.Println("Starting server on :8080")

	if err := http.ListenAndServe(":8080", handler); err != nil { //nolint:gosec
		panic(err.Error())
	}
}
