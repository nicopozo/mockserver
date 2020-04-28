package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/utils/log"
)

const (
	defaultPageSize = 30
)

func GetLogger(c *gin.Context) log.ILogger {
	trackingID := c.GetHeader("x-tracking-id")
	if len(trackingID) == 0 {
		return log.DefaultLogger()
	}

	return log.NewLogger(trackingID)
}

func GetPagingFromRequest(request *http.Request) (*model.Paging, error) {
	paging := &model.Paging{Limit: defaultPageSize}

	offset := request.URL.Query().Get("offset")
	if offset != "" {
		o, err := strconv.ParseInt(offset, 10, 32)
		if err != nil {
			return nil, err
		}

		paging.Offset = int32(o)
	}

	limit := request.URL.Query().Get("limit")
	if limit != "" {
		l, err := strconv.ParseInt(limit, 10, 32)
		if err != nil {
			return nil, err
		}

		paging.Limit = int32(l)
	}

	return paging, nil
}

func GetParametersFromRequest(request *http.Request) map[string]interface{} {
	queryParams := request.URL.Query()
	params := make(map[string]interface{}, len(queryParams))

	for key, values := range queryParams {
		if key != "offset" && key != "limit" {
			params[key] = values[0]
		}
	}

	return params
}

func endSegment(controller interface{}, segment *newrelic.Segment, logger log.ILogger) {
	if err := segment.End(); err != nil {
		logger.Info(controller, nil, "error closing datadog segment")
	}
}
