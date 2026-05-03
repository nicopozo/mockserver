package service

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/nicopozo/mockserver/internal/model"
)

const maxLogEntries = 500

// LogService stores request/response log entries in memory.
type LogService interface {
	Add(entry model.LogEntry)
	GetAll() model.LogList
	Clear()
}

type logService struct {
	mu      sync.RWMutex
	entries []model.LogEntry
	counter int64
}

// NewLogService creates a new in-memory LogService.
func NewLogService() LogService {
	return &logService{
		entries: make([]model.LogEntry, 0, maxLogEntries),
	}
}

func (s *logService) Add(entry model.LogEntry) {
	entry.ID = atomic.AddInt64(&s.counter, 1)
	entry.Timestamp = time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries = append(s.entries, entry)

	// Keep only the last maxLogEntries entries.
	if len(s.entries) > maxLogEntries {
		s.entries = s.entries[len(s.entries)-maxLogEntries:]
	}
}

func (s *logService) GetAll() model.LogList {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy in reverse order (newest first).
	result := make([]model.LogEntry, len(s.entries))
	for i, entry := range s.entries {
		result[len(s.entries)-1-i] = entry
	}

	return model.LogList{
		Results: result,
		Total:   len(result),
	}
}

func (s *logService) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries = make([]model.LogEntry, 0, maxLogEntries)
}
