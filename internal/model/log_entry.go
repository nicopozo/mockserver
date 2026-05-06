package model

import "time"

// LogEntry represents a captured request/response pair from the mock endpoint.
type LogEntry struct {
	ID              int64             `json:"id"`
	Timestamp       time.Time         `json:"timestamp"`
	Method          string            `json:"method"`
	URL             string            `json:"url"`
	RequestBody     string            `json:"request_body"`
	RequestHeaders  map[string]string `json:"request_headers"`
	QueryParams     map[string]string `json:"query_params"`
	ResponseStatus  int               `json:"response_status"`
	ResponseBody    string            `json:"response_body"`
	AssertionErrors []string          `json:"assertion_errors,omitempty"`
}

// LogList wraps a slice of LogEntry for API responses.
type LogList struct {
	Results []LogEntry `json:"results"`
	Total   int        `json:"total"`
}
