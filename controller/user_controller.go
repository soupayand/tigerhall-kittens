package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"tigerhall-kittens/database"
	"tigerhall-kittens/logger"
	"tigerhall-kittens/model"
	"time"
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
		WriteJSONResponse(w, errRes, http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var user model.User
	err := decoder.Decode(&user)
	if err != nil || user.Username == "" || user.Email == "" || user.Password == "" {
		errRes := ErrorResponse{Error: "Invalid request payload"}
		WriteJSONResponse(w, errRes, http.StatusBadRequest)
		return
	}
	userResponse, err := uc.user.CreateUser(&user)
	if err != nil {
		logger.LogError(err)
		errRes := ErrorResponse{Error: fmt.Sprintf("Failed to create user")}
		WriteJSONResponse(w, errRes, http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"username": userResponse.Username,
		"userId":   userResponse.ID,
	}
	WriteJSONResponse(w, response, http.StatusOK)
}

func (uc *UserController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		errRes := ErrorResponse{Error: "Failed to parse login request"}
		WriteJSONResponse(w, errRes, http.StatusBadRequest)
		return
	}
	logger.LogInfo("Username : ", loginReq.Username, " tried to login")
	user, err := authenticateUser(loginReq.Username, loginReq.Password, uc)
	if err != nil {
		errRes := ErrorResponse{Error: "Invalid credentials"}
		WriteJSONResponse(w, errRes, http.StatusUnauthorized)
		return
	}
	token, err := generateJWT(user)
	if err != nil {
		errRes := ErrorResponse{Error: "Failed to generate JWT"}
		WriteJSONResponse(w, errRes, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"token": token,
	}
	json.NewEncoder(w).Encode(response)
	logger.LogInfo(loginReq.Username, " successfully logged in")
}

func authenticateUser(username, password string, uc *UserController) (*model.User, error) {
	user, err := uc.user.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}
	return user, nil
}

func generateJWT(user *model.User) (string, error) {
	claims := model.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logger.LogError(errors.New("JWT_SECRET missing from secret/env file"))
		return "", errors.New("JWT_SECRET missing from secret/env file")
	}
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		logger.LogError(err)
		return "", err
	}

	return signedToken, nil
}

func WriteJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
