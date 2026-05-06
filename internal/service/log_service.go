package service

import (
	"context"

	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/repository"
)

// LogService stores request/response log entries.
type LogService interface {
	Add(entry model.LogEntry)
	GetAll(paging model.Paging) model.LogList
	Clear()
}

type logService struct {
	repo repository.LogRepository
}

// NewLogService creates a new LogService with the provided repository.
func NewLogService(repo repository.LogRepository) LogService {
	return &logService{
		repo: repo,
	}
}

func (s *logService) Add(entry model.LogEntry) {
	_ = s.repo.Add(context.Background(), entry)
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
