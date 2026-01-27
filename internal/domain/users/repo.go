package users

import "context"

type Repo interface {
	Insert(ctx context.Context, u User) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
}
