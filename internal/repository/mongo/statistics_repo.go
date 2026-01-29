package mongorepo

import (
	"context"
	"fmt"
	"time"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/statistics"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatisticsRepo struct {
	ordersCol     *mongo.Collection
	productsCol   *mongo.Collection
	categoriesCol *mongo.Collection
}

func NewStatisticsRepo(db *mongo.Database) *StatisticsRepo {
	return &StatisticsRepo{
		ordersCol:     db.Collection("orders"),
		productsCol:   db.Collection("products"),
		categoriesCol: db.Collection("categories"),
	}
}

func (r *StatisticsRepo) GetSalesStatsByDateRange(ctx context.Context, filter statistics.DateRangeFilter) (statistics.SalesStatistics, error) {
	dateFilter := bson.M{
		"createdAt": bson.M{
			"$gte": filter.StartDate,
			"$lte": filter.EndDate,
		},
	}
	return r.getSalesStats(ctx, dateFilter)
}

func (r *StatisticsRepo) GetSalesStatsByYear(ctx context.Context, filter statistics.YearFilter) (statistics.SalesStatistics, error) {
	startDate := time.Date(filter.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(filter.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)

	dateFilter := bson.M{
		"createdAt": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
	}
	return r.getSalesStats(ctx, dateFilter)
}

func (r *StatisticsRepo) GetSalesStatsAll(ctx context.Context) (statistics.SalesStatistics, error) {
	return r.getSalesStats(ctx, bson.M{})
}

func (r *StatisticsRepo) getSalesStats(ctx context.Context, filter bson.M) (statistics.SalesStatistics, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$facet", Value: bson.M{
			"statusCounts": []bson.M{
				{"$group": bson.M{
					"_id":   "$status",
					"count": bson.M{"$sum": 1},
				}},
			},
			"totals": []bson.M{
				{"$unwind": "$items"},
				{"$lookup": bson.M{
					"from":         "products",
					"localField":   "items.productId",
					"foreignField": "_id",
					"as":           "product",
				}},
				{"$unwind": bson.M{"path": "$product", "preserveNullAndEmptyArrays": true}},
				{"$group": bson.M{
					"_id":          nil,
					"totalOrders":  bson.M{"$addToSet": "$_id"},
					"totalRevenue": bson.M{"$sum": bson.M{"$multiply": []interface{}{"$items.quantity", bson.M{"$ifNull": []interface{}{"$product.price", 0}}}}},
				}},
				{"$project": bson.M{
					"totalOrders":  bson.M{"$size": "$totalOrders"},
					"totalRevenue": 1,
				}},
			},
		}}},
	}

	cur, err := r.ordersCol.Aggregate(ctx, pipeline)
	if err != nil {
		return statistics.SalesStatistics{}, fmt.Errorf("aggregate sales stats: %w", err)
	}
	defer cur.Close(ctx)

	var results []struct {
		StatusCounts []struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		} `bson:"statusCounts"`
		Totals []struct {
			TotalOrders  int64   `bson:"totalOrders"`
			TotalRevenue float64 `bson:"totalRevenue"`
		} `bson:"totals"`
	}

	if err := cur.All(ctx, &results); err != nil {
		return statistics.SalesStatistics{}, fmt.Errorf("decode sales stats: %w", err)
	}

	stats := statistics.SalesStatistics{}

	if len(results) > 0 {
		for _, sc := range results[0].StatusCounts {
			switch sc.ID {
			case "pending":
				stats.PendingOrders = sc.Count
			case "shipped":
				stats.ShippedOrders = sc.Count
			case "delivered":
				stats.DeliveredOrders = sc.Count
			case "cancelled":
				stats.CancelledOrders = sc.Count
			}
		}

		if len(results[0].Totals) > 0 {
			stats.TotalOrders = results[0].Totals[0].TotalOrders
			stats.TotalRevenue = results[0].Totals[0].TotalRevenue
		}
	}

	// Calculate total from status counts if totals pipeline didn't return results
	if stats.TotalOrders == 0 {
		stats.TotalOrders = stats.PendingOrders + stats.ShippedOrders + stats.DeliveredOrders + stats.CancelledOrders
	}

	if stats.TotalOrders > 0 {
		stats.AverageOrder = stats.TotalRevenue / float64(stats.TotalOrders)
	}

	return stats, nil
}

