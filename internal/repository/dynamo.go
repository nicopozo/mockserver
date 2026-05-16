package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/nicopozo/mockserver/internal/configs"
)

// NewDynamoClient initializes and returns a new DynamoDB client.
func NewDynamoClient(ctx context.Context, cfg *configs.Config) (*dynamodb.Client, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.AWS.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	var client *dynamodb.Client

	// For local development with LocalStack or DynamoDB Local
	if cfg.Dynamo.Endpoint != "" {
		client = dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(cfg.Dynamo.Endpoint)
		})
	} else {
		client = dynamodb.NewFromConfig(awsCfg)
	}

	return client, nil
}
