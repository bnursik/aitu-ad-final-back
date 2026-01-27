package server

import (
	"net/http"

	"github.com/bnursik/aitu-ad-final-back/internal/app"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine, c *app.Container) {
	r.GET("/api/v1/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/api/v1/auth/register", c.Auth.Register)
	r.POST("/api/v1/auth/login", c.Auth.Login)
}
