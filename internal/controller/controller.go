package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nicopozo/mockserver/internal/model"
)

const (
	defaultPageSize = 30
)

var errPagingLimitZero = errors.New("error parsing paging limit: limit must be greater than 0")

func getPagingFromRequest(request *http.Request) (*model.Paging, error) {
	paging := &model.Paging{
		Total:  0,
		Limit:  defaultPageSize,
		Offset: 0,
	}

	offset := request.URL.Query().Get("offset")
	if offset != "" {
		o, err := strconv.ParseInt(offset, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing paging offset, %w", err)
		}

		paging.Offset = int32(o)
	}

	limit := request.URL.Query().Get("limit")
	if limit != "" {
		parsedLimit, err := strconv.ParseInt(limit, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing paging limit, %w", err)
		}

		if parsedLimit < 0 {
			parsedLimit = 999999
		} else if parsedLimit == 0 {
			return nil, errPagingLimitZero
		}

		paging.Limit = int32(parsedLimit)
	}

	paging.LastID = request.URL.Query().Get("last_id")

	return paging, nil
}

func getParametersFromRequest(request *http.Request) map[string]any {
	queryParams := request.URL.Query()
	params := make(map[string]any, len(queryParams))

	for key, values := range queryParams {
		if key != "offset" && key != "limit" && key != "last_id" {
			params[key] = values[0]
		}
	}

	return params
}
