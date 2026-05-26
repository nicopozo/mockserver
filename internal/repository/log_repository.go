package repository

import (
	"context"

	"github.com/nicopozo/mockserver/internal/model"
)

type LogRepository interface {
	Add(ctx context.Context, entry model.LogEntry) error
	Update(ctx context.Context, id string, updater func(entry *model.LogEntry)) error
	GetAll(ctx context.Context, paging model.Paging) (model.LogList, error)
	Clear(ctx context.Context) error
}
