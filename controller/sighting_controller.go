package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		decoder := json.NewDecoder(r.Body)
		var sighting model.SightingReqResp
		err := decoder.Decode(&sighting)
		if err != nil {
			errRes := ErrorResponse{Error: "Invalid request payload"}
			WriteJSONResponse(w, errRes, http.StatusBadRequest)
			return
		}
		userID, _ := r.Context().Value("user_id").(int64)
		sighting.Reporter.ID = userID
		createdAnimal, err := sc.sighting.CreateSighting(&sighting)
		if err != nil {
			logger.LogError(err)
			errRes := ErrorResponse{Error: fmt.Sprintf("Failed to create animal : %v", err)}
			WriteJSONResponse(w, errRes, http.StatusInternalServerError)
			return
		}
		WriteJSONResponse(w, createdAnimal, http.StatusCreated)
	default:
		errRes := ErrorResponse{Error: "Method not allowed"}
		WriteJSONResponse(w, errRes, http.StatusMethodNotAllowed)
	}
}
