package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/nicopozo/mockserver/internal/model"
)

const (
	defaultPageSize = 30
)

func getPagingFromRequest(request *http.Request) (*model.Paging, error) {
	paging := &model.Paging{
		Total: 0,
		Limit: defaultPageSize,
	}

	limit := request.URL.Query().Get("limit")
	if limit != "" {
		l, err := strconv.ParseInt(limit, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing paging limit, %w", err)
		}

		paging.Limit = int32(l)
	}

	paging.LastID = request.URL.Query().Get("last_id")

	return paging, nil
}

func getParametersFromRequest(request *http.Request) map[string]interface{} {
	queryParams := request.URL.Query()
	params := make(map[string]interface{}, len(queryParams))

	for key, values := range queryParams {
		if key != "limit" && key != "last_id" {
			params[key] = values[0]
		}
	}

	return params
}
