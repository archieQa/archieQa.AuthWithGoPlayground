# Microservice Documentation

This document provides instructions on how to set up, run, and use the microservice.

## Table of Contents

1. [Installation](#installation)
2. [Environment Setup](#environment-setup)
3. [Database Configuration](#database-configuration)
4. [Running the Application](#running-the-application)
5. [API Usage Examples](#api-usage-examples)

## Installation

1. Clone the repository:

   ```
   git clone https://github.com/your-username/your-project.git
   cd your-project
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

## Environment Setup

1. Create a `.env` file in the root directory of the project.
2. Add the following environment variables:
   ```
   APP_NAME=myapp
   ENV=development
   PORT=8080
   DB_HOST=localhost
   DB_PORT=5432
   DB_NAME=your_database_name
   DB_USER=your_database_user
   DB_PASSWORD=your_database_password
   JWT_SECRET=your_jwt_secret_key
   ```

## Database Configuration

1. Ensure you have PostgreSQL installed and running.
2. Create a new database:
   ```
   createdb your_database_name
   ```
3. Run the database migrations:
   ```
   go run cmd/migrate/main.go
   ```

## Running the Application

1. Build the application:

   ```
   ./scripts/build.sh
   ```

2. Run the application:
   ```
   ./build/myapp
   ```

The application should now be running on `http://localhost:8080`.

## API Usage Examples

### User Registration
