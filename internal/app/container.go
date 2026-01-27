package app

import (
	"context"
	"time"

	"github.com/bnursik/aitu-ad-final-back/internal/http/handlers"
	"github.com/bnursik/aitu-ad-final-back/internal/http/middleware"
)

type Container struct {
	Auth *handlers.AuthHandler

	Shutdown func(ctx context.Context) error

	Now        func() time.Time
	Categories *handlers.CategoriesHandler
	JWT        *middleware.JWT
	Products   *handlers.ProductsHandler
	Orders     *handlers.OrdersHandler
}
