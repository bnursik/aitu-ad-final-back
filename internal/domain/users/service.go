package users

import "context"

type TokenIssuer interface {
	IssueAccessToken(userID string, role string) (string, error)
}

type Service interface {
	Register(ctx context.Context, name, email, password string) (token string, user PublicUser, err error)
	Login(ctx context.Context, email, password string) (token string, user PublicUser, err error)
	AdminRegister(ctx context.Context, name, email, password string) (token string, user PublicUser, err error)
	GetProfile(ctx context.Context, userID string) (PublicUser, error)
	UpdateProfile(ctx context.Context, userID string, in UpdateProfileInput) (PublicUser, error)
}
