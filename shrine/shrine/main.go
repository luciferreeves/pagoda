package main

import (
	"fmt"
	"os"
	"os/signal"
	"shrine/config"
	"shrine/database"
	"shrine/jobs"
	"shrine/middleware"
	"shrine/router"
	"shrine/utils/logger"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          router.ErrorHandler,
		BodyLimit:             int(config.Storage.MaxFileSize) * 2,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:  config.Server.CorsOrigins,
		AllowMethods:  "GET, HEAD, PUT, PATCH, POST, DELETE, OPTIONS",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization, X-Requested-With, X-API-Key, X-CSRF-Token",
		ExposeHeaders: "Content-Length, Content-Type, Content-Disposition, X-Pagination, X-Total-Count",
		MaxAge:        86400,
	}))
	app.Use(helmet.New())

	middleware.Initialize(app)
	router.Initialize(app)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)); err != nil {
			logger.Fatalf("Main", "Failed to start the server on %s:%d: %v", config.Server.Host, config.Server.Port, err)
		}
	}()

	jobs.Start()
	logger.Successf("Main", "Server started on %s:%d", config.Server.Host, config.Server.Port)

	<-quit
	logger.Infof("Main", "Shutting down gracefully...")
	jobs.Stop()

	if err := app.Shutdown(); err != nil {
		logger.Errorf("Main", "Error during server shutdown: %v", err)
	}

	if err := database.Close(); err != nil {
		logger.Errorf("Main", "Error closing database connection: %v", err)
	}

	logger.Successf("Main", "Shutdown complete")
}
