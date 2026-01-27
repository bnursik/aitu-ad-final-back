package handlers

import (
	"errors"
	"net/http"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/users"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc users.Service
}

func NewAuthHandler(svc users.Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=60"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register godoc
// @Summary Register new user
// @Description Create new account and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "Register data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	token, user, err := h.svc.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrEmailTaken):
			c.JSON(http.StatusConflict, gin.H{"error": "email already taken"})
		case errors.Is(err, users.ErrInvalidEmail), errors.Is(err, users.ErrInvalidPassword):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case err.Error() == "invalid name":
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"access_token": token, "user": user})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	token, user, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, users.ErrInvalidCreds) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token, "user": user})
}
