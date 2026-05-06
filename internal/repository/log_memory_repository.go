package repository

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/nicopozo/mockserver/internal/model"
	"github.com/oklog/ulid/v2"
)

const maxLogEntries = 500

type logMemoryRepository struct {
	mu      sync.RWMutex
	entries []model.LogEntry
}

func NewLogMemoryRepository() LogRepository {
	return &logMemoryRepository{
		entries: make([]model.LogEntry, 0, maxLogEntries),
	}
}

func (r *logMemoryRepository) Add(ctx context.Context, entry model.LogEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if entry.ID == "" {
		entry.ID = ulid.Make().String()
	}

	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	r.entries = append(r.entries, entry)

	if len(r.entries) > maxLogEntries {
		r.entries = r.entries[len(r.entries)-maxLogEntries:]
	}

	return nil
}

func (r *logMemoryRepository) GetAll(ctx context.Context, paging model.Paging) (model.LogList, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Sort a copy of all entries newest first
	allEntries := make([]model.LogEntry, len(r.entries))
	copy(allEntries, r.entries)

	sort.Slice(allEntries, func(i, j int) bool {
		return allEntries[i].ID > allEntries[j].ID
	})

	total := int64(len(allEntries))
	paging.Total = total

	var result []model.LogEntry

	start := int(paging.Offset)

	if paging.LastID != "" {
		start = findStartIndex(allEntries, paging.LastID)
	}

	if start != -1 {
		end := start + int(paging.Limit)

		if start > len(allEntries) {
			start = len(allEntries)
		}

		if end > len(allEntries) {
			end = len(allEntries)
		}

		result = allEntries[start:end]
	}

	return model.LogList{
		Results: result,
		Paging:  paging,
	}, nil
}

func findStartIndex(entries []model.LogEntry, lastID string) int {
	for index, entry := range entries {
		if entry.ID == lastID {
			return index + 1
		}
	}

	return -1
}

func (r *logMemoryRepository) Clear(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.entries = make([]model.LogEntry, 0, maxLogEntries)

	return nil
}
