package middleware_test

import (
	"github.com/joho/godotenv"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"tigerhall-kittens/middleware"
)

func TestJWTMiddleware_ValidToken(t *testing.T) {
	err := os.Chdir("../")
	if err != nil {
		t.Fatalf("Error changing working directory: %s", err)
	}
	err = godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %s", err)
	}
	//Place a new valid token string below
	tokenString := ""
	if tokenString == "" {
		t.Fatalf("In order to test this please assign a new valid token in the above tokenString variable")
	}
	// Create a test request with a valid JWT token in the Authorization header
	req := httptest.NewRequest("POST", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	// Create a test response recorder
	rr := httptest.NewRecorder()

	// Create a mock handler for the next middleware in the chain
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure that the user_id is set in the request context
		userID, ok := r.Context().Value("user_id").(int64)
		if !ok {
			t.Error("user_id not set in the request context")
			return
		}

		// Verify the user_id value based on your test data
		expectedUserID := int64(1) // Change this to the expected user_id for your test data
		if userID != expectedUserID {
			t.Errorf("expected user_id %d, but got %d", expectedUserID, userID)
		}

		// Write a response indicating the successful execution of the handler
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the mock handler with the JWTMiddleware
	handler := middleware.JWTMiddleware(mockHandler)

	// Call the handler with the test request and response recorder
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, rr.Code)
	}
}

func TestJWTMiddleware_MissingToken(t *testing.T) {
	// Create a test request without a JWT token in the Authorization header
	req := httptest.NewRequest("POST", "/test", nil)

	// Create a test response recorder
	rr := httptest.NewRecorder()

	// Create a mock handler for the next middleware in the chain
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("mockHandler should not be called due to missing token")
	})

	// Wrap the mock handler with the JWTMiddleware
	handler := middleware.JWTMiddleware(mockHandler)

	// Call the handler with the test request and response recorder
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status code %d, but got %d", http.StatusUnauthorized, rr.Code)
	}
}

// Add more test cases to cover other scenarios, such as invalid tokens and invalid claims.
// Remember to mock the necessary dependencies, such as jwt.ParseWithClaims and controller.WriteJSONResponse.
// The above tests demonstrate basic scenarios for the JWTMiddleware function.
