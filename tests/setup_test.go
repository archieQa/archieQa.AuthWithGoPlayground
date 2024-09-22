package tests

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var (
	testDB *sql.DB
)

func TestMain(m *testing.M) {
	// Setup
	err := setup()
	if err != nil {
		log.Fatalf("Failed to set up test environment: %v", err)
	}

	// Run tests
	code := m.Run()

	// Teardown
	err = teardown()
	if err != nil {
		log.Printf("Failed to tear down test environment: %v", err)
	}

	os.Exit(code)
}

func setup() error {
	// Set up test database connection
	dbHost := os.Getenv("TEST_DB_HOST")
	dbPort := os.Getenv("TEST_DB_PORT")
	dbUser := os.Getenv("TEST_DB_USER")
	dbPassword := os.Getenv("TEST_DB_PASSWORD")
	dbName := os.Getenv("TEST_DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to test database: %v", err)
	}

	// Run migrations or set up schema
	err = runMigrations()
	if err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	// Set up mock data
	err = setupMockData()
	if err != nil {
		return fmt.Errorf("failed to set up mock data: %v", err)
	}

	return nil
}

func teardown() error {
	// Clean up mock data
	err := cleanupMockData()
	if err != nil {
		return fmt.Errorf("failed to clean up mock data: %v", err)
	}

	// Close database connection
	if testDB != nil {
		err := testDB.Close()
		if err != nil {
			return fmt.Errorf("failed to close test database connection: %v", err)
		}
	}

	return nil
}

func runMigrations() error {
	// Implement migration logic here
	// This could involve running SQL scripts or using a migration tool
	return nil
}

func setupMockData() error {
	// Insert mock data into the test database
	// Example:
	_, err := testDB.Exec(`
		INSERT INTO users (username, email, password) 
		VALUES ('testuser', 'test@example.com', 'hashedpassword')
	`)
	if err != nil {
		return err
	}
	return nil
}

func cleanupMockData() error {
	// Clean up all tables in the test database
	// Example:
	_, err := testDB.Exec(`
		TRUNCATE users CASCADE;
	`)
	if err != nil {
		return err
	}
	return nil
}
