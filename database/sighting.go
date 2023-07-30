package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"tigerhall-kittens/model"
)

type ISighting interface {
	CreateSighting(sighting *model.SightingReqResp) (*model.SightingReqResp, error)
	ListSightingInfo(name string, animalType string, limit string, offset string) ([]model.SightingReqResp, error)
}

type SightingDB struct {
	pool *pgxpool.Pool
}

func NewSightingDB(pool *pgxpool.Pool) *SightingDB {
	return &SightingDB{
		pool: pool,
	}
}

func (db *SightingDB) CreateSighting(sightingReq *model.SightingReqResp) (*model.SightingReqResp, error) {
	ctx := context.Background()
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to begin transaction")
	}
	if err = CreateSightingWithTransaction(ctx, tx, sightingReq.AnimalID, &sightingReq.Sighting); err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("Failed to commit transaction")
	}
	return sightingReq, nil
}

func (db *SightingDB) ListSightingInfo(name string, animalType string, limit string, offset string) ([]model.SightingReqResp, error) {
	//TODO implement me
	panic("implement me")
}
