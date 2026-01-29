package statisticssvc

import (
	"context"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/statistics"
)

type Service struct {
	repo statistics.Repo
}

func New(repo statistics.Repo) *Service {
	return &Service{repo: repo}
}

var _ statistics.Service = (*Service)(nil)

func (s *Service) GetSalesStatsByDateRange(ctx context.Context, filter statistics.DateRangeFilter) (statistics.SalesStatistics, error) {
	if filter.StartDate.After(filter.EndDate) {
		return statistics.SalesStatistics{}, statistics.ErrInvalidDateRange
	}
	return s.repo.GetSalesStatsByDateRange(ctx, filter)
}

func (s *Service) GetSalesStatsByYear(ctx context.Context, filter statistics.YearFilter) (statistics.SalesStatistics, error) {
	if filter.Year < 1900 || filter.Year > 2100 {
		return statistics.SalesStatistics{}, statistics.ErrInvalidYear
	}
	return s.repo.GetSalesStatsByYear(ctx, filter)
}

func (s *Service) GetSalesStatsAll(ctx context.Context) (statistics.SalesStatistics, error) {
	return s.repo.GetSalesStatsAll(ctx)
}

func (s *Service) GetProductsStatsByDateRange(ctx context.Context, filter statistics.DateRangeFilter) (statistics.ProductStatistics, error) {
	if filter.StartDate.After(filter.EndDate) {
		return statistics.ProductStatistics{}, statistics.ErrInvalidDateRange
	}
	return s.repo.GetProductsStatsByDateRange(ctx, filter)
}

func (s *Service) GetProductsStatsByYear(ctx context.Context, filter statistics.YearFilter) (statistics.ProductStatistics, error) {
	if filter.Year < 1900 || filter.Year > 2100 {
		return statistics.ProductStatistics{}, statistics.ErrInvalidYear
	}
	return s.repo.GetProductsStatsByYear(ctx, filter)
}

func (s *Service) GetProductsStatsAll(ctx context.Context) (statistics.ProductStatistics, error) {
	return s.repo.GetProductsStatsAll(ctx)
}
