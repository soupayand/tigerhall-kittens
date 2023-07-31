package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"math"
	"strings"
	"tigerhall-kittens/logger"
	"tigerhall-kittens/model"
)

const earthRadiusKm = 6371.0
const thresholdDistanceInKms = 5

type ISighting interface {
	CreateSighting(sighting *model.SightingReqResp) (*model.SightingReqResp, error)
	ListSightingInfo(animalId int64, limit string, offset string) ([]model.SightingReqResp, error)
	SpottedByEmailIds(animalId int64) (string, error)
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
	lastLocation, err := db.GetLastLocation(sightingReq.AnimalID)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	distance := getDistance(lastLocation, &sightingReq.Location)
	if distance <= thresholdDistanceInKms {
		err = errors.New(fmt.Sprintf("The supplied location is within the threshold distance of %v kms", thresholdDistanceInKms))
		logger.LogError(err)
		return nil, err
	}
	if err = createImageWithTransaction(ctx, tx, &sightingReq.Sighting.Image); err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			return nil, err
		}
		return nil, err
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

func createImageWithTransaction(ctx context.Context, tx pgx.Tx, image *model.Image) error {
	var id int64
	err := tx.QueryRow(ctx,
		`INSERT INTO image (filename, type, data)
         VALUES($1, $2, $3) RETURNING id`,
		image.FileName,
		image.Type,
		image.Data).Scan(&id)
	if err != nil {
		return err
	}
	image.ID = id
	return nil
}

func (db *SightingDB) GetLastLocation(animalId int64) (*model.Point, error) {
	params := make([]interface{}, 0)
	sqlQuery := `
				SELECT s.location[0] AS longitude, s.location[1] AS latitude
				FROM sighting s
				WHERE animal_id = $1
				AND spotting_timestamp =
				    (
    				SELECT MAX(spotting_timestamp)
    				FROM sighting
    				WHERE animal_id = s.animal_id
					);
				`
	params = append(params, animalId)
	rows, err := db.pool.Query(context.Background(), sqlQuery, params...)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	defer rows.Close()
	var response model.Point
	for rows.Next() {
		err = rows.Scan(
			&response.Longitude,
			&response.Latitude,
		)
		if err != nil {
			logger.LogError(err)
			return nil, err
		}
	}
	logger.LogInfo("Retrieved last location for animal with id", animalId)
	return &response, nil
}

func getDistance(lastLocation *model.Point, currentLocation *model.Point) float64 {
	lon1Rad := toRadians(lastLocation.Longitude)
	lat1Rad := toRadians(lastLocation.Latitude)
	lon2Rad := toRadians(currentLocation.Longitude)
	lat2Rad := toRadians(currentLocation.Latitude)

	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad
	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Pow(math.Sin(deltaLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadiusKm * c
	return distance
}

func toRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func (db *SightingDB) ListSightingInfo(animalId int64, limit string, offset string) ([]model.SightingReqResp, error) {
	sqlQuery := `
		SELECT a.id, s.location[0] as longitude, s.location[1] as latitude, TO_CHAR(s.spotting_timestamp, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), i.filename, i.type, i.data
		FROM animal a
		JOIN sighting s ON a.id = s.animal_id
		LEFT OUTER JOIN image i ON s.image_id = i.id
	`
	params := make([]interface{}, 0)
	sqlQuery += " WHERE a.id = $1"
	params = append(params, animalId)
	sqlQuery += " ORDER BY s.spotting_timestamp DESC"
	sqlQuery += " LIMIT $2 OFFSET $3"
	params = append(params, limit, offset)
	rows, err := db.pool.Query(context.Background(), sqlQuery, params...)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	defer rows.Close()
	var responseArray []model.SightingReqResp
	for rows.Next() {
		var response model.SightingReqResp
		var nullFileName, nullType sql.NullString
		err = rows.Scan(
			&response.AnimalID,
			&response.Sighting.Location.Longitude,
			&response.Sighting.Location.Latitude,
			&response.Sighting.SpottingTimestamp,
			&nullFileName,
			&nullType,
			&response.Sighting.Image.Data,
		)
		if err != nil {
			logger.LogError(err)
			return nil, err
		}
		if nullFileName.Valid {
			response.Sighting.Image.FileName = nullFileName.String
		} else {
			response.Sighting.Image.FileName = ""
		}
		if nullType.Valid {
			response.Sighting.Image.Type = nullType.String
		} else {
			response.Sighting.Image.Type = ""
		}
		responseArray = append(responseArray, response)
	}
	logger.LogInfo("Retrieved animal list info")
	return responseArray, nil
}

func (db *SightingDB) SpottedByEmailIds(animalId int64) (string, error) {
	sqlQuery := `
		SELECT distinct u.email
		FROM animal a
		JOIN sighting s ON a.id = s.animal_id
		JOIN "user" u ON s.reporter = u.id
	`
	params := make([]interface{}, 0)
	sqlQuery += " WHERE a.id = $1"
	params = append(params, animalId)
	rows, err := db.pool.Query(context.Background(), sqlQuery, params...)
	if err != nil {
		logger.LogError(err)
		return "", err
	}
	defer rows.Close()
	var emailList []string
	for rows.Next() {
		var email string
		err = rows.Scan(&email)
		if err != nil {
			logger.LogError(err)
			return "", err
		}
		emailList = append(emailList, email)
	}
	logger.LogInfo("Retrieved email addresses of users who spotted the animal")
	return strings.Join(emailList, ","), nil
}
