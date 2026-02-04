package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/products"
	"github.com/bnursik/aitu-ad-final-back/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

type ProductsHandler struct {
	svc products.Service
}

func NewProductsHandler(svc products.Service) *ProductsHandler {
	return &ProductsHandler{svc: svc}
}

type CreateProductRequest struct {
	CategoryID  string  `json:"categoryId" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Stock       int64   `json:"stock" binding:"required"`
}

type UpdateProductRequest struct {
	CategoryID  *string  `json:"categoryId"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Stock       *int64   `json:"stock"`
}

type AddReviewRequest struct {
	Rating  int64  `json:"rating" binding:"required"`
	Comment string `json:"comment"`
}

// ListProducts godoc
// @Summary List products
// @Tags Products
// @Produce json
// @Param categoryId query string false "Category ID (ObjectId hex)"
// @Param offset query int true "Offset for pagination"
// @Param limit query int true "Limit for pagination"
// @Success 200 {array} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /products [get]
func (h *ProductsHandler) List(c *gin.Context) {
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")

	if offsetStr == "" || limitStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "offset and limit are required"})
		return
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	var f products.ListFilter
	f.Offset = offset
	f.Limit = limit
	if v := c.Query("categoryId"); v != "" {
		f.CategoryID = &v
	}

	items, total, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		if errors.Is(err, products.ErrInvalidCategory) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid categoryId"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	out := make([]gin.H, 0, len(items))
	for _, it := range items {
		out = append(out, gin.H{
			"id":          it.ID,
			"categoryId":  it.CategoryID,
			"name":        it.Name,
			"description": it.Description,
			"price":       it.Price,
			"stock":       it.Stock,
			"createdAt":   it.CreatedAt,
			"updatedAt":   it.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"items":  out,
		"total":  total,
		"offset": offset,
		"limit":  limit,
	})
}

// GetProduct godoc
// @Summary Get product by ID
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func (h *ProductsHandler) Get(c *gin.Context) {
	id := c.Param("id")

	it, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, products.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		case errors.Is(err, products.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	reviews := make([]gin.H, 0, len(it.Reviews))
	for _, r := range it.Reviews {
		reviews = append(reviews, gin.H{
			"id":        r.ID,
			"userId":    r.UserID,
			"rating":    r.Rating,
			"comment":   r.Comment,
			"createdAt": r.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          it.ID,
		"categoryId":  it.CategoryID,
		"name":        it.Name,
		"description": it.Description,
		"price":       it.Price,
		"stock":       it.Stock,
		"createdAt":   it.CreatedAt,
		"updatedAt":   it.UpdatedAt,
		"reviews":     reviews,
	})
}

// CreateProduct godoc
// @Summary Create product
// @Tags Admin Products
// @Accept json
// @Produce json
// @Param body body CreateProductRequest true "Product"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /admin/products [post]
func (h *ProductsHandler) Create(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	it, err := h.svc.Create(c.Request.Context(), products.CreateInput{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	})
	if err != nil {
		switch {
		case errors.Is(err, products.ErrInvalidCategory):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid categoryId"})
		case errors.Is(err, products.ErrInvalidName):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
		case errors.Is(err, products.ErrInvalidPrice):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
		case errors.Is(err, products.ErrInvalidStock):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stock"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          it.ID,
		"categoryId":  it.CategoryID,
		"name":        it.Name,
		"description": it.Description,
		"price":       it.Price,
		"stock":       it.Stock,
		"createdAt":   it.CreatedAt,
		"updatedAt":   it.UpdatedAt,
	})
}

// UpdateProduct godoc
// @Summary Update product
// @Tags Admin Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param body body UpdateProductRequest true "Patch"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/products/{id} [put]
func (h *ProductsHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	it, err := h.svc.Update(c.Request.Context(), id, products.UpdateInput{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	})
	if err != nil {
		switch {
		case errors.Is(err, products.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		case errors.Is(err, products.ErrInvalidCategory):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid categoryId"})
		case errors.Is(err, products.ErrInvalidName):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
		case errors.Is(err, products.ErrInvalidPrice):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
		case errors.Is(err, products.ErrInvalidStock):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stock"})
		case errors.Is(err, products.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          it.ID,
		"categoryId":  it.CategoryID,
		"name":        it.Name,
		"description": it.Description,
		"price":       it.Price,
		"stock":       it.Stock,
		"createdAt":   it.CreatedAt,
		"updatedAt":   it.UpdatedAt,
	})
}

// DeleteProduct godoc
// @Summary Delete product
// @Tags Admin Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/products/{id} [delete]
func (h *ProductsHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, products.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		case errors.Is(err, products.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		case errors.Is(err, products.ErrCannotDeleteProduct):
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete product with stock; stock must be less than 1"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// AddReview godoc
// @Summary Add product review (auth required)
// @Tags Reviews
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param body body AddReviewRequest true "Review"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id}/reviews [post]
func (h *ProductsHandler) AddReview(c *gin.Context) {
	productID := c.Param("id")

	var req AddReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	userIDVal, ok := c.Get(middleware.CtxUserID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, _ := userIDVal.(string)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	rev, err := h.svc.AddReview(c.Request.Context(), productID, products.AddReviewInput{
		UserID:  userID,
		Rating:  req.Rating,
		Comment: req.Comment,
	})
	if err != nil {
		switch {
		case errors.Is(err, products.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		case errors.Is(err, products.ErrInvalidRating):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rating"})
		case errors.Is(err, products.ErrInvalidComment):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment"})
		case errors.Is(err, products.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        rev.ID,
		"userId":    rev.UserID,
		"rating":    rev.Rating,
		"comment":   rev.Comment,
		"createdAt": rev.CreatedAt,
	})
}

// DeleteReview godoc
// @Summary Delete product review (auth required)
// @Tags Reviews
// @Produce json
// @Param id path string true "Product ID"
// @Param reviewId path string true "Review ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id}/reviews/{reviewId} [delete]
func (h *ProductsHandler) DeleteReview(c *gin.Context) {
	productID := c.Param("id")
	reviewID := c.Param("reviewId")

	if _, ok := c.Get(middleware.CtxUserID); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := h.svc.DeleteReview(c.Request.Context(), productID, reviewID)
	if err != nil {
		switch {
		case errors.Is(err, products.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		case errors.Is(err, products.ErrInvalidReviewID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reviewId"})
		case errors.Is(err, products.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
