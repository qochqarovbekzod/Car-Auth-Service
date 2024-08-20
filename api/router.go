package api

import (
	"auth/api/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "auth/api/docs"
)

// @title Auth Service API
// @version 1.0
// @description This is a sample server for Auth Service.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /auth
// @schemes http
func NewRouter(handle *handler.Handler) *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(handler.CORSMiddleware())

	auth := router.Group("/auth")
	{
		auth.POST("/register", handle.RegisterHandler)
		auth.POST("/login", handle.LoginHandler)
		auth.POST("/refreshtoken", handle.RefreshToken)
		auth.POST("/logout", handle.LogoutHandler)
	}
	return router
}
