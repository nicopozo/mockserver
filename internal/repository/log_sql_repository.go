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

	query := FormatQuery(
		"INSERT INTO request_logs (id, timestamp, method, url, request_body, "+
			"request_headers, query_params, response_status, response_body, assertion_errors) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		r.db.DriverName(),
	)

	_, err := r.db.Exec(query, entry.ID, entry.Timestamp, entry.Method, entry.URL, entry.RequestBody,
		rawHeaders, rawParams, entry.ResponseStatus, entry.ResponseBody, rawAssertions)
	if err != nil {
		return fmt.Errorf("error inserting log into DB: %w", err)
	}

	return nil
}

func (r *logSQLRepository) GetAll(ctx context.Context, paging model.Paging) (model.LogList, error) {
	var rows []model.LogEntry

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
		query = FormatQuery("SELECT * FROM request_logs ORDER BY id DESC LIMIT ?", r.db.DriverName())
		errSelect = r.db.Select(&rows, query, paging.Limit)
	}

	if errSelect != nil {
		return model.LogList{}, fmt.Errorf("error fetching logs from DB: %w", errSelect)
	}

	for index := range rows {
		if rows[index].RawHeaders != nil {
			_ = jsonutils.Unmarshal(strings.NewReader(*rows[index].RawHeaders), &rows[index].RequestHeaders)
		}

		if rows[index].RawQueryParams != nil {
			_ = jsonutils.Unmarshal(strings.NewReader(*rows[index].RawQueryParams), &rows[index].QueryParams)
		}

		if rows[index].RawAssertions != nil {
			_ = jsonutils.Unmarshal(strings.NewReader(*rows[index].RawAssertions), &rows[index].AssertionErrors)
		}
	}

	return model.LogList{
		Results: rows,
		Paging:  paging,
	}, nil
}

func (r *logSQLRepository) Clear(ctx context.Context) error {
	query := FormatQuery("DELETE FROM request_logs", r.db.DriverName())

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("error clearing logs from DB: %w", err)
	}

	return nil
}
