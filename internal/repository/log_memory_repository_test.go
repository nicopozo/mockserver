package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestLogMemoryRepository(t *testing.T) {
	repo := repository.NewLogMemoryRepository()
	ctx := context.Background()

	t.Run("Add and GetAll", func(t *testing.T) {
		entry1 := model.LogEntry{ID: "1", Method: "GET", URL: "/test1", Timestamp: time.Now()}
		entry2 := model.LogEntry{ID: "2", Method: "POST", URL: "/test2", Timestamp: time.Now()}

		err := repo.Add(ctx, entry1)
		assert.NoError(t, err)
		err = repo.Add(ctx, entry2)
		assert.NoError(t, err)

		list, err := repo.GetAll(ctx, model.Paging{Limit: 10})
		assert.NoError(t, err)
		assert.Equal(t, int64(2), list.Paging.Total)
		assert.Len(t, list.Results, 2)
		// Newest first
		assert.Equal(t, "2", list.Results[0].ID)
		assert.Equal(t, "1", list.Results[1].ID)
	})

	t.Run("Keyset Pagination", func(t *testing.T) {
		list, err := repo.GetAll(ctx, model.Paging{Limit: 1, LastID: "2"})
		assert.NoError(t, err)
		assert.Len(t, list.Results, 1)
		assert.Equal(t, "1", list.Results[0].ID)
	})

	t.Run("Clear", func(t *testing.T) {
		err := repo.Clear(ctx)
		assert.NoError(t, err)

		list, err := repo.GetAll(ctx, model.Paging{Limit: 10})
		assert.NoError(t, err)
		assert.Equal(t, int64(0), list.Paging.Total)
		assert.Len(t, list.Results, 0)
	})
}
