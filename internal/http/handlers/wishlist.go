package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/wishlist"
	"github.com/bnursik/aitu-ad-final-back/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

type WishlistHandler struct {
	svc wishlist.Service
}

func NewWishlistHandler(svc wishlist.Service) *WishlistHandler {
	return &WishlistHandler{svc: svc}
}

type AddToWishlistRequest struct {
	ProductID string `json:"product_id" binding:"required"`
}

// AddToWishlist godoc
// @Summary Add product to wishlist (auth required)
// @Tags Wishlist
// @Accept json
// @Produce json
// @Param body body AddToWishlistRequest true "Product ID"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /wishlist [post]
func (h *WishlistHandler) Add(c *gin.Context) {
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

	var req AddToWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	item, err := h.svc.Add(c.Request.Context(), userID, wishlist.AddItemInput{
		ProductID: req.ProductID,
	})
	if err != nil {
		switch {
		case errors.Is(err, wishlist.ErrInvalidProduct):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		case errors.Is(err, wishlist.ErrAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "product already in wishlist"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        item.ID,
		"productId": item.ProductID,
		"createdAt": item.CreatedAt,
	})
}

// ListWishlist godoc
// @Summary Get user's wishlist (auth required)
// @Tags Wishlist
// @Produce json
// @Param offset query int true "Offset for pagination"
// @Param limit query int true "Limit for pagination"
// @Success 200 {array} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /wishlist [get]
func (h *WishlistHandler) List(c *gin.Context) {
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

	items, err := h.svc.List(c.Request.Context(), userID, wishlist.ListFilter{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	out := make([]gin.H, 0, len(items))
	for _, item := range items {
		out = append(out, gin.H{
			"id":        item.ID,
			"productId": item.ProductID,
			"createdAt": item.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, out)
}

// DeleteFromWishlist godoc
// @Summary Remove product from wishlist (auth required)
// @Tags Wishlist
// @Produce json
// @Param id path string true "Wishlist Item ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /wishlist/{id} [delete]
func (h *WishlistHandler) Delete(c *gin.Context) {
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

	id := c.Param("id")

	err := h.svc.Delete(c.Request.Context(), userID, id)
	if err != nil {
		switch {
		case errors.Is(err, wishlist.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		case errors.Is(err, wishlist.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
