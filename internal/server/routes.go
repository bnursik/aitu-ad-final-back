package server

import (
	"net/http"

	"github.com/bnursik/aitu-ad-final-back/internal/app"
	"github.com/bnursik/aitu-ad-final-back/internal/http/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine, c *app.Container) {
	v1 := r.Group("/api/v1")

	// public
	v1.GET("/categories", c.Categories.List)
	v1.GET("/categories/:id", c.Categories.Get)

	admin := v1.Group("/admin")
	admin.Use(middleware.AuthRequired(c.JWT), middleware.AdminOnly())
	admin.POST("/categories", c.Categories.Create)
	admin.PUT("/categories/:id", c.Categories.Update)
	admin.DELETE("/categories/:id", c.Categories.Delete)

	v1.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1.POST("/auth/register", c.Auth.Register)
	v1.POST("/auth/login", c.Auth.Login)
}
