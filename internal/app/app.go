package app

import (
	"context"
	"time"

	"github.com/bnursik/aitu-ad-final-back/internal/config"
	"github.com/bnursik/aitu-ad-final-back/internal/db"
	"github.com/bnursik/aitu-ad-final-back/internal/http/handlers"
	"github.com/bnursik/aitu-ad-final-back/internal/http/middleware"
	mongorepo "github.com/bnursik/aitu-ad-final-back/internal/repository/mongo"
	categoriessvc "github.com/bnursik/aitu-ad-final-back/internal/services/categories"
	orderssvc "github.com/bnursik/aitu-ad-final-back/internal/services/orders"
	productssvc "github.com/bnursik/aitu-ad-final-back/internal/services/products"
	userssvc "github.com/bnursik/aitu-ad-final-back/internal/services/users"
)

func Build(cfg *config.Config) (*Container, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := db.Connect(ctx, cfg.MongoURI)
	if err != nil {
		return nil, err
	}
	dbase := client.Database(cfg.DBName)

	usersRepo := mongorepo.NewUsersRepo(dbase)
	_ = usersRepo.EnsureIndexes(context.Background())

	jwtIssuer := middleware.NewJWT(cfg.JWTSecret, 24*time.Hour)

	authSvc := userssvc.New(usersRepo, jwtIssuer)

	authHandler := handlers.NewAuthHandler(authSvc)
	categoriesRepo := mongorepo.NewCategoriesRepo(dbase)
	productsCounter := mongorepo.NewProductsCounterRepo(dbase)
	categoriesSvc := categoriessvc.New(categoriesRepo, productsCounter)
	categoriesHandler := handlers.NewCategoriesHandler(categoriesSvc)

	productsRepo := mongorepo.NewProductsRepo(dbase)
	productsSvc := productssvc.New(productsRepo)
	productsHandler := handlers.NewProductsHandler(productsSvc)

	ordersRepo := mongorepo.NewOrdersRepo(dbase)
	ordersSvc := orderssvc.New(ordersRepo)
	ordersHandler := handlers.NewOrdersHandler(ordersSvc)

	return &Container{
		Auth: authHandler,
		Shutdown: func(ctx context.Context) error {
			return client.Disconnect(ctx)
		},
		JWT:        jwtIssuer,
		Now:        func() time.Time { return time.Now().UTC() },
		Categories: categoriesHandler,
		Products:   productsHandler,
		Orders:     ordersHandler,
	}, nil
}
