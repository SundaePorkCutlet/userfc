package main

import (
	"userfc/cmd/user/handler"
	"userfc/cmd/user/repository"
	"userfc/cmd/user/resource"
	"userfc/config"
	"userfc/infrastructure/log"
	"userfc/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	log.SetupLogger()

	redis := resource.InitRedis(cfg.Redis)
	db := resource.InitDB(cfg.Database)

	userRepository := repository.NewUserRepository(db, redis)
	userHandler := handler.NewUserHandler()

	port := cfg.App.Port
	router := gin.Default()

	routes.SetupRoutes(router, *userHandler)

	router.Run(":" + port)

	log.Logger.Info().Msgf("Server is running on port %s", port)
}
