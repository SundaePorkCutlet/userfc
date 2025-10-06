package main

import (
	"userfc/cmd/user/handler"
	"userfc/config"
	"userfc/infrastructure/log"
	"userfc/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	log.SetupLogger()

	userHandler := handler.NewUserHandler()

	port := cfg.App.Port
	router := gin.Default()

	routes.SetupRoutes(router, *userHandler)

	router.Run(":" + port)

	log.Logger.Info().Msgf("Server is running on port %s", port)
}
