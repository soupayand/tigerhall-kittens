package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"tigerhall-kittens/logger"
	"tigerhall-kittens/model"
)

type IAnimal interface {
	CreateAnimal(animal *model.Animal, sighting *model.Sighting) (*model.Animal, error)
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

func (db *AnimalDB) CreateAnimal(animal *model.Animal, sighting *model.Sighting) (*model.Animal, error) {
	ctx := context.Background()
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to begin transaction")
	}

	if err = createAnimalWithTransaction(ctx, tx, animal); err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			logger.LogError(errors.New("Error rolling back database transaction"))
		}
		return nil, err
	}

	if err = createSightingWithTransaction(ctx, tx, animal, sighting); err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("Failed to commit transaction")
	}
	return animal, nil
}

func createAnimalWithTransaction(ctx context.Context, tx pgx.Tx, animal *model.Animal) error {
	var id int64
	err := tx.QueryRow(ctx,
		`INSERT INTO animal (name, type, variant, date_of_birth, description)
         VALUES($1, $2, $3, TO_DATE($4, 'YYYY-MM-DD'), $5) RETURNING id`,
		animal.Name, animal.Type, animal.Variant, animal.DateOfBirth, animal.Description).Scan(&id)
	if err != nil {
		logger.LogError(err)
		return errors.New("An animal with same name, type or variant already exists")
	}
	animal.ID = id
	return nil
}

func createSightingWithTransaction(ctx context.Context, tx pgx.Tx, animal *model.Animal, sighting *model.Sighting) error {
	var id int64
	geoPoint := fmt.Sprintf("POINT(%f %f)", sighting.LastLocation.Longitude, sighting.LastLocation.Latitude)
	err := tx.QueryRow(ctx,
		`INSERT INTO sighting (animal_id, reporter, last_location, last_seen)
         VALUES($1, $2, ST_GeographyFromText($3), $4) RETURNING id`,
		animal.ID,
		sighting.Reporter.ID,
		geoPoint,
		sighting.LastSeen).Scan(&id)
	if err != nil {
		return err
	}
	sighting.ID = id
	return nil
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
