package repository

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/nicopozo/mockserver/internal/model"
)

//nolint:lll
//go:generate mockgen -destination=../utils/test/mocks/rule_repository_mock.go -package=mocks -source=./rule_repository.go

type IRuleRepository interface {
	Create(ctx context.Context, rule *model.Rule) (*model.Rule, error)
	Update(ctx context.Context, rule *model.Rule) (*model.Rule, error)
	Get(ctx context.Context, key string) (*model.Rule, error)
	Search(ctx context.Context, params map[string]interface{}, paging model.Paging) (*model.RuleList, error)
	SearchByMethodAndPath(ctx context.Context, method string, path string) (*model.Rule, error)
	Delete(ctx context.Context, key string) error
}

func CreateExpression(path string) string {
	paramRegex := regexp.MustCompile("{.+?}/")
	params := paramRegex.FindAllString(path, -1)

	for _, param := range params {
		path = strings.ReplaceAll(path, param, "[^/]+?/")
	}

	paramRegex = regexp.MustCompile("{.+?}")
	params = paramRegex.FindAllString(path, -1)

	for _, param := range params {
		path = strings.ReplaceAll(path, param, "[^/]+")
	}

	return fmt.Sprintf("^%s$", path)
}
