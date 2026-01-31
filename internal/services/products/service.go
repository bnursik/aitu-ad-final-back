package productssvc

import (
	"context"
	"strings"
	"time"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/products"
)

type Service struct {
	repo products.Repo
	now  func() time.Time
}

func New(repo products.Repo) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

var _ products.Service = (*Service)(nil)

func (s *Service) List(ctx context.Context, f products.ListFilter) ([]products.Product, int64, error) {
	items, err := s.repo.List(ctx, f)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx, f)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (s *Service) Get(ctx context.Context, id string) (products.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, in products.CreateInput) (products.Product, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return products.Product{}, products.ErrInvalidName
	}
	if strings.TrimSpace(in.CategoryID) == "" {
		return products.Product{}, products.ErrInvalidCategory
	}
	if in.Price <= 0 {
		return products.Product{}, products.ErrInvalidPrice
	}
	if in.Stock < 0 {
		return products.Product{}, products.ErrInvalidStock
	}

	now := s.now()
	p := products.Product{
		CategoryID:  in.CategoryID,
		Name:        name,
		Description: strings.TrimSpace(in.Description),
		Price:       in.Price,
		Stock:       in.Stock,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return s.repo.Create(ctx, p)
}

func (s *Service) Update(ctx context.Context, id string, in products.UpdateInput) (products.Product, error) {
	if in.Name != nil {
		n := strings.TrimSpace(*in.Name)
		if n == "" {
			return products.Product{}, products.ErrInvalidName
		}
		in.Name = &n
	}
	if in.CategoryID != nil {
		c := strings.TrimSpace(*in.CategoryID)
		if c == "" {
			return products.Product{}, products.ErrInvalidCategory
		}
		in.CategoryID = &c
	}
	if in.Description != nil {
		d := strings.TrimSpace(*in.Description)
		in.Description = &d
	}
	if in.Price != nil && *in.Price <= 0 {
		return products.Product{}, products.ErrInvalidPrice
	}
	if in.Stock != nil && *in.Stock < 0 {
		return products.Product{}, products.ErrInvalidStock
	}

	return s.repo.Update(ctx, id, in)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) AddReview(ctx context.Context, productID string, in products.AddReviewInput) (products.Review, error) {
	if strings.TrimSpace(in.UserID) == "" {
		// можешь заменить на отдельную ошибку, но ок
		return products.Review{}, products.ErrInvalidComment
	}
	if in.Rating < 1 || in.Rating > 5 {
		return products.Review{}, products.ErrInvalidRating
	}
	comment := strings.TrimSpace(in.Comment)
	if len(comment) > 500 {
		return products.Review{}, products.ErrInvalidComment
	}

	r := products.Review{
		UserID:    in.UserID,
		Rating:    in.Rating,
		Comment:   comment,
		CreatedAt: s.now(),
	}
	return s.repo.AddReview(ctx, productID, r)
}

func (s *Service) DeleteReview(ctx context.Context, productID, reviewID string) error {
	return s.repo.DeleteReview(ctx, productID, reviewID)
}
