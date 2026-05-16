package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/nicopozo/mockserver/internal/configs"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/oklog/ulid/v2"
)

type DynamoLogRepository struct {
	client    *dynamodb.Client
	tableName string
}

// NewDynamoLogRepository creates a new LogRepository for DynamoDB.
func NewDynamoLogRepository(client *dynamodb.Client, cfg *configs.Config) LogRepository {
	return &DynamoLogRepository{
		client:    client,
		tableName: cfg.Dynamo.TablePrefix + "logs",
	}
}

func (r *DynamoLogRepository) Add(ctx context.Context, entry model.LogEntry) error {
	if entry.ID == "" {
		entry.ID = ulid.Make().String()
	}

	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	item := logItem{
		ID:              entry.ID,
		Timestamp:       entry.Timestamp,
		Method:          entry.Method,
		URL:             entry.URL,
		RequestBody:     entry.RequestBody,
		RequestHeaders:  entry.RequestHeaders,
		QueryParams:     entry.QueryParams,
		ResponseStatus:  entry.ResponseStatus,
		ResponseBody:    entry.ResponseBody,
		AssertionErrors: entry.AssertionErrors,
	}

	attributes, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("error marshaling log item: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      attributes,
	})
	if err != nil {
		return fmt.Errorf("error adding log to DynamoDB: %w", err)
	}

	return nil
}

func (r *DynamoLogRepository) GetAll(ctx context.Context, paging model.Paging) (model.LogList, error) {
	var exclusiveStartKey map[string]types.AttributeValue

	if paging.LastID != "" {
		exclusiveStartKey = map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: paging.LastID},
		}
	}

	input := &dynamodb.ScanInput{
		TableName:         aws.String(r.tableName),
		Limit:             aws.Int32(paging.Limit),
		ExclusiveStartKey: exclusiveStartKey,
	}

	result, err := r.client.Scan(ctx, input)
	if err != nil {
		return model.LogList{}, fmt.Errorf("error scanning logs in DynamoDB: %w", err)
	}

	var items []logItem

	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return model.LogList{}, fmt.Errorf("error unmarshaling logs: %w", err)
	}

	results := make([]model.LogEntry, 0, len(items))
	for _, item := range items {
		results = append(results, model.LogEntry{
			ID:              item.ID,
			Timestamp:       item.Timestamp,
			Method:          item.Method,
			URL:             item.URL,
			RequestBody:     item.RequestBody,
			RequestHeaders:  item.RequestHeaders,
			QueryParams:     item.QueryParams,
			ResponseStatus:  item.ResponseStatus,
			ResponseBody:    item.ResponseBody,
			AssertionErrors: item.AssertionErrors,
		})
	}

	paging.Total = int64(result.ScannedCount)

	return model.LogList{
		Results: results,
		Paging:  paging,
	}, nil
}

func (r *DynamoLogRepository) Clear(ctx context.Context) error {
	input := &dynamodb.ScanInput{
		TableName:            aws.String(r.tableName),
		ProjectionExpression: aws.String("id"),
	}

	paginator := dynamodb.NewScanPaginator(r.client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error scanning for clear: %w", err)
		}

		for _, item := range page.Items {
			id := item["id"]

			_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
				TableName: aws.String(r.tableName),
				Key: map[string]types.AttributeValue{
					"id": id,
				},
			})
			if err != nil {
				return fmt.Errorf("error deleting item during clear: %w", err)
			}
		}
	}

	return nil
}

type logItem struct {
	ID              string            `dynamodbav:"id"`
	Timestamp       time.Time         `dynamodbav:"timestamp"`
	Method          string            `dynamodbav:"method"`
	URL             string            `dynamodbav:"url"`
	RequestBody     string            `dynamodbav:"request_body"`
	RequestHeaders  map[string]string `dynamodbav:"request_headers"`
	QueryParams     map[string]string `dynamodbav:"query_params"`
	ResponseStatus  int               `dynamodbav:"response_status"`
	ResponseBody    string            `dynamodbav:"response_body"`
	AssertionErrors []string          `dynamodbav:"assertion_errors"`
}
