package handlers

import (
	"errors"
	"net/http"
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

type DateRangeRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

type YearRequest struct {
	Year int `json:"year" binding:"required"`
}

// GetSalesStatsByDateRange godoc
// @Summary Get sales statistics by date range (admin only)
// @Tags Admin Statistics
// @Accept json
// @Produce json
// @Param body body DateRangeRequest true "Date Range"
// @Success 200 {object} statistics.SalesStatistics
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/statistics/sales/date-range [post]
func (h *StatisticsHandler) GetSalesStatsByDateRange(c *gin.Context) {
	var req DateRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, use YYYY-MM-DD"})
		return
	}

	// Set end date to end of day
	endDate = endDate.Add(24*time.Hour - time.Second)

	stats, err := h.svc.GetSalesStatsByDateRange(c.Request.Context(), statistics.DateRangeFilter{
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		if errors.Is(err, statistics.ErrInvalidDateRange) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date range"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetSalesStatsByYear godoc
// @Summary Get sales statistics by year (admin only)
// @Tags Admin Statistics
// @Accept json
// @Produce json
// @Param body body YearRequest true "Year"
// @Success 200 {object} statistics.SalesStatistics
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/statistics/sales/year [post]
func (h *StatisticsHandler) GetSalesStatsByYear(c *gin.Context) {
	var req YearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	stats, err := h.svc.GetSalesStatsByYear(c.Request.Context(), statistics.YearFilter{
		Year: req.Year,
	})
	if err != nil {
		if errors.Is(err, statistics.ErrInvalidYear) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetSalesStatsAll godoc
// @Summary Get all sales statistics (admin only)
// @Tags Admin Statistics
// @Produce json
// @Success 200 {object} statistics.SalesStatistics
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/statistics/sales [get]
func (h *StatisticsHandler) GetSalesStatsAll(c *gin.Context) {
	stats, err := h.svc.GetSalesStatsAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetProductsStatsByDateRange godoc
// @Summary Get products statistics by date range (admin only)
// @Tags Admin Statistics
// @Accept json
// @Produce json
// @Param body body DateRangeRequest true "Date Range"
// @Success 200 {object} statistics.ProductStatistics
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/statistics/products/date-range [post]
func (h *StatisticsHandler) GetProductsStatsByDateRange(c *gin.Context) {
	var req DateRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, use YYYY-MM-DD"})
		return
	}

	// Set end date to end of day
	endDate = endDate.Add(24*time.Hour - time.Second)

	stats, err := h.svc.GetProductsStatsByDateRange(c.Request.Context(), statistics.DateRangeFilter{
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		if errors.Is(err, statistics.ErrInvalidDateRange) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date range"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetProductsStatsByYear godoc
// @Summary Get products statistics by year (admin only)
// @Tags Admin Statistics
// @Accept json
// @Produce json
// @Param body body YearRequest true "Year"
// @Success 200 {object} statistics.ProductStatistics
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/statistics/products/year [post]
func (h *StatisticsHandler) GetProductsStatsByYear(c *gin.Context) {
	var req YearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	stats, err := h.svc.GetProductsStatsByYear(c.Request.Context(), statistics.YearFilter{
		Year: req.Year,
	})
	if err != nil {
		if errors.Is(err, statistics.ErrInvalidYear) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetProductsStatsAll godoc
// @Summary Get all products statistics (admin only)
// @Tags Admin Statistics
// @Produce json
// @Success 200 {object} statistics.ProductStatistics
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/statistics/products [get]
func (h *StatisticsHandler) GetProductsStatsAll(c *gin.Context) {
	stats, err := h.svc.GetProductsStatsAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
