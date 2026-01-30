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
	admin.GET("/orders/:id", c.Orders.Get)
	admin.POST("/orders/find", c.Orders.FindOrderByID)

	// admin statistics
	admin.POST("/statistics/sales/date-range", c.Statistics.GetSalesStatsByDateRange)
	admin.POST("/statistics/sales/year", c.Statistics.GetSalesStatsByYear)
	admin.GET("/statistics/sales", c.Statistics.GetSalesStatsAll)
	admin.POST("/statistics/products/date-range", c.Statistics.GetProductsStatsByDateRange)
	admin.POST("/statistics/products/year", c.Statistics.GetProductsStatsByYear)
	admin.GET("/statistics/products", c.Statistics.GetProductsStatsAll)

	v1.POST("/auth/register", c.Auth.Register)
	v1.POST("/auth/login", c.Auth.Login)
	admin.POST("/auth/register", c.Auth.AdminRegister)

	// profile: auth required
	v1.GET("/profile", middleware.AuthRequired(c.JWT), c.Auth.GetProfile)
	v1.PUT("/profile", middleware.AuthRequired(c.JWT), c.Auth.UpdateProfile)

	// wishlist: auth required
	wishlistGroup := v1.Group("/wishlist")
	wishlistGroup.Use(middleware.AuthRequired(c.JWT))
	wishlistGroup.POST("", c.Wishlist.Add)
	wishlistGroup.GET("", c.Wishlist.List)
	wishlistGroup.DELETE("/:id", c.Wishlist.Delete)

}
