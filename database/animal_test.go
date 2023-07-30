package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"testing"
	"tigerhall-kittens/model"
	"time"
)

func TestCreateAnimal(t *testing.T) {
	connConfig, err := pgxpool.ParseConfig("postgresql://root:secret@localhost:5432/tigerhall_kittens?sslmode=disable")
	connConfig.MaxConns = 10
	pool, err := pgxpool.ConnectConfig(context.Background(), connConfig)
	defer pool.Close()
	animalDB := NewAnimalDB(pool)
	animal := &model.Animal{
		Name:        "testanimal",
		Type:        "testtype",
		Variant:     "testvariant",
		DateOfBirth: "2018-05-15",
		Description: "testdescription",
	}
	lastSeenStr := "2023-07-30T12:34:56Z"
	lastSeenTime, err := time.Parse(time.RFC3339, lastSeenStr)
	if err != nil {
		fmt.Printf("Error while parsing the date-time string: %v", err)
		return
	}
	sighting := &model.Sighting{
		LastLocation: model.Point{
			Latitude:  51.503399,
			Longitude: -0.119519,
		},
		LastSeen: lastSeenTime,
		Reporter: model.User{
			ID: 0,
		},
	}
	var createdAnimal *model.AnimalSighting

	createdAnimal, err = animalDB.CreateAnimal(animal, sighting)
	if err != nil {
		t.Fatalf("Failed to create animal: %v", err)
	}
	defer func() {
		if createdAnimal != nil && (*createdAnimal).Animal.ID != 0 {
			_, err := pool.Exec(context.Background(), `DELETE FROM sighting WHERE id = $1`, sighting.ID)
			if err != nil {
				t.Errorf("Failed to delete test user: %v", err)
			}
			_, err = pool.Exec(context.Background(), `DELETE FROM animal WHERE id = $1`, animal.ID)
			if err != nil {
				t.Errorf("Failed to delete test user: %v", err)
			}
		}
	}()
	assert.NotEmpty(t, createdAnimal.Animal.ID, "Animal ID should not be empty")
	assert.Equal(t, animal.Name, createdAnimal.Name, "Animal name mismatch")
}
