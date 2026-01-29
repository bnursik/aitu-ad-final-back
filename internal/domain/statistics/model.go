package statistics

import "time"

type SalesStatistics struct {
	TotalOrders     int64   `json:"total_orders"`
	TotalRevenue    float64 `json:"total_revenue"`
	AverageOrder    float64 `json:"average_order"`
	PendingOrders   int64   `json:"pending_orders"`
	ShippedOrders   int64   `json:"shipped_orders"`
	DeliveredOrders int64   `json:"delivered_orders"`
	CancelledOrders int64   `json:"cancelled_orders"`
}

type ProductStatistics struct {
	TotalProducts   int64   `json:"total_products"`
	TotalStock      int64   `json:"total_stock"`
	OutOfStock      int64   `json:"out_of_stock"`
	TotalReviews    int64   `json:"total_reviews"`
	AverageRating   float64 `json:"average_rating"`
	TotalCategories int64   `json:"total_categories"`
}

type DateRangeFilter struct {
	StartDate time.Time
	EndDate   time.Time
}

type YearFilter struct {
	Year int
}
