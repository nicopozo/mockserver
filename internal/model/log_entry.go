package model

import (
	"time"
)

// LogEntry represents a captured request/response pair from the mock endpoint.
type LogEntry struct {
	ID              string            `json:"id" db:"id"`
	Timestamp       time.Time         `json:"timestamp" db:"timestamp"`
	Method          string            `json:"method" db:"method"`
	URL             string            `json:"url" db:"url"`
	RequestBody     string            `json:"request_body" db:"request_body"`
	RequestHeaders  map[string]string `json:"request_headers" db:"-"`
	QueryParams     map[string]string `json:"query_params" db:"-"`
	ResponseStatus  int               `json:"response_status" db:"response_status"`
	ResponseBody    string            `json:"response_body" db:"response_body"`
	AssertionErrors []string          `json:"assertion_errors,omitempty" db:"-"`

	// Helper fields for SQL persistence
	RawHeaders     *string `json:"-" db:"request_headers"`
	RawQueryParams *string `json:"-" db:"query_params"`
	RawAssertions  *string `json:"-" db:"assertion_errors"`
}

// LogList wraps a slice of LogEntry for API responses.
type LogList struct {
	Paging  Paging     `json:"paging"`
	Results []LogEntry `json:"results"`
}
