#!/bin/bash

set -e

# Define variables
APP_NAME="myapp"
MAIN_PATH="cmd/main.go"
OUTPUT_DIR="./build"
BINARY_NAME="${APP_NAME}"

# Ensure the output directory exists
mkdir -p ${OUTPUT_DIR}

echo "Checking and downloading dependencies..."
go mod tidy
go mod verify

echo "Running tests..."
go test ./...

echo "Checking for compilation errors..."
go vet ./...

echo "Building the application..."
go build -o ${OUTPUT_DIR}/${BINARY_NAME} ${MAIN_PATH}

if [ $? -eq 0 ]; then
    echo "Build successful! Binary located at ${OUTPUT_DIR}/${BINARY_NAME}"
else
    echo "Build failed."
    exit 1
fi

echo "Running linter..."
if command -v golangci-lint &> /dev/null; then
    golangci-lint run
else
    echo "golangci-lint not found. Skipping linting."
fi

echo "Build process completed."