func (r *StatisticsRepo) GetProductsStatsByDateRange(ctx context.Context, filter statistics.DateRangeFilter) (statistics.ProductStatistics, error) {
	dateFilter := bson.M{
		"createdAt": bson.M{
			"$gte": filter.StartDate,
			"$lte": filter.EndDate,
		},
	}
	return r.getProductsStats(ctx, dateFilter)
}

func (r *StatisticsRepo) GetProductsStatsByYear(ctx context.Context, filter statistics.YearFilter) (statistics.ProductStatistics, error) {
	startDate := time.Date(filter.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(filter.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)

	dateFilter := bson.M{
		"createdAt": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
	}
	return r.getProductsStats(ctx, dateFilter)
}

func (r *StatisticsRepo) GetProductsStatsAll(ctx context.Context) (statistics.ProductStatistics, error) {
	return r.getProductsStats(ctx, bson.M{})
}

func (r *StatisticsRepo) getProductsStats(ctx context.Context, filter bson.M) (statistics.ProductStatistics, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$facet", Value: bson.M{
			"products": []bson.M{
				{"$group": bson.M{
					"_id":           nil,
					"totalProducts": bson.M{"$sum": 1},
					"totalStock":    bson.M{"$sum": "$stock"},
					"outOfStock":    bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$lte": []interface{}{"$stock", 0}}, 1, 0}}},
				}},
			},
			"reviews": []bson.M{
				{"$unwind": bson.M{"path": "$reviews", "preserveNullAndEmptyArrays": true}},
				{"$group": bson.M{
					"_id":          nil,
					"totalReviews": bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$ifNull": []interface{}{"$reviews", false}}, 1, 0}}},
					"totalRating":  bson.M{"$sum": bson.M{"$ifNull": []interface{}{"$reviews.rating", 0}}},
					"reviewCount":  bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$ifNull": []interface{}{"$reviews.rating", false}}, 1, 0}}},
				}},
			},
		}}},
	}

	cur, err := r.productsCol.Aggregate(ctx, pipeline)
	if err != nil {
		return statistics.ProductStatistics{}, fmt.Errorf("aggregate products stats: %w", err)
	}
	defer cur.Close(ctx)

	var results []struct {
		Products []struct {
			TotalProducts int64 `bson:"totalProducts"`
			TotalStock    int64 `bson:"totalStock"`
			OutOfStock    int64 `bson:"outOfStock"`
		} `bson:"products"`
		Reviews []struct {
			TotalReviews int64 `bson:"totalReviews"`
			TotalRating  int64 `bson:"totalRating"`
			ReviewCount  int64 `bson:"reviewCount"`
		} `bson:"reviews"`
	}

	if err := cur.All(ctx, &results); err != nil {
		return statistics.ProductStatistics{}, fmt.Errorf("decode products stats: %w", err)
	}

	stats := statistics.ProductStatistics{}

	if len(results) > 0 {
		if len(results[0].Products) > 0 {
			stats.TotalProducts = results[0].Products[0].TotalProducts
			stats.TotalStock = results[0].Products[0].TotalStock
			stats.OutOfStock = results[0].Products[0].OutOfStock
		}

		if len(results[0].Reviews) > 0 {
			stats.TotalReviews = results[0].Reviews[0].TotalReviews
			if results[0].Reviews[0].ReviewCount > 0 {
				stats.AverageRating = float64(results[0].Reviews[0].TotalRating) / float64(results[0].Reviews[0].ReviewCount)
			}
		}
	}

	// Get total categories count
	catCount, err := r.categoriesCol.CountDocuments(ctx, bson.M{})
	if err != nil {
		return statistics.ProductStatistics{}, fmt.Errorf("count categories: %w", err)
	}
	stats.TotalCategories = catCount

	return stats, nil
}
