package repository

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/nicopozo/mockserver/internal/configs"
	mockscontext "github.com/nicopozo/mockserver/internal/context"
	mockserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/oklog/ulid/v2"
)

type DynamoRuleRepository struct {
	client    *dynamodb.Client
	tableName string
}

// NewDynamoRuleRepository creates a new RuleRepository for DynamoDB.
func NewDynamoRuleRepository(client *dynamodb.Client, cfg *configs.Config) RuleRepository {
	return &DynamoRuleRepository{
		client:    client,
		tableName: cfg.Dynamo.TablePrefix + "rules",
	}
}

func (r *DynamoRuleRepository) Create(ctx context.Context, rule *model.Rule) (*model.Rule, error) {
	if rule.Key == "" {
		rule.Key = ulid.Make().String()
	}

	itemStruct := toRuleItem(rule)
	itemStruct.Pattern = CreateExpression(rule.Path)

	item, err := attributevalue.MarshalMap(itemStruct)
	if err != nil {
		return nil, fmt.Errorf("error marshaling rule: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating rule in DynamoDB: %w", err)
	}

	return rule, nil
}

func (r *DynamoRuleRepository) Update(ctx context.Context, rule *model.Rule) (*model.Rule, error) {
	// Check if rule exists first to stay consistent with other implementations
	_, err := r.Get(ctx, rule.Key)
	if err != nil {
		return nil, err
	}

	return r.Create(ctx, rule)
}

func (r *DynamoRuleRepository) Get(ctx context.Context, key string) (*model.Rule, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"key": &types.AttributeValueMemberS{Value: key},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error getting rule from DynamoDB: %w", err)
	}

	if result.Item == nil {
		return nil, mockserrors.RuleNotFoundError{
			Message: fmt.Sprintf("no rule found with key: %s", key),
		}
	}

	var item ruleItem

	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling rule: %w", err)
	}

	return toRuleModel(&item), nil
}

func (r *DynamoRuleRepository) Search(
	ctx context.Context,
	params map[string]interface{},
	paging model.Paging,
) (*model.RuleList, error) {
	exclusiveStartKey := r.getExclusiveStartKey(paging.LastID)

	filterExpression, attrValues, attrNames := r.buildSearchExpression(params)

	input := &dynamodb.ScanInput{
		TableName:                 aws.String(r.tableName),
		FilterExpression:          filterExpression,
		ExpressionAttributeValues: attrValues,
		ExpressionAttributeNames:  attrNames,
		Limit:                     aws.Int32(paging.Limit),
		ExclusiveStartKey:         exclusiveStartKey,
	}

	result, err := r.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error scanning rules in DynamoDB: %w", err)
	}

	var items []ruleItem

	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling rules: %w", err)
	}

	rules := make([]*model.Rule, 0, len(items))
	for i := range items {
		rules = append(rules, toRuleModel(&items[i]))
	}

	paging.Total = int64(result.ScannedCount)

	return &model.RuleList{
		Results: rules,
		Paging:  paging,
	}, nil
}

func (r *DynamoRuleRepository) buildSearchExpression(
	params map[string]interface{},
) (*string, map[string]types.AttributeValue, map[string]string) {
	if len(params) == 0 {
		return nil, nil, nil
	}

	attrValues := make(map[string]types.AttributeValue)
	attrNames := make(map[string]string)
	filters := make([]string, 0, len(params))
	idx := 0

	for key, val := range params {
		placeholder := fmt.Sprintf(":v%d", idx)
		namePlaceholder := fmt.Sprintf("#n%d", idx)

		filters = append(filters, fmt.Sprintf("contains(to_lower(%s), %s)", namePlaceholder, placeholder))
		attrValues[placeholder] = &types.AttributeValueMemberS{Value: strings.ToLower(fmt.Sprintf("%v", val))}
		attrNames[namePlaceholder] = key
		idx++
	}

	return aws.String(strings.Join(filters, " AND ")), attrValues, attrNames
}

func (r *DynamoRuleRepository) getExclusiveStartKey(lastID string) map[string]types.AttributeValue {
	if lastID == "" {
		return nil
	}

	return map[string]types.AttributeValue{
		"key": &types.AttributeValueMemberS{Value: lastID},
	}
}

func (r *DynamoRuleRepository) SearchByMethodAndPath(
	ctx context.Context,
	method string,
	path string,
) (*model.Rule, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("method-index"),
		KeyConditionExpression: aws.String("#m = :method"),
		ExpressionAttributeNames: map[string]string{
			"#m": "method",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":method": &types.AttributeValueMemberS{Value: strings.ToUpper(method)},
		},
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error querying rules by method: %w", err)
	}

	for _, itemAV := range result.Items {
		var item ruleItem
		if err := attributevalue.UnmarshalMap(itemAV, &item); err != nil {
			continue
		}

		rule, ok, err := r.matchItem(&item, path)
		if err != nil {
			return nil, err
		}

		if ok {
			return rule, nil
		}
	}

	return nil, mockserrors.RuleNotFoundError{
		Message: fmt.Sprintf("no rule found for path: %s and method %s", path, method),
	}
}

