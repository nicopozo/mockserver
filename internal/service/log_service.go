package service

import (
	"context"
	"sync"

	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/repository"
	"github.com/oklog/ulid/v2"
)

// LogService stores request/response log entries.
type LogService interface {
	Add(entry model.LogEntry) string
	Update(id string, updater func(entry *model.LogEntry))
	GetAll(paging model.Paging) model.LogList
	Clear()
}

type logService struct {
	repo                  repository.LogRepository
	pendingWebhookResults sync.Map // key: logID (string), value: []model.WebhookResult
}

// NewLogService creates a new LogService with the provided repository.
func NewLogService(repo repository.LogRepository) LogService {
	return &logService{
		repo: repo,
	}
}

func (s *logService) Add(entry model.LogEntry) string {
	if entry.ID == "" {
		entry.ID = ulid.Make().String()
	}

	// Merge any pending webhook results that arrived before Add was called
	val, ok := s.pendingWebhookResults.Load(entry.ID)
	if ok {
		results, isResults := val.([]model.WebhookResult)
		if isResults {
			entry.WebhookResults = append(entry.WebhookResults, results...)
		}

		s.pendingWebhookResults.Delete(entry.ID)
	}

	_ = s.repo.Add(context.Background(), entry)

	return entry.ID
}

func (s *logService) Update(logID string, updater func(entry *model.LogEntry)) {
	err := s.repo.Update(context.Background(), logID, updater)
	if err != nil {
		s.handlePendingWebhook(logID, updater)
	}
}

func (s *logService) handlePendingWebhook(logID string, updater func(entry *model.LogEntry)) {
	// Entry not found yet (race condition where webhook finishes before Add is called).
	// We execute the updater on a temporary entry to collect the webhook results,
	// and save them in the pending map.
	var tempEntry model.LogEntry

	updater(&tempEntry)

	if len(tempEntry.WebhookResults) == 0 {
		return
	}

	var newResults []model.WebhookResult

	val, ok := s.pendingWebhookResults.Load(logID)
	if ok {
		results, isResults := val.([]model.WebhookResult)
		if isResults {
			newResults = append(newResults, results...)
		}
	}

	newResults = append(newResults, tempEntry.WebhookResults...)

	s.pendingWebhookResults.Store(logID, newResults)
}

func (s *logService) GetAll(paging model.Paging) model.LogList {
	logs, err := s.repo.GetAll(context.Background(), paging)
	if err != nil {
		return model.LogList{
			Results: []model.LogEntry{},
			Paging:  paging,
		}
	}

	return logs
}

func (s *logService) Clear() {
	_ = s.repo.Clear(context.Background())
}
