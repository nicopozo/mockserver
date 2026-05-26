package model

import (
	"time"
)

// WebhookResult stores the outcome of a fired webhook.
type WebhookResult struct {
	URL          string `json:"url"`
	Method       string `json:"method"`
	StatusCode   int    `json:"status_code"`
	DurationMs   int64  `json:"duration_ms"`
	Error        string `json:"error,omitempty"`
	ResponseBody string `json:"response_body,omitempty"`
}

// LogEntry represents a captured request/response pair from the mock endpoint.
type LogEntry struct {
	ID              string            `json:"id"`
	Timestamp       time.Time         `json:"timestamp"`
	Method          string            `json:"method"`
	URL             string            `json:"url"`
	RequestBody     string            `json:"request_body"`
	RequestHeaders  map[string]string `json:"request_headers"`
	QueryParams     map[string]string `json:"query_params"`
	ResponseStatus  int               `json:"response_status"`
	ResponseBody    string            `json:"response_body"`
	AssertionErrors []string          `json:"assertion_errors,omitempty"`
	WebhookResults  []WebhookResult   `json:"webhook_results,omitempty"`
}

// LogList wraps a slice of LogEntry for API responses.
type LogList struct {
	Paging  Paging     `json:"paging"`
	Results []LogEntry `json:"results"`
}
