package main

import (
	"email-dispatcher/api"
	"email-dispatcher/config"
	"email-dispatcher/database"
	"email-dispatcher/mailer"
	"email-dispatcher/worker"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	database.ConnectDB()
	database.ConnectRedis()

	// Create router
	r := gin.Default()

	// Register routes
	api.RegisterRoutes(r)

	// Initialize mailer (mock)
	mockMailer := &mailer.MockMailer{}

	// Create worker
	w := worker.Worker{
		Mailer: mockMailer,
		Queue:  api.GetQueue(),
	}

	// Start worker

	go w.Start()

	// Start server
	r.Run(":" + config.GetEnv("PORT"))
}
