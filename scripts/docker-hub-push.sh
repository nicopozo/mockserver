#!/bin/bash

# Configuration
REPO_NAME="nicopozo/mock-service"
VERSION=$1

# Exit on error
set -e

echo "🚀 Pushing Docker image to Docker Hub..."

if [ -z "$VERSION" ]; then
    echo "🏷️  No version provided, pushing only 'latest'..."
    docker push "$REPO_NAME:latest"
else
    echo "🏷️  Version $VERSION provided, pushing 'latest' and '$VERSION'..."
    docker push "$REPO_NAME:latest"
    docker push "$REPO_NAME:$VERSION"
fi

echo "✅ Push finished successfully!"
