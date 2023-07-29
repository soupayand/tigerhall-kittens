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

func (ac *AnimalController) CreateAnimalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errRes := ErrorResponse{Error: "Method not allowed"}
		WriteJSONResponse(w, errRes, http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var animal model.Animal
	err := decoder.Decode(&animal)
	if err != nil {
		errRes := ErrorResponse{Error: "Invalid request payload"}
		WriteJSONResponse(w, errRes, http.StatusBadRequest)
		return
	}
	createdAnimal, err := ac.animal.CreateAnimal(&animal)
	if err != nil {
		logger.LogError(err)
		errRes := ErrorResponse{Error: fmt.Sprintf("Failed to create user")}
		WriteJSONResponse(w, errRes, http.StatusInternalServerError)
		return
	}
	WriteJSONResponse(w, createdAnimal, http.StatusOK)
}
