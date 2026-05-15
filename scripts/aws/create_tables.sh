#!/bin/bash

# Configuration
TABLE_PREFIX=${DYNAMO_TABLE_PREFIX:-"mockserver_"}
REGION=${AWS_REGION:-"us-east-1"}
ENDPOINT_URL=${DYNAMO_ENDPOINT:-""}

AWS_ARGS="--region $REGION"
if [ -n "$ENDPOINT_URL" ]; then
    AWS_ARGS="$AWS_ARGS --endpoint-url $ENDPOINT_URL"
fi

echo "Creating DynamoDB tables with prefix: $TABLE_PREFIX"

# Create Rules Table
echo "Creating table: ${TABLE_PREFIX}rules..."
aws dynamodb create-table $AWS_ARGS \
    --table-name "${TABLE_PREFIX}rules" \
    --attribute-definitions \
        AttributeName=key,AttributeType=S \
        AttributeName=method,AttributeType=S \
    --key-schema \
        AttributeName=key,KeyType=HASH \
    --global-secondary-indexes \
        "[
            {
                \"IndexName\": \"method-index\",
                \"KeySchema\": [{\"AttributeName\":\"method\",\"KeyType\":\"HASH\"}],
                \"Projection\": {\"ProjectionType\":\"ALL\"}
            }
        ]" \
    --billing-mode PAY_PER_REQUEST

# Create Logs Table
echo "Creating table: ${TABLE_PREFIX}logs..."
aws dynamodb create-table $AWS_ARGS \
    --table-name "${TABLE_PREFIX}logs" \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
    --key-schema \
        AttributeName=id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST

echo "Tables creation requested."
