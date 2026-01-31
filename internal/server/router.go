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
			"https://aitu-ad-final-back-production.up.railway.app",
			"https://mangustad.vercel.app",
			"http://localhost:5173",
			"http://localhost:8080",
			"http://localhost:3000",
		},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders: []string{
			"Authorization", "Content-Type", "Origin", "Accept", "X-Requested-With",
		},
		ExposeHeaders: []string{
			"Set-Cookie", "Content-Length",
		},
		AllowCredentials: true,
		AllowAllOrigins:  false,
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(corsConfig))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Static("/static", "./static")
	RegisterRoutes(r, c)
	return r
}
