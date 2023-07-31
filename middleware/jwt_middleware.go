package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"tigerhall-kittens/controller"
	"tigerhall-kittens/logger"
)

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				errRes := controller.ErrorResponse{Error: "Missing JWT token"}
				controller.WriteJSONResponse(w, errRes, http.StatusUnauthorized)
				return
			}
			tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
			token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				jwtSecret := os.Getenv("JWT_SECRET")
				if jwtSecret == "" {
					logger.LogError(errors.New("JWT_SECRET missing from secret/env file"))
					return nil, errors.New("JWT_SECRET missing from secret/env file")
				}
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				errRes := controller.ErrorResponse{Error: "Invalid JWT token"}
				controller.WriteJSONResponse(w, errRes, http.StatusUnauthorized)
				return
			}
			if claims, ok := token.Claims.(*CustomClaims); ok {
				ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				errRes := controller.ErrorResponse{Error: "Invalid JWT claims"}
				controller.WriteJSONResponse(w, errRes, http.StatusUnauthorized)
				return
			}
			break
		case http.MethodGet:
			next.ServeHTTP(w, r)

		}
	}
}
