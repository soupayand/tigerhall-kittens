// userdb_test.go
package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"testing"
	"tigerhall-kittens/model"
)

func TestCreateUser(t *testing.T) {
	connConfig, err := pgxpool.ParseConfig("postgresql://root:secret@localhost:5432/tigerhall_kittens?sslmode=disable")
	connConfig.MaxConns = 10
	pool, err := pgxpool.ConnectConfig(context.Background(), connConfig)
	defer pool.Close()
	userDB := NewUserDB(pool)
	testUser := &model.User{
		Username: "testuser",
		Password: "testpassword",
		Email:    "testuser@example.com",
	}
	var createdUser *model.User
	createdUser, err = userDB.CreateUser(testUser)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	defer func() {
		if createdUser != nil && (*createdUser).ID != 0 {
			_, err := pool.Exec(context.Background(), `DELETE FROM "user" WHERE username = 'testuser'`)
			if err != nil {
				t.Errorf("Failed to delete test user: %v", err)
			}
		}
	}()
	assert.NotEmpty(t, createdUser.ID, "User ID should not be empty")
	assert.Equal(t, testUser.Username, createdUser.Username, "Username mismatch")
}

func TestHashPassword(t *testing.T) {
	// Test data
	password := "testpassword"

	// Test hashPassword
	hashedPassword, err := hashPassword(password)

	// Assertions
	assert.NoError(t, err, "Unexpected error while hashing password")
	assert.NotEmpty(t, hashedPassword, "Hashed password should not be empty")
}

func TestGetUserByUsername(t *testing.T) {
	connConfig, err := pgxpool.ParseConfig("postgresql://root:secret@localhost:5432/tigerhall_kittens?sslmode=disable")
	connConfig.MaxConns = 10
	pool, err := pgxpool.ConnectConfig(context.Background(), connConfig)
	defer pool.Close()
	userDB := NewUserDB(pool)
	testUser := &model.User{
		Username: "testuser",
		Password: "testpassword",
		Email:    "testuser@example.com",
	}
	var createdUser *model.User
	createdUser, err = userDB.CreateUser(testUser)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	defer func() {
		if createdUser != nil && (*createdUser).ID != 0 {
			_, err := pool.Exec(context.Background(), `DELETE FROM "user" WHERE username = 'testuser'`)
			if err != nil {
				t.Errorf("Failed to delete test user: %v", err)
			}
		}
	}()
	db := NewUserDB(pool)
	username := "testuser"
	expectedUser := &model.User{
		Username: "testuser",
		Email:    "testuser@example.com",
	}
	sqlQuery := `SELECT id, username, password, email FROM "user" where username = $1`
	rows, err := db.pool.Query(context.Background(), sqlQuery, username)
	if err != nil {
		t.Fatalf("Unable to fetch user details from username")
	}
	defer rows.Close()

	var data model.User
	for rows.Next() {
		err = rows.Scan(
			&data.ID,
			&data.Username,
			&data.Password,
			&data.Email,
		)
	}
	// Call the GetUserByUsername method
	user, err := db.GetUserByUsername(username)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if user.Username != expectedUser.Username || user.Email != expectedUser.Email {
		t.Errorf("Returned user does not match the expected user")
	}
}
