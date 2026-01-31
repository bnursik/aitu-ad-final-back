package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/statistics"
	"github.com/gin-gonic/gin"
)

type StatisticsHandler struct {
	svc statistics.Service
}

func NewStatisticsHandler(svc statistics.Service) *StatisticsHandler {
	return &StatisticsHandler{svc: svc}
}

// parseStatsParams reads query params: year OR start&end. If year is present (optionally with start), use year.
func parseStatsParams(c *gin.Context) (useYear bool, year int, startDate, endDate time.Time, err error) {
	yearStr := c.Query("year")
	startStr := c.Query("start")
	endStr := c.Query("end")

	// If year is provided, use by-year (even if start/end also provided)
	if yearStr != "" {
		y, parseErr := strconv.Atoi(yearStr)
		if parseErr != nil || y < 1900 || y > 2100 {
			return false, 0, time.Time{}, time.Time{}, statistics.ErrInvalidYear
		}
		return true, y, time.Time{}, time.Time{}, nil
	}

	// Else use date range if both start and end provided
	if startStr != "" && endStr != "" {
		start, err1 := time.Parse("2006-01-02", startStr)
		if err1 != nil {
			return false, 0, time.Time{}, time.Time{}, err1
		}
		end, err2 := time.Parse("2006-01-02", endStr)
		if err2 != nil {
			return false, 0, time.Time{}, time.Time{}, err2
		}
		end = end.Add(24*time.Hour - time.Second)
		if start.After(end) {
			return false, 0, time.Time{}, time.Time{}, statistics.ErrInvalidDateRange
		}
		return false, 0, start, end, nil
	}

	// Neither valid year nor full date range: treat as "all" (no filter)
	return false, 0, time.Time{}, time.Time{}, nil
}

// GetSalesStats godoc
// @Summary Get sales statistics (admin only)
// @Tags Admin Stats
// @Produce json
// @Param year query int false "Filter by year (e.g. 2024). If year and start both present, year wins."
// @Param start query string false "Start date YYYY-MM-DD (use with end for date range)"
// @Param end query string false "End date YYYY-MM-DD (use with start for date range)"
// @Success 200 {object} statistics.SalesStatistics
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/stats/sales [get]
func (h *StatisticsHandler) GetSalesStats(c *gin.Context) {
	useYear, year, startDate, endDate, err := parseStatsParams(c)
	if err != nil {
		if errors.Is(err, statistics.ErrInvalidYear) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
			return
		}
		if errors.Is(err, statistics.ErrInvalidDateRange) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date range"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start or end date format, use YYYY-MM-DD"})
		return
	}

	var stats statistics.SalesStatistics
	if useYear {
		stats, err = h.svc.GetSalesStatsByYear(c.Request.Context(), statistics.YearFilter{Year: year})
	} else if !startDate.IsZero() && !endDate.IsZero() {
		stats, err = h.svc.GetSalesStatsByDateRange(c.Request.Context(), statistics.DateRangeFilter{StartDate: startDate, EndDate: endDate})
	} else {
		stats, err = h.svc.GetSalesStatsAll(c.Request.Context())
	}
	if err != nil {
		if errors.Is(err, statistics.ErrInvalidDateRange) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date range"})
			return
		}
		if errors.Is(err, statistics.ErrInvalidYear) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetProductsStats godoc
// @Summary Get products statistics (admin only)
// @Tags Admin Stats
// @Produce json
// @Param year query int false "Filter by year (e.g. 2024). If year and start both present, year wins."
// @Param start query string false "Start date YYYY-MM-DD (use with end for date range)"
// @Param end query string false "End date YYYY-MM-DD (use with start for date range)"
// @Success 200 {object} statistics.ProductStatistics
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/stats/products [get]
func (h *StatisticsHandler) GetProductsStats(c *gin.Context) {
	useYear, year, startDate, endDate, err := parseStatsParams(c)
	if err != nil {
		if errors.Is(err, statistics.ErrInvalidYear) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
			return
		}
		if errors.Is(err, statistics.ErrInvalidDateRange) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date range"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start or end date format, use YYYY-MM-DD"})
		return
	}

	var stats statistics.ProductStatistics
	if useYear {
		stats, err = h.svc.GetProductsStatsByYear(c.Request.Context(), statistics.YearFilter{Year: year})
	} else if !startDate.IsZero() && !endDate.IsZero() {
		stats, err = h.svc.GetProductsStatsByDateRange(c.Request.Context(), statistics.DateRangeFilter{StartDate: startDate, EndDate: endDate})
	} else {
		stats, err = h.svc.GetProductsStatsAll(c.Request.Context())
	}
	if err != nil {
		if errors.Is(err, statistics.ErrInvalidDateRange) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date range"})
			return
		}
		if errors.Is(err, statistics.ErrInvalidYear) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
