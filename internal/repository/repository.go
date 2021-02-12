package repository

import (
	"io"

	"github.com/nicopozo/mockserver/internal/utils/log"
)

func closeBody(body io.ReadCloser, repository IRuleRepository, logger log.ILogger) {
	err := body.Close()
	if err != nil {
		logger.Error(repository, nil, err, "error closing response body")
	}
}
