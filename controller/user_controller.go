package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tigerhall-kittens/database"
	"tigerhall-kittens/logger"
	"tigerhall-kittens/model"
)

type UserController struct {
	user database.IUser
}

func NewUserController(repo database.IUser) *UserController {
	return &UserController{
		user: repo,
	}
}

// ErrorResponse represents the JSON error response format.
type ErrorResponse struct {
	Error string `json:"error"`
}

func (uc *UserController) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errRes := ErrorResponse{Error: "Method not allowed"}
		writeJSONResponse(w, errRes, http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var user model.User
	err := decoder.Decode(&user)
	if err != nil {
		errRes := ErrorResponse{Error: "Invalid request payload"}
		writeJSONResponse(w, errRes, http.StatusBadRequest)
		return
	}
	userResponse, err := uc.user.CreateUser(&user)
	if err != nil {
		logger.LogError(err)
		errRes := ErrorResponse{Error: fmt.Sprintf("Failed to create user")}
		writeJSONResponse(w, errRes, http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"username": userResponse.Username,
		"userId":   userResponse.ID,
	}
	writeJSONResponse(w, response, http.StatusOK)
}

func writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
