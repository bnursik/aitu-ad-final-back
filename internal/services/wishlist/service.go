package wishlistsvc

import (
	"context"
	"strings"
	"time"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/wishlist"
)

type Service struct {
	repo wishlist.Repo
	now  func() time.Time
}

func New(repo wishlist.Repo) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
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

	item := wishlist.WishlistItem{
		UserID:    uid,
		ProductID: productID,
		CreatedAt: s.now(),
	}

	return s.repo.Add(ctx, item)
}

func (s *Service) List(ctx context.Context, userID string, f wishlist.ListFilter) ([]wishlist.WishlistItem, error) {
	uid := strings.TrimSpace(userID)
	if uid == "" {
		return nil, wishlist.ErrInvalidID
	}
	return s.repo.List(ctx, uid, f)
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
