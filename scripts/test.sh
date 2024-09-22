#!/bin/bash

set -e

# Define variables
PROJECT_ROOT=$(pwd)
TEST_OUTPUT_DIR="${PROJECT_ROOT}/test_results"
COVERAGE_FILE="${TEST_OUTPUT_DIR}/coverage.out"
REPORT_FILE="${TEST_OUTPUT_DIR}/test_report.txt"

# Ensure the test output directory exists
mkdir -p ${TEST_OUTPUT_DIR}

# Set up test environment variables
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5433
export TEST_DB_USER=testuser
export TEST_DB_PASSWORD=testpassword
export TEST_DB_NAME=testdb
export TEST_JWT_SECRET=testsecretkey

# Start test database
echo "Starting test PostgreSQL container..."
docker run --name postgres_test -e POSTGRES_USER=${TEST_DB_USER} -e POSTGRES_PASSWORD=${TEST_DB_PASSWORD} -e POSTGRES_DB=${TEST_DB_NAME} -p ${TEST_DB_PORT}:5432 -d postgres

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
sleep 10

# Run database migrations for test database
echo "Running database migrations for test database..."
TEST_DB_URL="postgres://${TEST_DB_USER}:${TEST_DB_PASSWORD}@localhost:${TEST_DB_PORT}/${TEST_DB_NAME}?sslmode=disable"
go run cmd/migrate/main.go -db ${TEST_DB_URL}

echo "Running tests and generating coverage report..."
go test -v -coverprofile=${COVERAGE_FILE} ./... 2>&1 | tee ${REPORT_FILE}

echo "Generating HTML coverage report..."
go tool cover -html=${COVERAGE_FILE} -o ${TEST_OUTPUT_DIR}/coverage.html

# Check if any tests failed
if grep -q "FAIL" ${REPORT_FILE}; then
    echo "Some tests failed. Please check the test report for details."
    TEST_STATUS="FAILED"
else
    echo "All tests passed successfully!"
    TEST_STATUS="PASSED"
fi

# Clean up
echo "Stopping and removing test PostgreSQL container..."
docker stop postgres_test
docker rm postgres_test

echo "Test execution completed. Status: ${TEST_STATUS}"
echo "Test report: ${REPORT_FILE}"
echo "Coverage report: ${TEST_OUTPUT_DIR}/coverage.html"

# Exit with appropriate status
if [ "${TEST_STATUS}" = "FAILED" ]; then
    exit 1
else
    exit 0
fi
