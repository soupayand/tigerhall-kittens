package main

import "tigerhall-kittens/worker"
import "tigerhall-kittens/logger"

func main() {
	logger.LogInfo("Starting worker..............................")
	worker.StartConsumer()
}
