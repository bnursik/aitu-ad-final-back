package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/orders"
	"github.com/bnursik/aitu-ad-final-back/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

type OrdersHandler struct {
	svc orders.Service
}

func NewOrdersHandler(svc orders.Service) *OrdersHandler {
	return &OrdersHandler{svc: svc}
}

type CreateOrderRequest struct {
	Items []struct {
		ProductID string `json:"productId" binding:"required"`
		Quantity  int64  `json:"quantity" binding:"required"`
	} `json:"items" binding:"required"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type FindOrderByIDRequest struct {
	OrderID string `json:"order_id" binding:"required"`
}

func isAdminFromCtx(c *gin.Context) bool {
	roleVal, _ := c.Get(middleware.CtxRole)
	role, _ := roleVal.(string)
	return role == "admin"
}

func userIDFromCtx(c *gin.Context) (string, bool) {
	v, ok := c.Get(middleware.CtxUserID)
	if !ok {
		return "", false
	}
	s, _ := v.(string)
	return s, s != ""
}

// ListOrders godoc
// @Summary List orders (user: own, admin: all)
// @Tags Orders
// @Produce json
// @Param offset query int true "Offset for pagination"
// @Param limit query int true "Limit for pagination"
// @Success 200 {array} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /orders [get]
func (h *OrdersHandler) List(c *gin.Context) {
	uid, ok := userIDFromCtx(c)
	if !ok {
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

	admin := isAdminFromCtx(c)

	items, total, err := h.svc.List(c.Request.Context(), uid, admin, orders.ListFilter{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	out := make([]gin.H, 0, len(items))
	for _, it := range items {
		out = append(out, orderToJSON(it, admin))
	}

	c.JSON(http.StatusOK, gin.H{
		"items":  out,
		"total":  total,
		"offset": offset,
		"limit":  limit,
	})
}

// GetOrder godoc
// @Summary Get order by ID (user: own, admin: any)
// @Tags Orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /orders/{id} [get]
func (h *OrdersHandler) Get(c *gin.Context) {
	uid, ok := userIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	admin := isAdminFromCtx(c)

	id := c.Param("id")
	it, err := h.svc.Get(c.Request.Context(), id, uid, admin)
	if err != nil {
		switch {
		case errors.Is(err, orders.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		case errors.Is(err, orders.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, orderToJSON(it, admin))
}

// CreateOrder godoc
// @Summary Create order (auth required)
// @Tags Orders
// @Accept json
// @Produce json
// @Param body body CreateOrderRequest true "Order"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /orders [post]
func (h *OrdersHandler) Create(c *gin.Context) {
	uid, ok := userIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	in := orders.CreateInput{Items: make([]orders.Item, 0, len(req.Items))}
	for _, it := range req.Items {
		in.Items = append(in.Items, orders.Item{
			ProductID: it.ProductID,
			Quantity:  it.Quantity,
		})
	}

	created, err := h.svc.Create(c.Request.Context(), uid, in)
	if err != nil {
		switch {
		case errors.Is(err, orders.ErrInvalidItems):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid items"})
		case errors.Is(err, orders.ErrInvalidProduct):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid productId"})
		case errors.Is(err, orders.ErrInvalidQty):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quantity"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusCreated, orderToJSON(created, false))
}

// UpdateOrderStatus godoc
// @Summary Update order status (admin only)
// @Tags Admin Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param body body UpdateOrderStatusRequest true "Status"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/orders/{id}/status [put]
func (h *OrdersHandler) UpdateStatus(c *gin.Context) {
	// admin guard лучше делать в routes (AdminOnly), но на всякий:
	if !isAdminFromCtx(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
		return
	}

	id := c.Param("id")

	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	updated, err := h.svc.UpdateStatus(c.Request.Context(), id, orders.Status(req.Status))
	if err != nil {
		switch {
		case errors.Is(err, orders.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		case errors.Is(err, orders.ErrInvalidStatus):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		case errors.Is(err, orders.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, orderToJSON(updated, true))
}

// FindOrderByID godoc
// @Summary Find order by ID (admin only)
// @Tags Admin Orders
// @Accept json
// @Produce json
// @Param body body FindOrderByIDRequest true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/orders/find [post]
func (h *OrdersHandler) FindOrderByID(c *gin.Context) {
	if !isAdminFromCtx(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
		return
	}

	var req FindOrderByIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	order, err := h.svc.Get(c.Request.Context(), req.OrderID, "", true)
	if err != nil {
		switch {
		case errors.Is(err, orders.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order_id"})
		case errors.Is(err, orders.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, orderToJSON(order, true))
}

func orderToJSON(o orders.Order, admin bool) gin.H {
	items := make([]gin.H, 0, len(o.Items))
	for _, it := range o.Items {
		items = append(items, gin.H{
			"productId": it.ProductID,
			"quantity":  it.Quantity,
			"unitPrice": it.UnitPrice,
			"lineTotal": it.LineTotal,
		})
	}

	out := gin.H{
		"id":         o.ID,
		"items":      items,
		"status":     o.Status,
		"totalPrice": o.TotalPrice,
		"createdAt":  o.CreatedAt,
		"updatedAt":  o.UpdatedAt,
	}
	if admin {
		out["userId"] = o.UserID
	}
	return out
}