func (r *DynamoRuleRepository) matchItem(item *ruleItem, path string) (*model.Rule, bool, error) {
	pattern := item.Pattern
	if pattern == "" {
		pattern = CreateExpression(item.Path)
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, false, fmt.Errorf("invalid regex pattern %s: %w", pattern, err)
	}

	if regex.MatchString(path) {
		if item.Status == model.RuleStatusEnabled {
			return toRuleModel(item), true, nil
		}
	}

	return nil, false, nil
}

func (r *DynamoRuleRepository) Delete(ctx context.Context, key string) error {
	logger := mockscontext.Logger(ctx)

	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"key": &types.AttributeValueMemberS{Value: key},
		},
	})
	if err != nil {
		logger.Error(r, nil, err, "error deleting rule from DynamoDB")

		return fmt.Errorf("error deleting rule: %w", err)
	}

	return nil
}

// Internal Item Structs for DynamoDB mapping (LOWERCASE as per user request).
type ruleItem struct {
	Key       string         `dynamodbav:"key"`
	Group     string         `dynamodbav:"group"`
	Name      string         `dynamodbav:"name"`
	Path      string         `dynamodbav:"path"`
	Strategy  string         `dynamodbav:"strategy"`
	Method    string         `dynamodbav:"method"`
	Status    string         `dynamodbav:"status"`
	Responses []responseItem `dynamodbav:"responses"`
	Variables []variableItem `dynamodbav:"variables"`
	Pattern   string         `dynamodbav:"pattern"`
}

type responseItem struct {
	Body        string `dynamodbav:"body"`
	ContentType string `dynamodbav:"content_type"`
	HTTPStatus  int    `dynamodbav:"http_status"`
	Delay       int    `dynamodbav:"delay"`
	Scene       string `dynamodbav:"scene"`
	Description string `dynamodbav:"description"`
}

type variableItem struct {
	Type       string          `dynamodbav:"type"`
	Name       string          `dynamodbav:"name"`
	Key        string          `dynamodbav:"key"`
	Min        *float64        `dynamodbav:"min,omitempty"`
	Max        *float64        `dynamodbav:"max,omitempty"`
	Decimals   *int            `dynamodbav:"decimals,omitempty"`
	Assertions []assertionItem `dynamodbav:"assertions"`
}

type assertionItem struct {
	FailOnError bool    `dynamodbav:"fail_on_error"`
	Type        string  `dynamodbav:"type"`
	Value       string  `dynamodbav:"value"`
	Min         float64 `dynamodbav:"min"`
	Max         float64 `dynamodbav:"max"`
}

// Mappers.
func toRuleItem(rule *model.Rule) *ruleItem {
	responses := make([]responseItem, 0, len(rule.Responses))
	for _, resp := range rule.Responses {
		responses = append(responses, responseItem{
			Body:        resp.Body,
			ContentType: resp.ContentType,
			HTTPStatus:  resp.HTTPStatus,
			Delay:       resp.Delay,
			Scene:       resp.Scene,
			Description: resp.Description,
		})
	}

	variables := make([]variableItem, 0, len(rule.Variables))
	for _, variable := range rule.Variables {
		assertions := make([]assertionItem, 0, len(variable.Assertions))
		for _, assertion := range variable.Assertions {
			assertions = append(assertions, assertionItem{
				FailOnError: assertion.FailOnError,
				Type:        assertion.Type,
				Value:       assertion.Value,
				Min:         assertion.Min,
				Max:         assertion.Max,
			})
		}

		variables = append(variables, variableItem{
			Type:       variable.Type,
			Name:       variable.Name,
			Key:        variable.Key,
			Min:        variable.Min,
			Max:        variable.Max,
			Decimals:   variable.Decimals,
			Assertions: assertions,
		})
	}

	return &ruleItem{
		Key:       rule.Key,
		Group:     rule.Group,
		Name:      rule.Name,
		Path:      rule.Path,
		Strategy:  rule.Strategy,
		Method:    rule.Method,
		Status:    rule.Status,
		Responses: responses,
		Variables: variables,
	}
}

func toRuleModel(item *ruleItem) *model.Rule {
	responses := make([]model.Response, 0, len(item.Responses))
	for _, resp := range item.Responses {
		responses = append(responses, model.Response{
			Body:        resp.Body,
			ContentType: resp.ContentType,
			HTTPStatus:  resp.HTTPStatus,
			Delay:       resp.Delay,
			Scene:       resp.Scene,
			Description: resp.Description,
		})
	}

	variables := make([]*model.Variable, 0, len(item.Variables))
	for _, variable := range item.Variables {
		assertions := make([]*model.Assertion, 0, len(variable.Assertions))
		for _, assertion := range variable.Assertions {
			assertions = append(assertions, &model.Assertion{
				FailOnError: assertion.FailOnError,
				Type:        assertion.Type,
				Value:       assertion.Value,
				Min:         assertion.Min,
				Max:         assertion.Max,
			})
		}

		variables = append(variables, &model.Variable{
			Type:       variable.Type,
			Name:       variable.Name,
			Key:        variable.Key,
			Min:        variable.Min,
			Max:        variable.Max,
			Decimals:   variable.Decimals,
			Assertions: assertions,
		})
	}

	return &model.Rule{
		Key:       item.Key,
		Group:     item.Group,
		Name:      item.Name,
		Path:      item.Path,
		Strategy:  item.Strategy,
		Method:    item.Method,
		Status:    item.Status,
		Responses: responses,
		Variables: variables,
	}
}
