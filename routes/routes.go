package routes

import (
	"userfc/cmd/user/handler"
	"userfc/config"
	"userfc/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userHandler *handler.UserHandler) {

	router.Use(middleware.RequestLogger())
	//public API
	router.GET("/ping", userHandler.Ping())
	router.POST("/v1/register", userHandler.Register)
	router.POST("/v1/login", userHandler.Login)
	//private API
	private := router.Group("/api")
	private.Use(middleware.AuthMiddleware(config.GetJwtSecret()))
	{
		private.GET("/v1/user-info", userHandler.GetUserInfo)
	}
}
