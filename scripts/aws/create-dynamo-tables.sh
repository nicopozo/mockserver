#!/bin/bash

# Configuration
REGION="us-east-1"
RULES_TABLE="mockserver_rules"
LOGS_TABLE="mockserver_logs"

# Exit on error
set -e

echo "🚀 Initializing DynamoDB tables in $REGION..."

# Helper function to check if a table exists
table_exists() {
    aws dynamodb describe-table --table-name "$1" --region "$REGION" > /dev/null 2>&1
    return $?
}

# 1. Create Rules Table
if table_exists "$RULES_TABLE"; then
    echo "✅ Table '$RULES_TABLE' already exists."
else
    echo "✨ Creating table '$RULES_TABLE'..."
    aws dynamodb create-table \
        --table-name "$RULES_TABLE" \
        --attribute-definitions \
            AttributeName=key,AttributeType=S \
            AttributeName=method,AttributeType=S \
        --key-schema AttributeName=key,KeyType=HASH \
        --global-secondary-indexes \
            "[
                {
                    \"IndexName\": \"method-index\",
                    \"KeySchema\": [{\"AttributeName\": \"method\",\"KeyType\": \"HASH\"}],
                    \"Projection\": {\"ProjectionType\": \"ALL\"}
                }
            ]" \
        --billing-mode PAY_PER_REQUEST \
        --region "$REGION"
    echo "⏳ Waiting for table '$RULES_TABLE' to be created..."
    aws dynamodb wait table-exists --table-name "$RULES_TABLE" --region "$REGION"
fi

# 2. Create Logs Table
if table_exists "$LOGS_TABLE"; then
    echo "✅ Table '$LOGS_TABLE' already exists."
else
    echo "✨ Creating table '$LOGS_TABLE'..."
    aws dynamodb create-table \
        --table-name "$LOGS_TABLE" \
        --attribute-definitions AttributeName=id,AttributeType=S \
        --key-schema AttributeName=id,KeyType=HASH \
        --billing-mode PAY_PER_REQUEST \
        --region "$REGION"
    echo "⏳ Waiting for table '$LOGS_TABLE' to be created..."
    aws dynamodb wait table-exists --table-name "$LOGS_TABLE" --region "$REGION"
fi

echo "---------------------------------------------------------"
echo "✅ DynamoDB provisioning finished successfully!"
echo "---------------------------------------------------------"
