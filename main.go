package main

import (
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"tigerhall-kittens/logger"
)

func main() {
	logger.InitLogger()
	logger.LogInfo("Server Starting...........................................................")
	err := godotenv.Load()
	if err != nil {
		logger.LogError(err)
	}

	logger.LogError(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
	logger.LogInfo("Server listening at port ", os.Getenv("PORT"))
	logger.LogInfo("Server exited and released port ", os.Getenv("PORT"))
}
