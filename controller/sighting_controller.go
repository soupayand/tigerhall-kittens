package controller

import (
	"bytes"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"tigerhall-kittens/database"
	"tigerhall-kittens/logger"
	"tigerhall-kittens/model"
)

type SightingController struct {
	sighting database.ISighting
}

func NewSightingController(s database.ISighting) *SightingController {
	return &SightingController{
		sighting: s,
	}
}

func (sc *SightingController) SightingHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		/*
			queryParams := r.URL.Query()
			name := queryParams.Get("name")
			animalType := queryParams.Get("type")
			limit := queryParams.Get("limit")
			offset := queryParams.Get("offset")
			if limit == "" || offset == "" {
				errRes := ErrorResponse{Error: "limit or offset query parameter(s) missing"}
				WriteJSONResponse(w, errRes, http.StatusBadRequest)
				return
			}
			WriteJSONResponse(w, nil, http.StatusOK)
		*/
	case http.MethodPost:
		err := r.ParseMultipartForm(10) // Max size of 10 MB for the image
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
		animalID := r.FormValue("animal_id")
		latitude := r.FormValue("latitude")
		longitude := r.FormValue("longitude")
		spottingTimestamp := r.FormValue("spotting_timestamp")
		animalIDInt, err := strconv.ParseInt(animalID, 10, 64)
		if err != nil {
			http.Error(w, "Invalid value for animal_id", http.StatusBadRequest)
			return
		}
		latitudeFloat, err := strconv.ParseFloat(latitude, 64)
		if err != nil {
			http.Error(w, "Invalid value for location.latitude", http.StatusBadRequest)
			return
		}
		longitudeFloat, err := strconv.ParseFloat(longitude, 64)
		if err != nil {
			http.Error(w, "Invalid value for location.longitude", http.StatusBadRequest)
			return
		}
		sightingReq := model.SightingReqResp{
			AnimalID: animalIDInt,
			Sighting: model.Sighting{
				Location: model.Point{
					Latitude:  latitudeFloat,
					Longitude: longitudeFloat,
				},
				SpottingTimestamp: spottingTimestamp,
			},
		}
		imageFile, imageHeader, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Error retrieving image from form", http.StatusBadRequest)
			return
		}
		defer func(imageFile multipart.File) {
			err := imageFile.Close()
			if err != nil {
				logger.LogError(err)
			}
		}(imageFile)
		var i model.Image
		i.FileName = imageHeader.Filename
		i.Type = imageHeader.Header.Get("Content-Type")
		i.Data, err = resizeImage(imageFile)
		if err != nil {
			http.Error(w, "Failed to resize i", http.StatusInternalServerError)
			return
		}
		userID, _ := r.Context().Value("user_id").(int64)
		sightingReq.Reporter.ID = userID
		sightingReq.Image = i
		createdAnimal, err := sc.sighting.CreateSighting(&sightingReq)
		if err != nil {
			logger.LogError(err)
			errRes := ErrorResponse{Error: fmt.Sprintf("Failed to create sighting : %v", err)}
			WriteJSONResponse(w, errRes, http.StatusInternalServerError)
			return
		}
		WriteJSONResponse(w, createdAnimal, http.StatusCreated)
	default:
		errRes := ErrorResponse{Error: "Method not allowed"}
		WriteJSONResponse(w, errRes, http.StatusMethodNotAllowed)
	}

}

func resizeImage(imageFile io.Reader) ([]byte, error) {
	img, _, err := image.Decode(imageFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}
	resizedImage := resize.Resize(250, 200, img, resize.Lanczos3)
	resizedImageBytes := new(bytes.Buffer)
	if err := jpeg.Encode(resizedImageBytes, resizedImage, nil); err != nil {
		return nil, fmt.Errorf("failed to encode resized image: %v", err)
	}

	return resizedImageBytes.Bytes(), nil
}
