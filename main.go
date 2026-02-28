// Package main USERFC - User & Auth service.
package main

import (
	"context"
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
	"userfc/tracing"

	_ "userfc/docs"

	"github.com/gin-gonic/gin"
)

// @title           USERFC API
// @version         1.0
// @description     User registration, authentication (JWT), and user info for Go Commerce.
// @host            localhost:28080
// @BasePath        /
// @schemes         http
func main() {
	cfg := config.LoadConfig()

	log.SetupLogger()

	// Vault에서 secrets 로드
	config.LoadVaultSecrets(&cfg)

	// Tracing 초기화
	shutdownTracer, err := tracing.InitTracer(cfg.Tracing)
	if err != nil {
		log.Logger.Warn().Err(err).Msg("Failed to initialize tracing - continuing without tracing")
	} else {
		defer shutdownTracer(context.Background())
	}

	redis := resource.InitRedis(cfg.Redis)
	db := resource.InitDB(cfg.Database)

	// AutoMigrate: 데이터베이스 테이블 자동 생성/업데이트
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Logger.Fatal().Err(err).Msg("Failed to migrate database")
	}
	log.Logger.Info().Msg("Database migration completed")

	userRepository := repository.NewUserRepository(db, redis)
	userService := service.NewUserService(userRepository)
	userUsecase := usecase.NewUserUsecase(*userService)
	userHandler := handler.NewUserHandler(*userUsecase)

	// gRPC 서버 시작 (별도 고루틴)
	go usergrpc.StartGRPCServer(cfg.App.GRPCPort, userService)

	port := cfg.App.Port
	router := gin.Default()

	// 트레이싱 미들웨어 추가
	if cfg.Tracing.Enabled {
		router.Use(tracing.GinMiddleware(cfg.Tracing.ServiceName))
	}

	routes.SetupRoutes(router, userHandler)

	log.Logger.Info().Msgf("HTTP server is running on port %s", port)
	router.Run(":" + port)
}
