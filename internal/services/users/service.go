package userssvc

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/users"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo   users.Repo
	tokens users.TokenIssuer
	now    func() time.Time
}

func New(repo users.Repo, tokens users.TokenIssuer) *Service {
	return &Service{
		repo:   repo,
		tokens: tokens,
		now:    func() time.Time { return time.Now().UTC() },
	}
}

var _ users.Service = (*Service)(nil)

func (s *Service) Register(ctx context.Context, name, email, password string) (string, users.PublicUser, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))

	if len(name) < 2 || len(name) > 60 {
		return "", users.PublicUser{}, fmt.Errorf("invalid name")
	}
	if email == "" || !strings.Contains(email, "@") {
		return "", users.PublicUser{}, users.ErrInvalidEmail
	}
	if len(password) < 6 || len(password) > 72 {
		return "", users.PublicUser{}, users.ErrInvalidPassword
	}

	_, err := s.repo.FindByEmail(ctx, email)
	if err == nil {
		return "", users.PublicUser{}, users.ErrEmailTaken
	}
	if err != nil && !errors.Is(err, users.ErrUserNotFound) {
		return "", users.PublicUser{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", users.PublicUser{}, fmt.Errorf("hash password: %w", err)
	}

	created, err := s.repo.Insert(ctx, users.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         users.RoleUser,
		CreatedAt:    s.now(),
	})
	if err != nil {
		if errors.Is(err, users.ErrEmailTaken) {
			return "", users.PublicUser{}, users.ErrEmailTaken
		}
		return "", users.PublicUser{}, err
	}

	token, err := s.tokens.IssueAccessToken(created.ID, string(created.Role))
	if err != nil {
		return "", users.PublicUser{}, err
	}

	return token, created.Public(), nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, users.PublicUser, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return "", users.PublicUser{}, users.ErrInvalidCreds
		}
		return "", users.PublicUser{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", users.PublicUser{}, users.ErrInvalidCreds
	}

	token, err := s.tokens.IssueAccessToken(u.ID, string(u.Role))
	if err != nil {
		return "", users.PublicUser{}, err
	}

	return token, u.Public(), nil
}


func (s *Service) AdminRegister(ctx context.Context, name, email, password string) (string, users.PublicUser, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))

	if len(name) < 2 || len(name) > 60 {
		return "", users.PublicUser{}, fmt.Errorf("invalid name")
	}
	if email == "" || !strings.Contains(email, "@") {
		return "", users.PublicUser{}, users.ErrInvalidEmail
	}
	if len(password) < 6 || len(password) > 72 {
		return "", users.PublicUser{}, users.ErrInvalidPassword
	}

	_, err := s.repo.FindByEmail(ctx, email)
	if err == nil {
		return "", users.PublicUser{}, users.ErrEmailTaken
	}
	if err != nil && !errors.Is(err, users.ErrUserNotFound) {
		return "", users.PublicUser{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", users.PublicUser{}, fmt.Errorf("hash password: %w", err)
	}

	created, err := s.repo.Insert(ctx, users.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         users.RoleAdmin,
		CreatedAt:    s.now(),
	})
	if err != nil {
		if errors.Is(err, users.ErrEmailTaken) {
			return "", users.PublicUser{}, users.ErrEmailTaken
		}
		return "", users.PublicUser{}, err
	}

	token, err := s.tokens.IssueAccessToken(created.ID, string(created.Role))
	if err != nil {
		return "", users.PublicUser{}, err
	}

	return token, created.Public(), nil
}

func (s *Service) GetProfile(ctx context.Context, userID string) (users.PublicUser, error) {
	uid := strings.TrimSpace(userID)
	if uid == "" {
		return users.PublicUser{}, users.ErrUserNotFound
	}

	u, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		return users.PublicUser{}, err
	}

	return u.Public(), nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID string, in users.UpdateProfileInput) (users.PublicUser, error) {
	uid := strings.TrimSpace(userID)
	if uid == "" {
		return users.PublicUser{}, users.ErrUserNotFound
	}

	if in.Name != nil {
		n := strings.TrimSpace(*in.Name)
		if len(n) < 2 || len(n) > 60 {
			return users.PublicUser{}, fmt.Errorf("invalid name")
		}
		in.Name = &n
	}
	if in.Address != nil {
		a := strings.TrimSpace(*in.Address)
		in.Address = &a
	}
	if in.Phone != nil {
		p := strings.TrimSpace(*in.Phone)
		in.Phone = &p
	}
	if in.Bio != nil {
		b := strings.TrimSpace(*in.Bio)
		if len(b) > 500 {
			return users.PublicUser{}, fmt.Errorf("invalid bio")
		}
		in.Bio = &b
	}

	u, err := s.repo.Update(ctx, uid, in)
	if err != nil {
		return users.PublicUser{}, err
	}

	return u.Public(), nil
}

func (s *Service) GetAllUsers(ctx context.Context) ([]users.PublicUser, error) {
	all, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]users.PublicUser, len(all))
	for i, u := range all {
		result[i] = u.Public()
	}

	return result, nil
}