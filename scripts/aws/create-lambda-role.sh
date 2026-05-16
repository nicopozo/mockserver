#!/bin/bash

# Configuration
ROLE_NAME="mockserver-lambda-role"

# Exit on error
set -e

echo "🔑 1. Creating Trust Policy for Lambda..."
cat > /tmp/mockserver-lambda-trust-policy.json << 'EOF'
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF

echo "✨ 2. Creating IAM Role: $ROLE_NAME..."
aws iam create-role \
    --role-name "$ROLE_NAME" \
    --assume-role-policy-document file:///tmp/mockserver-lambda-trust-policy.json > /dev/null

echo "🔗 3. Attaching Basic Execution Policy (for CloudWatch Logs)..."
aws iam attach-role-policy \
    --role-name "$ROLE_NAME" \
    --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole

echo "🗄️  4. Attaching DynamoDB Policy..."
aws iam attach-role-policy \
    --role-name "$ROLE_NAME" \
    --policy-arn arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess

echo "🔐  5. Attaching Secrets Manager Policy..."
aws iam attach-role-policy \
    --role-name "$ROLE_NAME" \
    --policy-arn arn:aws:iam::aws:policy/SecretsManagerReadWrite

echo "⏳ Waiting 10 seconds for AWS to propagate role permissions..."
sleep 10

echo "✅ Success! The role '$ROLE_NAME' has been created."
echo "You can now run: make aws-lambda-deploy"
