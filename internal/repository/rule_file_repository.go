package repository

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	guuid "github.com/google/uuid"
	mockscontext "github.com/nicopozo/mockserver/internal/context"
	mockserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	jsonutils "github.com/nicopozo/mockserver/internal/utils/json"
)

const defaultFilePath = "/tmp/mocks.json"

type ruleFileRepository struct {
	rules    []fileRule
	filePath string
}

type fileRule struct {
	model.Rule
	NextResponseIndex int `json:"next_response_index"`
}

func NewRuleFileRepository(filePath string) (RuleRepository, error) {
	if filePath == "" {
		filePath = defaultFilePath
	}

	repo := ruleFileRepository{
		rules:    make([]fileRule, 0),
		filePath: filePath,
	}

	file, err := os.Open(filePath)
	if err != nil {
		errMsg := err.Error()
		if !strings.Contains(errMsg, "no such file or directory") {
			return nil, fmt.Errorf("error creating file repository when reading file: %s - %w", filePath, err)
		}

		err = repo.SaveFile()
		if err != nil {
			return nil, fmt.Errorf("error creating file repository when creating file: %s - %w", filePath, err)
		}

		file, err = os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("error creating file repository when reading file: %s - %w", filePath, err)
		}
	}

	defer func(f *os.File) { _ = f.Close() }(file)

	repo.rules = make([]fileRule, 0)

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("error reading file: %s - %w", filePath, err)
	}

	if stat.Size() > 0 {
		rules, err := unmarshalRules(file)
		if err != nil {
			return nil, fmt.Errorf("error creating file repository when unmarshaling file: %s - %w", filePath, err)
		}

		repo.rules = rules
	}

	return &repo, nil
}

func (repository *ruleFileRepository) Create(ctx context.Context, rule *model.Rule) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(repository, nil, "Saving new rule into file")

	rule.Key = fmt.Sprintf("%v", guuid.New())

	fRule := fileRule{
		Rule:              *rule,
		NextResponseIndex: rule.NextResponseIndex,
	}

	repository.rules = append(repository.rules, fRule)

	return rule, repository.SaveFile()
}

func (repository *ruleFileRepository) Update(ctx context.Context, rule *model.Rule) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(repository, nil, "Updating rule.")

	for index := range repository.rules {
		if repository.rules[index].Key == rule.Key {
			repository.rules[index] = fileRule{
				Rule:              *rule,
				NextResponseIndex: rule.NextResponseIndex,
			}
		}
	}

	return rule, repository.SaveFile()
}

func (repository *ruleFileRepository) Get(ctx context.Context, key string) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(repository, nil, "Updating rule.")

	for index := range repository.rules {
		if repository.rules[index].Key == key {
			result := &repository.rules[index].Rule
			result.NextResponseIndex = repository.rules[index].NextResponseIndex

			return result, nil
		}
	}

	msg := fmt.Sprintf("no rule found with key: %s", key)

	err := mockserrors.RuleNotFoundError{
		Message: msg,
	}

	logger.Error(repository, nil, err, msg)

	return nil, err
}

func (repository *ruleFileRepository) Search(ctx context.Context, params map[string]interface{},
	paging model.Paging,
) (*model.RuleList, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(repository, nil, "Searching rules.")

	ruleList := new(model.RuleList)
	ruleList.Paging = paging
	ruleList.Results = make([]*model.Rule, 0)

	filtered := make([]*model.Rule, 0)

	for index := range repository.rules {
		if applies(repository.rules[index].Rule, params) {
			result := &repository.rules[index].Rule
			result.NextResponseIndex = repository.rules[index].NextResponseIndex

			filtered = append(filtered, result)
		}
	}

	if paging.Offset > int32(len(filtered)) {
		return ruleList, nil
	}

	to := paging.Offset + paging.Limit
	if to > int32(len(filtered)) {
		to = int32(len(filtered))
	}

	ruleList.Results = filtered[paging.Offset:to]

	ruleList.Paging.Total = int64(len(filtered))

	return ruleList, nil
}

//nolint:cyclop
func applies(rule model.Rule, params map[string]any) bool {
	result := true

	for key, value := range params {
		formattedValue := strings.ToLower(fmt.Sprintf("%v", value))

		switch key {
		case "group":
			result = result && strings.Contains(strings.ToLower(rule.Group), formattedValue)
		case "status":
			result = result && strings.Contains(strings.ToLower(rule.Status), formattedValue)
		case "method":
			result = result && strings.Contains(strings.ToLower(rule.Method), formattedValue)
		case "strategy":
			result = result && strings.Contains(strings.ToLower(rule.Strategy), formattedValue)
		case "path":
			result = result && strings.Contains(strings.ToLower(rule.Path), formattedValue)
		case "name":
			result = result && strings.Contains(strings.ToLower(rule.Name), formattedValue)
		case "key":
			result = result && strings.Contains(strings.ToLower(rule.Key), formattedValue)
		}
	}

	return result
}

func (repository *ruleFileRepository) Delete(ctx context.Context, key string) error {
	logger := mockscontext.Logger(ctx)

	logger.Debug(repository, nil, "Updating rule.")

	var result []fileRule

	for index := range repository.rules {
		if repository.rules[index].Key == key {
			switch index {
			case 0:
				result = repository.rules[index+1 : len(repository.rules)]
			case len(repository.rules) - 1:
				result = repository.rules[0:index]
			default:
				result = repository.rules[0:index]
				result = append(result, repository.rules[index+1:len(repository.rules)]...)
			}

			repository.rules = result

			break
		}
	}

	return repository.SaveFile()
}

func (repository *ruleFileRepository) SearchByMethodAndPath(ctx context.Context, method string,
	path string,
) (*model.Rule, error) {
	logger := mockscontext.Logger(ctx)

	logger.Debug(repository, nil, "Searching by method and path rule.")

	for _, rule := range repository.rules {
		expr := CreateExpression(rule.Path)
		regex := regexp.MustCompile(expr)

		if rule.Method == method && rule.Status == model.RuleStatusEnabled && regex.MatchString(path) {
			result := rule.Rule
			result.NextResponseIndex = rule.NextResponseIndex

			return &result, nil
		}
	}

	return nil, mockserrors.RuleNotFoundError{
		Message: fmt.Sprintf("no rule found for path: %s and method %s", path, method),
	}
}

func (repository *ruleFileRepository) SaveFile() error {
	file, err := os.Create(repository.filePath)
	if err != nil {
		return fmt.Errorf("error saving file: %s - %w", repository.filePath, err)
	}

	defer func(f *os.File) { _ = f.Close() }(file)

	content := jsonutils.Marshal(repository.rules)

	_, err = file.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("error writing file content: %s - %w", repository.filePath, err)
	}

	return nil
}

func unmarshalRules(content io.Reader) ([]fileRule, error) {
	rules := new([]fileRule)

	err := jsonutils.Unmarshal(content, rules)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling content, %w", err)
	}

	return *rules, nil
}
