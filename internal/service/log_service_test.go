package service_test

import (
	"testing"
	"time"

	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/repository"
	"github.com/nicopozo/mockserver/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestLogService_RaceCondition(t *testing.T) {
	// Create a log service with a real memory repository.
	repo := repository.NewLogMemoryRepository()
	logSvc := service.NewLogService(repo)

	logID := "test-log-id"

	// 1. Webhook completes first and calls Update (the log entry doesn't exist yet).
	webhookResult := model.WebhookResult{
		URL:          "https://example.com/callback",
		Method:       "POST",
		StatusCode:   200,
		DurationMs:   150,
		ResponseBody: `{"status":"success"}`,
	}

	logSvc.Update(logID, func(entry *model.LogEntry) {
		entry.WebhookResults = append(entry.WebhookResults, webhookResult)
	})

	// 2. Main thread completes and calls Add.
	logEntry := model.LogEntry{
		ID:             logID,
		Timestamp:      time.Now(),
		Method:         "GET",
		URL:            "/api/test",
		RequestBody:    "hello",
		ResponseStatus: 200,
		ResponseBody:   "world",
	}

	savedID := logSvc.Add(logEntry)
	assert.Equal(t, logID, savedID)

	// 3. Fetch all logs and verify that the webhook results were correctly merged!
	logs := logSvc.GetAll(model.Paging{Limit: 10})
	assert.Len(t, logs.Results, 1)

	savedEntry := logs.Results[0]
	assert.Equal(t, logID, savedEntry.ID)
	assert.Equal(t, "GET", savedEntry.Method)
	assert.Equal(t, "/api/test", savedEntry.URL)
	assert.Len(t, savedEntry.WebhookResults, 1)

	savedWebhookResult := savedEntry.WebhookResults[0]
	assert.Equal(t, webhookResult.URL, savedWebhookResult.URL)
	assert.Equal(t, webhookResult.Method, savedWebhookResult.Method)
	assert.Equal(t, webhookResult.StatusCode, savedWebhookResult.StatusCode)
	assert.Equal(t, webhookResult.DurationMs, savedWebhookResult.DurationMs)
	assert.Equal(t, webhookResult.ResponseBody, savedWebhookResult.ResponseBody)
}

func TestLogService_NormalFlow(t *testing.T) {
	// Create a log service with a real memory repository.
	repo := repository.NewLogMemoryRepository()
	logSvc := service.NewLogService(repo)

	logID := "test-log-id-normal"

	// 1. Main thread completes and calls Add.
	logEntry := model.LogEntry{
		ID:             logID,
		Timestamp:      time.Now(),
		Method:         "GET",
		URL:            "/api/test",
		RequestBody:    "hello",
		ResponseStatus: 200,
		ResponseBody:   "world",
	}

	savedID := logSvc.Add(logEntry)
	assert.Equal(t, logID, savedID)

	// 2. Webhook completes and calls Update.
	webhookResult := model.WebhookResult{
		URL:          "https://example.com/callback",
		Method:       "POST",
		StatusCode:   200,
		DurationMs:   150,
		ResponseBody: `{"status":"success"}`,
	}

	logSvc.Update(logID, func(entry *model.LogEntry) {
		entry.WebhookResults = append(entry.WebhookResults, webhookResult)
	})

	// 3. Fetch all logs and verify everything is correct.
	logs := logSvc.GetAll(model.Paging{Limit: 10})
	assert.Len(t, logs.Results, 1)

	savedEntry := logs.Results[0]
	assert.Equal(t, logID, savedEntry.ID)
	assert.Len(t, savedEntry.WebhookResults, 1)

	savedWebhookResult := savedEntry.WebhookResults[0]
	assert.Equal(t, webhookResult.URL, savedWebhookResult.URL)
	assert.Equal(t, webhookResult.ResponseBody, savedWebhookResult.ResponseBody)
}
