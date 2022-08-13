package repository

import (
	"io"

	"github.com/nicopozo/mockserver/internal/utils/log"
)

func closeBody(body io.ReadCloser, repository RuleRepository, logger log.ILogger) {
	if err := body.Close(); err != nil {
		logger.Error(repository, nil, err, "error closing response body")
	}
}
