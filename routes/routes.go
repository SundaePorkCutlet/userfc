package routes

import (
	"net/http"
	"userfc/cmd/user/handler"
	"userfc/cmd/user/resource"
	"userfc/config"
	"userfc/infrastructure/ratelimiter"
	"userfc/infrastructure/tokenblacklist"
	"userfc/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine, userHandler *handler.UserHandler, rl *ratelimiter.RateLimiter, bl *tokenblacklist.TokenBlacklist) {

	router.Use(middleware.RequestLogger())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/ping", userHandler.Ping())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "userfc",
		})
	})

	router.GET("/debug/queries", func(c *gin.Context) {
		if resource.DBMonitor == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "monitor not initialized"})
			return
		}
		c.JSON(http.StatusOK, resource.DBMonitor.GetDebugInfo())
	})

	router.GET("/debug/redis", func(c *gin.Context) {
		ctx := c.Request.Context()
		info := map[string]interface{}{
			"status": "connected",
		}
		dbSize, err := resource.RedisClient.DBSize(ctx).Result()
		if err == nil {
			info["db_size"] = dbSize
		}
		c.JSON(http.StatusOK, info)
	})

	rateLimited := router.Group("/")
	rateLimited.Use(ratelimiter.RateLimitMiddleware(rl))
	{
		rateLimited.POST("/v1/register", userHandler.Register)
		rateLimited.POST("/v1/login", userHandler.Login)
	}

	authMiddleware := middleware.AuthMiddlewareWithBlacklist(config.GetJwtSecret(), bl)

	router.POST("/v1/logout", authMiddleware, userHandler.Logout)

	private := router.Group("/api")
	private.Use(authMiddleware)
	{
		private.GET("/v1/user-info", userHandler.GetUserInfo)
	}
}
