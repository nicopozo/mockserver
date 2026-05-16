#!/bin/bash

# Configuration
REGION="us-east-1"
REPO_NAME="mockserver"
IMAGE_TAG=${1:-"latest"}

# Exit on error
set -e

echo "🚀 Pushing Docker image to AWS ECR..."

# 1. Get AWS Account ID
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
ECR_URL="${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com"

# 2. Login to ECR
echo "🔐 Logging in to ECR..."
aws ecr get-login-password --region "$REGION" | docker login --username AWS --password-stdin "$ECR_URL"

# 3. Create repository if it doesn't exist
REPO_EXISTS=$(aws ecr describe-repositories --repository-names "$REPO_NAME" --region "$REGION" > /dev/null 2>&1 && echo "YES" || echo "NO")
if [ "$REPO_EXISTS" == "NO" ]; then
    echo "✨ Creating ECR repository: $REPO_NAME..."
    aws ecr create-repository --repository-name "$REPO_NAME" --region "$REGION"
fi

# 4. Tag and Push
echo "📦 Tagging and pushing image..."
docker tag nicopozo/mock-service:latest "${ECR_URL}/${REPO_NAME}:${IMAGE_TAG}"
docker push "${ECR_URL}/${REPO_NAME}:${IMAGE_TAG}"

if [ "$IMAGE_TAG" != "latest" ]; then
    docker tag nicopozo/mock-service:latest "${ECR_URL}/${REPO_NAME}:latest"
    docker push "${ECR_URL}/${REPO_NAME}:latest"
fi

echo "✅ Push finished successfully!"
