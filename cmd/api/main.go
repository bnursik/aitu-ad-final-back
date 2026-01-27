package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/bnursik/aitu-ad-final-back/internal/config"
	"github.com/bnursik/aitu-ad-final-back/internal/db"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := db.Connect(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	defer func() {
		_ = mongoClient.Disconnect(context.Background())
	}()

	_ = mongoClient.Database(cfg.DBName) // на будущее

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// health check
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Printf("listening on :%s (db=%s)", cfg.Port, cfg.DBName)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server run: %v", err)
	}
}
