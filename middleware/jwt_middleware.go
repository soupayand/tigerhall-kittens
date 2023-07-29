package middleware

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"tigerhall-kittens/controller"
	"tigerhall-kittens/logger"
	"tigerhall-kittens/model"
)

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			errRes := controller.ErrorResponse{Error: "Missing JWT token"}
			controller.WriteJSONResponse(w, errRes, http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
			jwtSecret := os.Getenv("JWT_SECRET")
			if jwtSecret == "" {
				logger.LogError(errors.New("JWT_SECRET missing from secret/env file"))
				return nil, errors.New("JWT_SECRET missing from secret/env file")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			errRes := controller.ErrorResponse{Error: "Invalid JWT token"}
			controller.WriteJSONResponse(w, errRes, http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), "user", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			errRes := controller.ErrorResponse{Error: "Invalid JWT token"}
			controller.WriteJSONResponse(w, errRes, http.StatusUnauthorized)
			return
		}
	}
}
