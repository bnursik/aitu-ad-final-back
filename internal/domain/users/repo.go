package users

import "context"

type Repo interface {
	Insert(ctx context.Context, u User) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByID(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, id string, in UpdateProfileInput) (User, error)
	GetAll(ctx context.Context) ([]User, error)
}
