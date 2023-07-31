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
	CreateAnimal(animal *model.Animal, sighting *model.Sighting) (*model.AnimalReqResp, error)
	ListAnimalInfo(name string, animalType string, limit string, offset string) ([]model.AnimalReqResp, error)
}

type AnimalDB struct {
	pool *pgxpool.Pool
}

func NewAnimalDB(pool *pgxpool.Pool) *AnimalDB {
	return &AnimalDB{
		pool: pool,
	}
}

func (db *AnimalDB) CreateAnimal(animal *model.Animal, sighting *model.Sighting) (*model.AnimalReqResp, error) {
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

	if err = CreateSightingWithTransaction(ctx, tx, animal.ID, sighting); err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("Failed to commit transaction")
	}
	return &model.AnimalReqResp{
		AnimalID: animal.ID,
		Animal:   *animal,
		Sighting: *sighting,
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

func CreateSightingWithTransaction(ctx context.Context, tx pgx.Tx, animalId int64, sighting *model.Sighting) error {
	var id int64
	err := tx.QueryRow(ctx,
		`INSERT INTO sighting (animal_id, reporter, location, spotting_timestamp)
         VALUES($1, $2, point($3, $4), $5) RETURNING id`,
		animalId,
		sighting.Reporter.ID,
		sighting.Location.Longitude,
		sighting.Location.Latitude,
		sighting.SpottingTimestamp).Scan(&id)
	if err != nil {
		return err
	}
	sighting.ID = id
	return nil
}

func (db *AnimalDB) ListAnimalInfo(name string, animalType string, limit string, offset string) ([]model.AnimalReqResp, error) {
	sqlQuery := `
		SELECT a.id, a.name, a.type, a.variant, TO_CHAR(a.date_of_birth, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), a.description, longitude, latitude, TO_CHAR(s.spotting_timestamp, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		FROM animal a
		JOIN (
			SELECT animal_id, location[0] as longitude,location[1] as latitude,MAX(spotting_timestamp) as spotting_timestamp
			FROM sighting
			GROUP BY animal_id, longitude, latitude
		) s ON a.id = s.animal_id
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
	sqlQuery += " ORDER BY s.spotting_timestamp DESC"
	sqlQuery += " LIMIT $" + strconv.Itoa(len(params)+1) + " OFFSET $" + strconv.Itoa(len(params)+2)
	params = append(params, limit, offset)
	rows, err := db.pool.Query(context.Background(), sqlQuery, params...)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	defer rows.Close()
	var responseArray []model.AnimalReqResp
	for rows.Next() {
		var response model.AnimalReqResp
		err = rows.Scan(
			&response.AnimalID,
			&response.Animal.Name,
			&response.Animal.Type,
			&response.Animal.Variant,
			&response.Animal.DateOfBirth,
			&response.Animal.Description,
			&response.Sighting.Location.Longitude,
			&response.Sighting.Location.Latitude,
			&response.Sighting.SpottingTimestamp,
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
