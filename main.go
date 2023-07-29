package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"strconv"
	"tigerhall-kittens/controller"
	"tigerhall-kittens/database"
	"tigerhall-kittens/logger"
	"tigerhall-kittens/middleware"
)

func main() {
	logger.InitLogger()
	logger.LogInfo("Server Starting...........................................................")
	err := godotenv.Load()
	if err != nil {
		logger.LogError(err)
	}
	pool, err := setupDatabase()
	if err != nil {
		logger.LogError(err)
		return
	}
	defer pool.Close()

	//DAO
	user := database.NewUserDB(pool)
	animal := database.NewAnimalDB(pool)

	//Controllers
	userController := controller.NewUserController(user)
	animalController := controller.NewAnimalController(animal)

	//Middlewares
	jwlMiddleWare := middleware.JWTMiddleware

	//Register handlers/controllers
	http.HandleFunc("/user", userController.CreateUserHandler)
	http.HandleFunc("/user/login", userController.LoginHandler)
	http.HandleFunc("/animal", jwlMiddleWare(animalController.CreateAnimalHandler))

	logger.LogError(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
	logger.LogInfo("Server listening at port ", os.Getenv("PORT"))
	logger.LogInfo("Server exited and released port ", os.Getenv("PORT"))
}

func setupDatabase() (*pgxpool.Pool, error) {
	connConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("error parsing DATABASE_URL: %w", err)
	}
	maxConnections, err := strconv.Atoi(os.Getenv("MAX_CONNECTIONS"))
	if err != nil {
		return nil, fmt.Errorf("error converting MAX_CONNECTIONS to int: %w", err)
	}
	connConfig.MaxConns = int32(maxConnections)
	pool, err := pgxpool.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}
	logger.LogInfo("Database connection pool is set up successfully!")
	return pool, nil
}
