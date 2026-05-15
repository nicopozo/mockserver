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

	item, err := attributevalue.MarshalMap(rule)
	if err != nil {
		return nil, fmt.Errorf("error marshaling rule: %w", err)
	}

	// Add pattern attribute (calculated field)
	item["pattern"] = &types.AttributeValueMemberS{Value: CreateExpression(rule.Path)}

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

	var rule model.Rule

	err = attributevalue.UnmarshalMap(result.Item, &rule)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling rule: %w", err)
	}

	return &rule, nil
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

	var rules []*model.Rule

	err = attributevalue.UnmarshalListOfMaps(result.Items, &rules)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling rules: %w", err)
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

	for _, item := range result.Items {
		rule, ok, err := r.matchItem(item, path)
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

func (r *DynamoRuleRepository) matchItem(item map[string]types.AttributeValue, path string) (*model.Rule, bool, error) {
	pattern := ""
	if p, ok := item["pattern"].(*types.AttributeValueMemberS); ok {
		pattern = p.Value
	}

	if pattern == "" {
		var rule model.Rule

		_ = attributevalue.UnmarshalMap(item, &rule)
		pattern = CreateExpression(rule.Path)
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, false, fmt.Errorf("invalid regex pattern %s: %w", pattern, err)
	}

	if regex.MatchString(path) {
		status := ""
		if s, ok := item["status"].(*types.AttributeValueMemberS); ok {
			status = s.Value
		}

		if status == model.RuleStatusEnabled {
			var rule model.Rule

			err = attributevalue.UnmarshalMap(item, &rule)
			if err != nil {
				return nil, false, fmt.Errorf("error unmarshaling matching rule: %w", err)
			}

			return &rule, true, nil
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
