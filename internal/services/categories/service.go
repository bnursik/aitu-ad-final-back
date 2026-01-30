package categoriessvc

import (
	"context"
	"strings"
	"time"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/categories"
)

type Service struct {
	repo     categories.Repo
	products categories.ProductsCounter
	now      func() time.Time
}

func New(repo categories.Repo, products categories.ProductsCounter) *Service {
	return &Service{
		repo:     repo,
		products: products,
		now:      func() time.Time { return time.Now().UTC() },
	}
}

var _ categories.Service = (*Service)(nil)

func (s *Service) List(ctx context.Context, f categories.ListFilter) ([]categories.Category, error) {
	return s.repo.List(ctx, f)
}

func (s *Service) Get(ctx context.Context, id string) (categories.Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, in categories.CreateInput) (categories.Category, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return categories.Category{}, categories.ErrInvalidName
	}

	now := s.now()
	c := categories.Category{
		Name:        name,
		Description: strings.TrimSpace(in.Description),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return s.repo.Create(ctx, c)
}

func (s *Service) Update(ctx context.Context, id string, in categories.UpdateInput) (categories.Category, error) {
	if in.Name != nil {
		n := strings.TrimSpace(*in.Name)
		if n == "" {
			return categories.Category{}, categories.ErrInvalidName
		}
		in.Name = &n
	}
	if in.Description != nil {
		d := strings.TrimSpace(*in.Description)
		in.Description = &d
	}

	return s.repo.Update(ctx, id, in)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	n, err := s.products.CountByCategoryID(ctx, id)
	if err != nil {
		return err
	}
	if n > 0 {
		return categories.ErrHasProducts
	}
	return s.repo.Delete(ctx, id)
}
