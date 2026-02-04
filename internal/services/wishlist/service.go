package wishlistsvc

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/products"
	"github.com/bnursik/aitu-ad-final-back/internal/domain/wishlist"
)

type Service struct {
	repo         wishlist.Repo
	productsRepo products.Repo
	now          func() time.Time
}

func New(repo wishlist.Repo, productsRepo products.Repo) *Service {
	return &Service{
		repo:         repo,
		productsRepo: productsRepo,
		now:          func() time.Time { return time.Now().UTC() },
	}
}

var _ wishlist.Service = (*Service)(nil)

func (s *Service) Add(ctx context.Context, userID string, in wishlist.AddItemInput) (wishlist.WishlistItem, error) {
	uid := strings.TrimSpace(userID)
	if uid == "" {
		return wishlist.WishlistItem{}, wishlist.ErrInvalidID
	}

	productID := strings.TrimSpace(in.ProductID)
	if productID == "" {
		return wishlist.WishlistItem{}, wishlist.ErrInvalidProduct
	}

	prod, err := s.productsRepo.GetByID(ctx, productID)
	if err != nil {
		if errors.Is(err, products.ErrNotFound) {
			return wishlist.WishlistItem{}, wishlist.ErrInvalidProduct
		}
		return wishlist.WishlistItem{}, err
	}
	if prod.Stock < 1 {
		return wishlist.WishlistItem{}, wishlist.ErrProductOutOfStock
	}

	item := wishlist.WishlistItem{
		UserID:    uid,
		ProductID: productID,
		CreatedAt: s.now(),
	}

	return s.repo.Add(ctx, item)
}

func (s *Service) List(ctx context.Context, userID string, f wishlist.ListFilter) ([]wishlist.WishlistItem, int64, error) {
	uid := strings.TrimSpace(userID)
	if uid == "" {
		return nil, 0, wishlist.ErrInvalidID
	}

	items, err := s.repo.List(ctx, uid, f)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx, uid)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (s *Service) Delete(ctx context.Context, userID string, id string) error {
	uid := strings.TrimSpace(userID)
	if uid == "" {
		return wishlist.ErrInvalidID
	}

	itemID := strings.TrimSpace(id)
	if itemID == "" {
		return wishlist.ErrInvalidID
	}

	return s.repo.Delete(ctx, uid, itemID)
}
