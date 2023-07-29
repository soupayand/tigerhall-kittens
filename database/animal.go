package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"tigerhall-kittens/logger"
	"tigerhall-kittens/model"
)

type IAnimal interface {
	CreateAnimal(animal *model.Animal) (*model.Animal, error)
	ListAnimalInfo(name string, animalType string, limit string, offset string) ([]model.Animal, error)
}

type AnimalDB struct {
	pool *pgxpool.Pool
}

func NewAnimalDB(pool *pgxpool.Pool) *AnimalDB {
	return &AnimalDB{
		pool: pool,
	}
}

func (db *AnimalDB) CreateAnimal(animal *model.Animal) (*model.Animal, error) {
	var id int
	err := db.pool.QueryRow(context.Background(),
		`INSERT INTO animal (name, type, variant, date_of_birth, description) VALUES($1, $2, $3, TO_DATE($4, 'YYYY-MM-DD'), $5) RETURNING id`,
		animal.Name, animal.Type, animal.Variant, animal.DateOfBirth, animal.Description).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	animal.ID = int64(id)
	return animal, nil
}

func (db *AnimalDB) ListAnimalInfo(name string, animalType string, limit string, offset string) ([]model.Animal, error) {
	sqlQuery := "SELECT name, type, variant, date_of_birth, description FROM animal WHERE 1=1"
	params := []interface{}{}

	if name != "" {
		sqlQuery += " AND name = $1"
		params = append(params, name)
	}
	if animalType != "" {
		sqlQuery += " AND type = $2"
		params = append(params, animalType)
	}
	sqlQuery += " LIMIT $3 OFFSET $4"
	params = append(params, limit, offset)

	rows, err := db.pool.Query(context.Background(), sqlQuery, params...)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	defer rows.Close()

	var animals []model.Animal
	for rows.Next() {
		var animal model.Animal
		err = rows.Scan(
			&animal.Name,
			&animal.Type,
			&animal.Variant,
			&animal.DateOfBirth,
			&animal.Description,
		)
		if err != nil {
			logger.LogError(err)
			return nil, err
		}
		animals = append(animals, animal)
	}
	logger.LogInfo("Retrieved animal list info")
	return animals, nil
}
