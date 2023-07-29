package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"tigerhall-kittens/logger"
	"tigerhall-kittens/model"
)

type IUser interface {
	CreateUser(user *model.User) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
}

type UserDB struct {
	pool *pgxpool.Pool
}

func NewUserDB(pool *pgxpool.Pool) *UserDB {
	return &UserDB{
		pool: pool,
	}
}

func (db *UserDB) CreateUser(user *model.User) (*model.User, error) {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	var id int
	err = db.pool.QueryRow(context.Background(),
		`INSERT INTO "user" (username, password, email) VALUES($1, $2, $3) RETURNING id`,
		user.Username, hashedPassword, user.Email).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	logger.LogInfo("User with username", user.Username, " created")
	return &model.User{
		ID:       int64(id),
		Username: user.Username,
	}, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	return string(hashedPassword), nil
}

func (db *UserDB) GetUserByUsername(username string) (*model.User, error) {
	sqlQuery := `SELECT id, username, password, email FROM "user" where username = $1`
	rows, err := db.pool.Query(context.Background(), sqlQuery, username)
	if err != nil {
		logger.LogError(err)
		return nil, err
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
	return &data, nil
}
