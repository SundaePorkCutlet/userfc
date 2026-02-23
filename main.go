package main

import (
	"userfc/cmd/user/handler"
	"userfc/cmd/user/repository"
	"userfc/cmd/user/resource"
	"userfc/cmd/user/service"
	"userfc/cmd/user/usecase"
	"userfc/config"
	usergrpc "userfc/grpc"
	"userfc/infrastructure/log"
	"userfc/models"
	"userfc/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	log.SetupLogger()

	redis := resource.InitRedis(cfg.Redis)
	db := resource.InitDB(cfg.Database)

	// AutoMigrate: 데이터베이스 테이블 자동 생성/업데이트
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Logger.Fatal().Err(err).Msg("Failed to migrate database")
	}
	log.Logger.Info().Msg("Database migration completed")

	userRepository := repository.NewUserRepository(db, redis)
	userService := service.NewUserService(*userRepository)
	userUsecase := usecase.NewUserUsecase(*userService)
	userHandler := handler.NewUserHandler(*userUsecase)

	// gRPC 서버 시작 (별도 고루틴)
	go usergrpc.StartGRPCServer(cfg.App.GRPCPort, userService)

	port := cfg.App.Port
	router := gin.Default()

	routes.SetupRoutes(router, userHandler)

	log.Logger.Info().Msgf("HTTP server is running on port %s", port)
	router.Run(":" + port)
}
