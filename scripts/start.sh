#!/bin/bash

# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=myuser
export DB_PASSWORD=mypassword
export DB_NAME=myapp
export JWT_SECRET=mysecretkey
export PORT=8080

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "Docker is not running. Please start Docker and try again."
    exit 1
fi

# Start the PostgreSQL container if it's not already running
if ! docker ps | grep -q postgres; then
    echo "Starting PostgreSQL container..."
    docker run --name postgres -e POSTGRES_USER=$DB_USER -e POSTGRES_PASSWORD=$DB_PASSWORD -e POSTGRES_DB=$DB_NAME -p $DB_PORT:5432 -d postgres
    
    # Wait for PostgreSQL to be ready
    echo "Waiting for PostgreSQL to be ready..."
    sleep 10
fi

# Run database migrations
echo "Running database migrations..."
go run cmd/migrate/main.go

# Build and run the application
echo "Building and starting the application..."
go build -o myapp cmd/main.go
./myapp
