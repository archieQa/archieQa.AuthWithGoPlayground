#!/bin/bash

set -e

# Define variables
DB_URL=${DB_URL:-"postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"}
MIGRATIONS_DIR="./migrations"

# Function to check if migrate tool is installed
check_migrate() {
    if ! command -v migrate &> /dev/null; then
        echo "Error: 'migrate' tool is not installed. Please install it first."
        echo "You can install it using: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
        exit 1
    fi
}

# Function to run migrations
run_migrations() {
    echo "Applying database migrations..."
    migrate -path ${MIGRATIONS_DIR} -database ${DB_URL} up
    if [ $? -eq 0 ]; then
        echo "Migrations applied successfully."
    else
        echo "Error applying migrations."
        exit 1
    fi
}

# Function to rollback migrations
rollback_migrations() {
    echo "Rolling back the last migration..."
    migrate -path ${MIGRATIONS_DIR} -database ${DB_URL} down 1
    if [ $? -eq 0 ]; then
        echo "Rollback successful."
    else
        echo "Error rolling back migration."
        exit 1
    fi
}

# Function to check migration status
check_migration_status() {
    echo "Checking migration status..."
    migrate -path ${MIGRATIONS_DIR} -database ${DB_URL} version
}

# Main script logic
check_migrate

case "$1" in
    up)
        run_migrations
        ;;
    down)
        rollback_migrations
        ;;
    status)
        check_migration_status
        ;;
    *)
        echo "Usage: $0 {up|down|status}"
        echo "  up     : Apply all available migrations"
        echo "  down   : Rollback the last applied migration"
        echo "  status : Show the current migration version"
        exit 1
        ;;
esac
