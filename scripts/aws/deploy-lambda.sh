#!/bin/bash

# Configuration
REGION="us-east-1"
REPO_NAME="mockserver"
IMAGE_TAG=${1:-"latest"}
FUNCTION_NAME="mockserver-lambda"
ROLE_NAME="mockserver-lambda-role"

# Exit on error
set -e

# Disable AWS CLI pager
export AWS_PAGER=""

echo "🚀 Deploying image tag '$IMAGE_TAG' to AWS Lambda..."

# 1. Get AWS Account ID and ECR URI
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
ECR_URL="${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com"
IMAGE_URI="${ECR_URL}/${REPO_NAME}:${IMAGE_TAG}"
ROLE_ARN="arn:aws:iam::${ACCOUNT_ID}:role/${ROLE_NAME}"

# 2. Parse aws-lambda.env to JSON for Lambda environment
echo "🔑 Reading aws-lambda.env for configuration..."
if [ ! -f "aws-lambda.env" ]; then
    echo "⚠️  No aws-lambda.env file found. Deploying without environment variables."
    ENV_JSON='{"Variables":{}}'
else
    ENV_VARS=$(grep -v '^#' aws-lambda.env | grep -v '^$' | jq -R -n '[inputs | split("=") | {(.[0]): (.[1:] | join("="))}] | add')
    ENV_JSON="{\"Variables\": $ENV_VARS}"
fi

echo "📦 Image URI: $IMAGE_URI"

# 3. Check if the Lambda function exists
FUNCTION_EXISTS=$(aws lambda get-function --function-name "$FUNCTION_NAME" --region "$REGION" > /dev/null 2>&1 && echo "YES" || echo "NO")

if [ "$FUNCTION_EXISTS" == "YES" ]; then
    echo "🔄 Updating existing Lambda function code: $FUNCTION_NAME..."
    aws lambda update-function-code \
        --function-name "$FUNCTION_NAME" \
        --image-uri "$IMAGE_URI" \
        --region "$REGION" \
        --no-cli-pager
    
    echo "⏳ Waiting for update to complete..."
    aws lambda wait function-updated --function-name "$FUNCTION_NAME" --region "$REGION"

    echo "⚙️  Updating Lambda configuration..."
    aws lambda update-function-configuration \
        --function-name "$FUNCTION_NAME" \
        --environment "$ENV_JSON" \
        --region "$REGION" \
        --no-cli-pager
else
    echo "✨ Creating NEW Lambda function: $FUNCTION_NAME..."
    aws lambda create-function \
        --function-name "$FUNCTION_NAME" \
        --package-type Image \
        --code ImageUri="$IMAGE_URI" \
        --role "$ROLE_ARN" \
        --architectures arm64 \
        --environment "$ENV_JSON" \
        --timeout 30 \
        --memory-size 256 \
        --region "$REGION" \
        --no-cli-pager
fi

echo "---------------------------------------------------------"
echo "✅ Lambda Deployment Process Finished!"
echo "---------------------------------------------------------"
