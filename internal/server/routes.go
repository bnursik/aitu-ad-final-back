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

	v1.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// public
	v1.GET("/categories", c.Categories.List)
	v1.GET("/categories/:id", c.Categories.Get)

	// public products
	v1.GET("/products", c.Products.List)
	v1.GET("/products/:id", c.Products.Get)

	v1.POST("/products/:id/reviews", middleware.AuthRequired(c.JWT), c.Products.AddReview)
	v1.DELETE("/products/:id/reviews/:reviewId", middleware.AuthRequired(c.JWT), c.Products.DeleteReview)

	// orders: auth required (user + admin)
	ordersGroup := v1.Group("/orders")
	ordersGroup.Use(middleware.AuthRequired(c.JWT))
	ordersGroup.POST("", c.Orders.Create)
	ordersGroup.GET("", c.Orders.List)
	ordersGroup.GET("/:id", c.Orders.Get)

	// admin products
	admin := v1.Group("/admin")
	admin.Use(middleware.AuthRequired(c.JWT), middleware.AdminOnly())

	admin.POST("/products", c.Products.Create)
	admin.PUT("/products/:id", c.Products.Update)
	admin.DELETE("/products/:id", c.Products.Delete)

	admin.POST("/categories", c.Categories.Create)
	admin.PUT("/categories/:id", c.Categories.Update)
	admin.DELETE("/categories/:id", c.Categories.Delete)

	admin.PUT("/orders/:id/status", c.Orders.UpdateStatus)

	v1.POST("/auth/register", c.Auth.Register)
	v1.POST("/auth/login", c.Auth.Login)
	admin.POST("/auth/admin/register", c.Auth.AdminRegister)

}
