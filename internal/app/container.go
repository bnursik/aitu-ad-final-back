package app

import (
	"context"
	"time"

	"github.com/bnursik/aitu-ad-final-back/internal/http/handlers"
)

type Container struct {
	Auth *handlers.AuthHandler

	Shutdown func(ctx context.Context) error

	Now func() time.Time
}
