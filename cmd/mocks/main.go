package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/nicopozo/mockserver/internal/configs"
	httputils "github.com/nicopozo/mockserver/internal/utils/http"
)

func main() {
	mux := http.NewServeMux()

	cfg := configs.New()
	api := BuildContainer(cfg)

	mapRoutes(mux, api)

	// Build handler chain
	var handler http.Handler = mux

	// Apply Recovery middleware
	handler = httputils.Recovery(handler)

	// Conditionally apply CORS middleware
	if os.Getenv("MOCKS_MODE") != "release" {
		handler = httputils.CORS(handler)
	}

	if cfg.IsLambda {
		log.Println("Starting AWS Lambda handler")
		lambda.Start(httpadapter.NewV2(handler).ProxyWithContext)
	} else {
		log.Println("Starting server on :8080")

		if err := http.ListenAndServe(":8080", handler); err != nil { //nolint:gosec
			panic(err.Error())
		}
	}
}
