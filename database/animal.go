package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"strconv"
	"tigerhall-kittens/logger"
	"tigerhall-kittens/model"
)

type IAnimal interface {
	CreateAnimal(animal *model.Animal, sighting *model.Sighting) (*model.AnimalSighting, error)
	ListAnimalInfo(name string, animalType string, limit string, offset string) ([]model.AnimalSighting, error)
}

type AnimalDB struct {
	pool *pgxpool.Pool
}

func NewAnimalDB(pool *pgxpool.Pool) *AnimalDB {
	return &AnimalDB{
		pool: pool,
	}
}

func (db *AnimalDB) CreateAnimal(animal *model.Animal, sighting *model.Sighting) (*model.AnimalSighting, error) {
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
	return &model.AnimalSighting{
		*animal,
		*sighting,
	}, nil
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
	err := tx.QueryRow(ctx,
		`INSERT INTO sighting (animal_id, reporter, last_location, last_seen)
         VALUES($1, $2, point($3, $4), $5) RETURNING id`,
		animal.ID,
		sighting.Reporter.ID,
		sighting.LastLocation.Longitude,
		sighting.LastLocation.Latitude,
		sighting.LastSeen).Scan(&id)
	if err != nil {
		return err
	}
	sighting.ID = id
	return nil
}

func (db *AnimalDB) ListAnimalInfo(name string, animalType string, limit string, offset string) ([]model.AnimalSighting, error) {
	sqlQuery := `
		SELECT a.name, a.type, a.variant, TO_CHAR(a.date_of_birth, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), a.description, last_location[0] AS longitude, last_location[1] AS latitude,TO_CHAR(s.last_seen, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		FROM animal a
		JOIN sighting s ON a.id = s.animal_id
		WHERE 1=1
	`
	params := make([]interface{}, 0)
	if name != "" {
		sqlQuery += " AND a.name = $" + strconv.Itoa(len(params)+1)
		params = append(params, name)
	}
	if name == "" {
		sqlQuery += " AND a.type = $" + strconv.Itoa(len(params)+1)
		params = append(params, "tiger")
	} else {
		sqlQuery += " AND a.type = $" + strconv.Itoa(len(params)+2)
		params = append(params, animalType)
	}
	sqlQuery += " ORDER BY s.last_seen DESC"
	sqlQuery += " LIMIT $" + strconv.Itoa(len(params)+1) + " OFFSET $" + strconv.Itoa(len(params)+2)
	params = append(params, limit, offset)
	rows, err := db.pool.Query(context.Background(), sqlQuery, params...)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	defer rows.Close()
	var responseArray []model.AnimalSighting
	for rows.Next() {
		var response model.AnimalSighting
		err = rows.Scan(
			&response.Animal.Name,
			&response.Animal.Type,
			&response.Animal.Variant,
			&response.Animal.DateOfBirth,
			&response.Animal.Description,
			&response.Sighting.LastLocation.Longitude,
			&response.Sighting.LastLocation.Latitude,
			&response.Sighting.LastSeen,
		)
		if err != nil {
			logger.LogError(err)
			return nil, err
		}
		responseArray = append(responseArray, response)
	}
	logger.LogInfo("Retrieved animal list info")
	return responseArray, nil
}
