package main

import (
	"os"
	"testing"
)

func TestSetupDatabase(t *testing.T) {
	// Set up the required environment variables for testing
	// Note: You may need to adjust these values based on your specific test environment
	os.Setenv("DATABASE_URL", "postgresql://root:secret@localhost:5432/tigerhall_kittens?sslmode=disable")
	os.Setenv("MAX_CONNECTIONS", "10")

	// Call the function being tested
	pool, err := setupDatabase()

	// Clean up the environment variables after the test
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("MAX_CONNECTIONS")

	// Check for any errors returned by the function
	if err != nil {
		t.Fatalf("setupDatabase returned an error: %v", err)
	}

	// Check if the pool is not nil
	if pool == nil {
		t.Fatal("setupDatabase returned a nil pool")
	}

	// Perform additional tests on the pool if required
	// For example, you can check if the MaxConns setting is as expected
	expectedMaxConnections := int32(10)
	if pool.Config().MaxConns != expectedMaxConnections {
		t.Fatalf("setupDatabase: expected MaxConns to be %d, got %d", expectedMaxConnections, pool.Config().MaxConns)
	}
}
