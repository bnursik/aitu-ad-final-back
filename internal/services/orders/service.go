package orderssvc

import (
	"context"
	"strings"
	"time"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/orders"
)

type Service struct {
	repo orders.Repo
	now  func() time.Time
}

func New(repo orders.Repo) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

var _ orders.Service = (*Service)(nil)

func (s *Service) List(ctx context.Context, userID string, isAdmin bool, f orders.ListFilter) ([]orders.Order, error) {
	if isAdmin {
		return s.repo.List(ctx, nil, f)
	}
	uid := strings.TrimSpace(userID)
	return s.repo.List(ctx, &uid, f)
}

func (s *Service) Get(ctx context.Context, id string, userID string, isAdmin bool) (orders.Order, error) {
	if isAdmin {
		return s.repo.GetByID(ctx, id, nil)
	}
	uid := strings.TrimSpace(userID)
	return s.repo.GetByID(ctx, id, &uid)
}

func (s *Service) Create(ctx context.Context, userID string, in orders.CreateInput) (orders.Order, error) {
	uid := strings.TrimSpace(userID)
	if uid == "" {
		return orders.Order{}, orders.ErrForbidden
	}

	if len(in.Items) == 0 {
		return orders.Order{}, orders.ErrInvalidItems
	}

	items := make([]orders.Item, 0, len(in.Items))
	for _, it := range in.Items {
		p := strings.TrimSpace(it.ProductID)
		if p == "" {
			return orders.Order{}, orders.ErrInvalidProduct
		}
		if it.Quantity <= 0 {
			return orders.Order{}, orders.ErrInvalidQty
		}
		items = append(items, orders.Item{ProductID: p, Quantity: it.Quantity})
	}

	now := s.now()
	o := orders.Order{
		UserID:    uid,
		Items:     items,
		Status:    orders.StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return s.repo.Create(ctx, o)
}

func (s *Service) UpdateStatus(ctx context.Context, id string, status orders.Status) (orders.Order, error) {
	switch status {
	case orders.StatusPending, orders.StatusShipped, orders.StatusDelivered, orders.StatusCancelled:
		// ok
	default:
		return orders.Order{}, orders.ErrInvalidStatus
	}
	return s.repo.UpdateStatus(ctx, id, status)
}
