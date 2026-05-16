#!/bin/bash

# Configuration
REGION="us-east-1"
API_NAME="mockserver-api"
FUNCTION_NAME="mockserver-lambda"

# Exit on error
set -e

echo "🚀 Configuring API Gateway (HTTP API) for $FUNCTION_NAME..."

# 1. Get AWS Account ID
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
FUNCTION_ARN="arn:aws:lambda:${REGION}:${ACCOUNT_ID}:function:${FUNCTION_NAME}"

# 2. Check if API already exists
API_ID=$(aws apigatewayv2 get-apis --query "Items[?Name=='$API_NAME'].ApiId" --output text)

if [ -z "$API_ID" ]; then
    echo "✨ Creating NEW HTTP API: $API_NAME..."
    API_ID=$(aws apigatewayv2 create-api \
        --name "$API_NAME" \
        --protocol-type HTTP \
        --query ApiId --output text)
else
    echo "✅ API '$API_NAME' already exists with ID: $API_ID"
fi

# 3. Create Lambda Integration
INTEGRATION_ID=$(aws apigatewayv2 get-integrations --api-id "$API_ID" --query "Items[?IntegrationUri=='$FUNCTION_ARN'].IntegrationId" --output text)

if [ -z "$INTEGRATION_ID" ]; then
    echo "🔗 Creating Lambda integration..."
    INTEGRATION_ID=$(aws apigatewayv2 create-integration \
        --api-id "$API_ID" \
        --integration-type AWS_PROXY \
        --integration-uri "$FUNCTION_ARN" \
        --payload-format-version 2.0 \
        --query IntegrationId --output text)
else
    echo "✅ Lambda integration already exists."
fi

# 4. Create default route ($default)
ROUTE_EXISTS=$(aws apigatewayv2 get-routes --api-id "$API_ID" --query "Items[?RouteKey=='\$default']" --output text)

if [ -z "$ROUTE_EXISTS" ]; then
    echo "🛣️  Creating default route..."
    aws apigatewayv2 create-route \
        --api-id "$API_ID" \
        --route-key '$default' \
        --target "integrations/$INTEGRATION_ID"
else
    echo "✅ Default route already exists."
fi

# 5. Create default stage
STAGE_EXISTS=$(aws apigatewayv2 get-stages --api-id "$API_ID" --query "Items[?StageName=='\$default']" --output text)

if [ -z "$STAGE_EXISTS" ]; then
    echo "🏁 Creating default stage..."
    aws apigatewayv2 create-stage \
        --api-id "$API_ID" \
        --stage-name '$default' \
        --auto-deploy
else
    echo "✅ Default stage already exists."
fi

# 6. Add Lambda Permission for API Gateway
echo "🔐 Adding Lambda permission for API Gateway..."
aws lambda add-permission \
    --function-name "$FUNCTION_NAME" \
    --statement-id apigateway-access \
    --action lambda:InvokeFunction \
    --principal apigateway.amazonaws.com \
    --source-arn "arn:aws:execute-api:${REGION}:${ACCOUNT_ID}:${API_ID}/*/*" \
    --region "$REGION" || echo "⚠️  Permission might already exist, continuing..."

echo "---------------------------------------------------------"
echo "✅ API Gateway configuration finished!"
echo "Invoke URL: https://${API_ID}.execute-api.${REGION}.amazonaws.com/"
echo "---------------------------------------------------------"
