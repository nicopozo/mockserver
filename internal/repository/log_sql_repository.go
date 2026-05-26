package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nicopozo/mockserver/internal/model"
	jsonutils "github.com/nicopozo/mockserver/internal/utils/json"
	"github.com/oklog/ulid/v2"
)

type logSQLRepository struct {
	db Database
}

type LogRow struct {
	ID              string    `db:"id"`
	Timestamp       time.Time `db:"timestamp"`
	Method          string    `db:"method"`
	URL             string    `db:"url"`
	RequestBody     string    `db:"request_body"`
	ResponseStatus  int       `db:"response_status"`
	ResponseBody    string    `db:"response_body"`
	RequestHeaders  string    `db:"request_headers"`
	QueryParams     string    `db:"query_params"`
	AssertionErrors string    `db:"assertion_errors"`
	WebhookResults  *string   `db:"webhook_results"`
}

func rowToLogEntry(row LogRow) model.LogEntry {
	entry := model.LogEntry{
		ID:             row.ID,
		Timestamp:      row.Timestamp,
		Method:         row.Method,
		URL:            row.URL,
		RequestBody:    row.RequestBody,
		ResponseStatus: row.ResponseStatus,
		ResponseBody:   row.ResponseBody,
	}

	if row.RequestHeaders != "" {
		_ = jsonutils.Unmarshal(strings.NewReader(row.RequestHeaders), &entry.RequestHeaders)
	}

	if row.QueryParams != "" {
		_ = jsonutils.Unmarshal(strings.NewReader(row.QueryParams), &entry.QueryParams)
	}

	if row.AssertionErrors != "" {
		_ = jsonutils.Unmarshal(strings.NewReader(row.AssertionErrors), &entry.AssertionErrors)
	}

	if row.WebhookResults != nil && *row.WebhookResults != "" {
		_ = jsonutils.Unmarshal(strings.NewReader(*row.WebhookResults), &entry.WebhookResults)
	}

	return entry
}

func NewLogSQLRepository(db Database) LogRepository {
	return &logSQLRepository{
		db: db,
	}
}

func (r *logSQLRepository) Add(ctx context.Context, entry model.LogEntry) error {
	if entry.ID == "" {
		entry.ID = ulid.Make().String()
	}

	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	rawHeaders := jsonutils.Marshal(entry.RequestHeaders)
	rawParams := jsonutils.Marshal(entry.QueryParams)
	rawAssertions := jsonutils.Marshal(entry.AssertionErrors)
	webhookResultsJSON := jsonutils.Marshal(entry.WebhookResults)

	query := FormatQuery(
		"INSERT INTO request_logs (id, timestamp, method, url, request_body, "+
			"request_headers, query_params, response_status, response_body, assertion_errors, webhook_results) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		r.db.DriverName(),
	)

	_, err := r.db.Exec(query, entry.ID, entry.Timestamp, entry.Method, entry.URL, entry.RequestBody,
		rawHeaders, rawParams, entry.ResponseStatus, entry.ResponseBody, rawAssertions, webhookResultsJSON)
	if err != nil {
		return fmt.Errorf("error inserting log into DB: %w", err)
	}

	return nil
}

func (r *logSQLRepository) GetAll(ctx context.Context, paging model.Paging) (model.LogList, error) {
	var rows []LogRow

	// Get total count
	var total int64

	countQuery := FormatQuery("SELECT COUNT(*) FROM request_logs", r.db.DriverName())

	err := r.db.Get(&total, countQuery)
	if err != nil {
		return model.LogList{}, fmt.Errorf("error counting logs in DB: %w", err)
	}

	paging.Total = total

	// Get paginated results
	var query string

	var errSelect error

	if paging.LastID != "" {
		query = FormatQuery("SELECT * FROM request_logs WHERE id < ? ORDER BY id DESC LIMIT ?", r.db.DriverName())
		errSelect = r.db.Select(&rows, query, paging.LastID, paging.Limit)
	} else {
		query = FormatQuery("SELECT * FROM request_logs ORDER BY id DESC LIMIT ? OFFSET ?", r.db.DriverName())
		errSelect = r.db.Select(&rows, query, paging.Limit, paging.Offset)
	}

	if errSelect != nil {
		return model.LogList{}, fmt.Errorf("error fetching logs from DB: %w", errSelect)
	}

	results := make([]model.LogEntry, 0, len(rows))
	for _, row := range rows {
		results = append(results, rowToLogEntry(row))
	}

	return model.LogList{
		Results: results,
		Paging:  paging,
	}, nil
}

func (r *logSQLRepository) Update(ctx context.Context, logID string, updater func(entry *model.LogEntry)) error {
	// Fetch existing entry
	query := FormatQuery("SELECT * FROM request_logs WHERE id = ?", r.db.DriverName())
	row := LogRow{}

	err := r.db.Get(&row, query, logID)
	if err != nil {
		return fmt.Errorf("error fetching log entry for update: %w", err)
	}

	entry := rowToLogEntry(row)

	// Apply the updater function
	updater(&entry)

	// Serialize webhook results back
	webhookResultsJSON := jsonutils.Marshal(entry.WebhookResults)

	updateQuery := FormatQuery(
		"UPDATE request_logs SET webhook_results = ? WHERE id = ?",
		r.db.DriverName(),
	)

	_, err = r.db.Exec(updateQuery, webhookResultsJSON, logID)
	if err != nil {
		return fmt.Errorf("error updating webhook_results in DB: %w", err)
	}

	return nil
}

func (r *logSQLRepository) Clear(ctx context.Context) error {
	query := FormatQuery("DELETE FROM request_logs", r.db.DriverName())

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("error clearing logs from DB: %w", err)
	}

	return nil
}
