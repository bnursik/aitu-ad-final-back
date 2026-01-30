package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/bnursik/aitu-ad-final-back/docs"
	"github.com/bnursik/aitu-ad-final-back/internal/app"
	"github.com/bnursik/aitu-ad-final-back/internal/config"
	"github.com/bnursik/aitu-ad-final-back/internal/server"
)

// @title Peripherals Store API
// @version 1.0
// @description REST API for Computer Peripherals Store (MongoDB + Gin)
// @host localhost:8080
// @host aitu-ad-final-back-production.up.railway.app
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	container, err := app.Build(cfg)
	if err != nil {
		log.Fatalf("app build: %v", err)
	}

	router := server.NewRouter(container)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = srv.Shutdown(ctx)
	if container.Shutdown != nil {
		_ = container.Shutdown(ctx)
	}
}
