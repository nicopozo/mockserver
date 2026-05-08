package main

import (
	"net/http"
)

func mapRoutes(mux *http.ServeMux, api MockContainer) {
	// Serve Admin Web UI
	mux.Handle("/mock-service/admin/", http.StripPrefix("/mock-service/admin/", http.FileServer(http.Dir("web/dist"))))

	ruleController := api.Controllers.RuleController
	mux.HandleFunc("POST /mock-service/rules", ruleController.Create)
	mux.HandleFunc("GET /mock-service/rules/{key}", ruleController.Get)
	mux.HandleFunc("GET /mock-service/rules", ruleController.Search)
	mux.HandleFunc("DELETE /mock-service/rules/{key}", ruleController.Delete)
	mux.HandleFunc("PUT /mock-service/rules/{key}", ruleController.Update)
	mux.HandleFunc("PUT /mock-service/rules/{key}/status", ruleController.UpdateStatus)
	mux.HandleFunc("GET /mock-service/rules/export", ruleController.Export)
	mux.HandleFunc("POST /mock-service/rules/import", ruleController.Import)

	mockController := api.Controllers.MockController
	// Any method wildcard route
	mux.HandleFunc("/mock-service/mock/{rule...}", mockController.Execute)

	logController := api.Controllers.LogController
	mux.HandleFunc("GET /mock-service/logs", logController.GetLogs)
	mux.HandleFunc("DELETE /mock-service/logs", logController.ClearLogs)

	mux.HandleFunc("GET /ping", ping)
}

func ping(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("pong"))
}
