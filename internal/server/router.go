package server

import (
	"time"

	"github.com/bnursik/aitu-ad-final-back/internal/app"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(c *app.Container) *gin.Engine {
	r := gin.New()
	corsConfig := cors.Config{
		AllowOrigins: []string{
			"https://mangustad.vercel.app",
			"http://localhost:5173",
		},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Authorization", "Content-Type",
		},
		ExposeHeaders: []string{
			"Set-Cookie",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(corsConfig))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	RegisterRoutes(r, c)
	return r
}
