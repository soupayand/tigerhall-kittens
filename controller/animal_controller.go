package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tigerhall-kittens/database"
	"tigerhall-kittens/logger"
	"tigerhall-kittens/model"
)

type AnimalController struct {
	animal database.IAnimal
}

func NewAnimalController(a database.IAnimal) *AnimalController {
	return &AnimalController{
		animal: a,
	}
}

func (ac *AnimalController) AnimalHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
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
		animals, err := ac.animal.ListAnimalInfo(name, animalType, limit, offset)
		if err != nil {
			errRes := ErrorResponse{Error: "Failed to retrieve animal information"}
			WriteJSONResponse(w, errRes, http.StatusInternalServerError)
			return
		}
		WriteJSONResponse(w, animals, http.StatusOK)
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var animalReq model.AnimalReqResp
		err := decoder.Decode(&animalReq)
		if err != nil {
			errRes := ErrorResponse{Error: "Invalid request payload"}
			WriteJSONResponse(w, errRes, http.StatusBadRequest)
			return
		}
		userID, _ := r.Context().Value("user_id").(int64)
		animalReq.Reporter.ID = userID
		createdAnimal, err := ac.animal.CreateAnimal(&animalReq.Animal, &animalReq.Sighting)
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
