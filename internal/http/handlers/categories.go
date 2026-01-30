package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/categories"
	"github.com/gin-gonic/gin"
)

type CategoriesHandler struct {
	svc categories.Service
}

func NewCategoriesHandler(svc categories.Service) *CategoriesHandler {
	return &CategoriesHandler{svc: svc}
}

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

// ListCategories godoc
// @Summary List categories
// @Tags Categories
// @Produce json
// @Param offset query int true "Offset for pagination"
// @Param limit query int true "Limit for pagination"
// @Success 200 {array} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /categories [get]
func (h *CategoriesHandler) List(c *gin.Context) {
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

	items, err := h.svc.List(c.Request.Context(), categories.ListFilter{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	out := make([]gin.H, 0, len(items))
	for _, it := range items {
		out = append(out, gin.H{
			"id":          it.ID,
			"name":        it.Name,
			"description": it.Description,
			"createdAt":   it.CreatedAt,
			"updatedAt":   it.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, out)
}

// GetCategory godoc
// @Summary Get category by ID
// @Tags Categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /categories/{id} [get]
func (h *CategoriesHandler) Get(c *gin.Context) {
	id := c.Param("id")

	item, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, categories.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		case errors.Is(err, categories.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          item.ID,
		"name":        item.Name,
		"description": item.Description,
		"createdAt":   item.CreatedAt,
		"updatedAt":   item.UpdatedAt,
	})
}

// CreateCategory godoc
// @Summary Create category
// @Tags Admin Categories
// @Accept json
// @Produce json
// @Param body body CreateCategoryRequest true "Category"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /admin/categories [post]
func (h *CategoriesHandler) Create(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	item, err := h.svc.Create(c.Request.Context(), categories.CreateInput{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		switch {
		case errors.Is(err, categories.ErrInvalidName):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          item.ID,
		"name":        item.Name,
		"description": item.Description,
		"createdAt":   item.CreatedAt,
		"updatedAt":   item.UpdatedAt,
	})
}

// UpdateCategory godoc
// @Summary Update category
// @Tags Admin Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param body body UpdateCategoryRequest true "Patch"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/categories/{id} [put]
func (h *CategoriesHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	item, err := h.svc.Update(c.Request.Context(), id, categories.UpdateInput{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		switch {
		case errors.Is(err, categories.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		case errors.Is(err, categories.ErrInvalidName):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
		case errors.Is(err, categories.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          item.ID,
		"name":        item.Name,
		"description": item.Description,
		"createdAt":   item.CreatedAt,
		"updatedAt":   item.UpdatedAt,
	})
}

// DeleteCategory godoc
// @Summary Delete category
// @Tags Admin Categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /admin/categories/{id} [delete]
func (h *CategoriesHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.svc.Delete(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, categories.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		case errors.Is(err, categories.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		case errors.Is(err, categories.ErrHasProducts):
			c.JSON(http.StatusConflict, gin.H{"error": "category has products"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
