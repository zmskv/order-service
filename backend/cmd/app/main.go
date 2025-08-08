package main

import "github.com/zmskv/order-service/logger"

func main() {
	log := logger.New()
	defer log.Sync()

	log.Info("Application started")
	log.Debug("This is a debug message")
}
