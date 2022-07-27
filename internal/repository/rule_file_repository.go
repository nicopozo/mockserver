package repository

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	guuid "github.com/google/uuid"
	mockscontext "github.com/nicopozo/mockserver/internal/context"
	mockserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	jsonutils "github.com/nicopozo/mockserver/internal/utils/json"
)

type RuleFileRepository struct {
	rules    []model.Rule
	filePath string
}

func NewRuleFileRepository(filePath string) (IRuleRepository, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error creating file repository when reading file: %s - %w", filePath, err)
	}

	defer func(f *os.File) { _ = file.Close() }(file)

	rules, err := model.UnmarshalRules(file)
	if err != nil {
		return nil, fmt.Errorf("error creating file repository when unmarshaling file: %s - %w", filePath, err)
	}

	return &RuleFileRepository{
		rules:    rules,
		filePath: filePath,
	}, nil
}

func (repository *RuleFileRepository) Create(ctx context.Context, rule *model.Rule) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(repository, nil, "Saving new rule into file")

	rule.Key = fmt.Sprintf("%v", guuid.New())

	repository.rules = append(repository.rules, *rule)

	return rule, repository.SaveFile(ctx)
}

func (repository *RuleFileRepository) Update(ctx context.Context, rule *model.Rule) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(repository, nil, "Updating rule.")

	for index := range repository.rules {
		if repository.rules[index].Key == rule.Key {
			repository.rules[index] = *rule
		}
	}

	return rule, repository.SaveFile(ctx)
}

func (repository *RuleFileRepository) Get(ctx context.Context, key string) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(repository, nil, "Updating rule.")

	for index := range repository.rules {
		if repository.rules[index].Key == key {
			return &repository.rules[index], nil
		}
	}

	msg := fmt.Sprintf("no rule found with key: %s", key)

	err := mockserrors.RuleNotFoundError{
		Message: msg,
	}

	logger.Error(repository, nil, err, msg)

	return nil, err
}

func (repository *RuleFileRepository) Search(ctx context.Context, params map[string]interface{},
	paging model.Paging) (*model.RuleList, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(repository, nil, "Searching rules.")

	ruleList := new(model.RuleList)
	ruleList.Paging = paging
	ruleList.Results = make([]*model.Rule, 0)

	if paging.Offset > int32(len(repository.rules)) {
		return ruleList, nil
	}

	to := paging.Offset + paging.Limit
	if to > int32(len(repository.rules)) {
		to = int32(len(repository.rules))
	}
	for index := range repository.rules[paging.Offset:to] {
		if applies(repository.rules[index], params) {
			ruleList.Results = append(ruleList.Results, &repository.rules[index])
		}
	}

	ruleList.Paging.Total = int64(len(repository.rules))

	return ruleList, nil
}

func applies(rule model.Rule, params map[string]any) bool {
	result := true

	for key, value := range params {
		v := strings.ToLower(fmt.Sprintf("%v", value))
		switch key {
		case "application":
			result = result && strings.Contains(rule.Application, v)
		case "status":
			result = result && strings.Contains(rule.Status, v)
		case "method":
			result = result && strings.Contains(rule.Method, v)
		case "strategy":
			result = result && strings.Contains(rule.Strategy, v)
		case "path":
			result = result && strings.Contains(rule.Path, v)
		case "name":
			result = result && strings.Contains(rule.Name, v)
		case "key":
			result = result && strings.Contains(rule.Key, v)
		}
	}

	return result
}

func (repository *RuleFileRepository) Delete(ctx context.Context, key string) error {
	logger := mockscontext.Logger(ctx)

	logger.Debug(repository, nil, "Updating rule.")

	for index := range repository.rules {
		if repository.rules[index].Key == key {
			result := append(repository.rules[0:index], repository.rules[index+1:len(repository.rules)]...)
			repository.rules = result
		}
	}

	return repository.SaveFile(ctx)
}

func (repository *RuleFileRepository) SearchByMethodAndPath(ctx context.Context, method string,
	path string) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(repository, nil, "Searching by method and path rule.")

	for _, rule := range repository.rules {
		var regex = regexp.MustCompile(CreateExpression(rule.Path))

		if rule.Method == method && rule.Status == model.RuleStatusEnabled && regex.MatchString(path) {
			return &rule, nil
		}
	}

	return nil, mockserrors.RuleNotFoundError{
		Message: fmt.Sprintf("no rule found for path: %s and method %s", path, method),
	}
}

func (repository *RuleFileRepository) SaveFile(ctx context.Context) error {
	file, err := os.Create(repository.filePath)
	if err != nil {
		return fmt.Errorf("error saving file: %s - %w", repository.filePath, err)
	}

	defer func(f *os.File) { _ = file.Close() }(file)

	content := jsonutils.Marshal(repository.rules)

	_, err = file.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("error writing file content: %s - %w", repository.filePath, err)
	}

	return nil
}
