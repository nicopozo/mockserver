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

	item, err := attributevalue.MarshalMap(entry)
	if err != nil {
		return fmt.Errorf("error marshaling log entry: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
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

	var logs []model.LogEntry

	err = attributevalue.UnmarshalListOfMaps(result.Items, &logs)
	if err != nil {
		return model.LogList{}, fmt.Errorf("error unmarshaling logs: %w", err)
	}

	paging.Total = int64(result.ScannedCount)

	return model.LogList{
		Results: logs,
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
