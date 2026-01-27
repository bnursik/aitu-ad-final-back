package server

import (
	"github.com/bnursik/aitu-ad-final-back/internal/app"
	"github.com/gin-gonic/gin"
)

func NewRouter(c *app.Container) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	RegisterRoutes(r, c)
	return r
}
