package statistics

import "context"

type Service interface {
	GetSalesStatsByDateRange(ctx context.Context, filter DateRangeFilter) (SalesStatistics, error)
	GetSalesStatsByYear(ctx context.Context, filter YearFilter) (SalesStatistics, error)
	GetSalesStatsAll(ctx context.Context) (SalesStatistics, error)

	GetProductsStatsByDateRange(ctx context.Context, filter DateRangeFilter) (ProductStatistics, error)
	GetProductsStatsByYear(ctx context.Context, filter YearFilter) (ProductStatistics, error)
	GetProductsStatsAll(ctx context.Context) (ProductStatistics, error)
}
