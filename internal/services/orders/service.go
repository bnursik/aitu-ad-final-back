package orderssvc

import (
	"context"
	"strings"
	"time"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/orders"
	"github.com/bnursik/aitu-ad-final-back/internal/domain/products"
)

type Service struct {
	repo         orders.Repo
	productsRepo products.Repo
	now          func() time.Time
}

func New(repo orders.Repo, productsRepo products.Repo) *Service {
	return &Service{
		repo:         repo,
		productsRepo: productsRepo,
		now:          func() time.Time { return time.Now().UTC() },
	}
}

var _ orders.Service = (*Service)(nil)

func (s *Service) List(ctx context.Context, userID string, isAdmin bool, f orders.ListFilter) ([]orders.Order, error) {
	var (
		list []orders.Order
		err  error
	)

	if isAdmin {
		list, err = s.repo.List(ctx, nil, f)
	} else {
		uid := strings.TrimSpace(userID)
		list, err = s.repo.List(ctx, &uid, f)
	}
	if err != nil {
		return nil, err
	}

	// Compute totals for each returned order
	if err := s.fillTotals(ctx, list); err != nil {
		return nil, err
	}

	return list, nil
}

func (s *Service) Get(ctx context.Context, id string, userID string, isAdmin bool) (orders.Order, error) {
	var (
		o   orders.Order
		err error
	)

	if isAdmin {
		o, err = s.repo.GetByID(ctx, id, nil)
	} else {
		uid := strings.TrimSpace(userID)
		o, err = s.repo.GetByID(ctx, id, &uid)
	}
	if err != nil {
		return orders.Order{}, err
	}

	// compute total for this single order
	if err := s.fillTotals(ctx, []orders.Order{o}); err != nil {
		return orders.Order{}, err
	}

	// fillTotals works on slice elements; for single order easiest is recompute directly:
	total, err := s.computeOne(ctx, &o)
	if err != nil {
		return orders.Order{}, err
	}
	o.TotalPrice = total

	return o, nil
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
	default:
		return orders.Order{}, orders.ErrInvalidStatus
	}
	return s.repo.UpdateStatus(ctx, id, status)
}

// ---- helpers ----

// fillTotals computes unitPrice/lineTotal/totalPrice for all orders in list.
// IMPORTANT: it does NOT write to DB.
func (s *Service) fillTotals(ctx context.Context, list []orders.Order) error {
	// cache product prices within the request
	priceCache := make(map[string]float64, 128)

	for i := range list {
		var total float64

		for j := range list[i].Items {
			pid := list[i].Items[j].ProductID

			price, ok := priceCache[pid]
			if !ok {
				p, err := s.productsRepo.GetByID(ctx, pid)
				if err != nil {
					return orders.ErrInvalidProduct
				}
				price = p.Price
				priceCache[pid] = price
			}

			list[i].Items[j].UnitPrice = price
			list[i].Items[j].LineTotal = price * float64(list[i].Items[j].Quantity)
			total += list[i].Items[j].LineTotal
		}

		list[i].TotalPrice = total
	}

	return nil
}

func (s *Service) computeOne(ctx context.Context, o *orders.Order) (float64, error) {
	priceCache := make(map[string]float64, 32)
	var total float64

	for i := range o.Items {
		pid := o.Items[i].ProductID

		price, ok := priceCache[pid]
		if !ok {
			p, err := s.productsRepo.GetByID(ctx, pid)
			if err != nil {
				return 0, orders.ErrInvalidProduct
			}
			price = p.Price
			priceCache[pid] = price
		}

		o.Items[i].UnitPrice = price
		o.Items[i].LineTotal = price * float64(o.Items[i].Quantity)
		total += o.Items[i].LineTotal
	}
	return total, nil
}
