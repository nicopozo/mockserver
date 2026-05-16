#!/bin/bash

# Configuration
REPO_NAME="nicopozo/mock-service"
VERSION=$1

# Dynamically detect platform to avoid mismatch warnings on Mac M1/M2/M3
PLATFORM_FLAG=""
if [ -z "$DOCKER_DEFAULT_PLATFORM" ]; then
    DETECTED_ARCH=$(uname -m)
    # Convert macOS nomenclature to Docker nomenclature
    case "$DETECTED_ARCH" in
        x86_64) ARCH="amd64" ;;
        arm64)  ARCH="arm64" ;;
        *)      ARCH="amd64" ;; # Fallback
    esac
    
    # Docker images are ALWAYS linux based
    PLATFORM_FLAG="--platform linux/$ARCH"
fi

# Exit on error
set -e

echo "🔨 Building Docker image using $PLATFORM_FLAG..."

# Prepare build command
TAG_LATEST="$REPO_NAME:latest"
if [ -z "$VERSION" ]; then
    echo "🏷️  No version provided, building only 'latest'..."
    docker build $PLATFORM_FLAG --provenance=false -t "$TAG_LATEST" .
else
    echo "🏷️  Version $VERSION provided, building 'latest' and '$VERSION'..."
    docker build $PLATFORM_FLAG --provenance=false -t "$TAG_LATEST" -t "$REPO_NAME:$VERSION" .
fi

echo "✅ Image built successfully."
