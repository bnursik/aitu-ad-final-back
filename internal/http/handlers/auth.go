package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/bnursik/aitu-ad-final-back/internal/domain/users"
	"github.com/bnursik/aitu-ad-final-back/internal/http/middleware"
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

type UpdateProfileRequest struct {
	Name    *string `json:"name"`
	Address *string `json:"address"`
	Phone   *string `json:"phone"`
	Bio     *string `json:"bio"`
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
		log.Printf("login error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token, "user": user})
}

// Register godoc
// @Summary Register new admin
// @Description Create new account and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "Register data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/auth/register [post]
func (h *AuthHandler) AdminRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	token, user, err := h.svc.AdminRegister(c.Request.Context(), req.Name, req.Email, req.Password)
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

// GetProfile godoc
// @Summary Get user profile (auth required)
// @Tags Profile
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
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

	user, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile godoc
// @Summary Update user profile (auth required)
// @Tags Profile
// @Accept json
// @Produce json
// @Param body body UpdateProfileRequest true "Profile data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
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

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	user, err := h.svc.UpdateProfile(c.Request.Context(), userID, users.UpdateProfileInput{
		Name:    req.Name,
		Address: req.Address,
		Phone:   req.Phone,
		Bio:     req.Bio,
	})
	if err != nil {
		switch {
		case errors.Is(err, users.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case err.Error() == "invalid name":
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
		case err.Error() == "invalid bio":
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bio"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetAllUsers godoc
// @Summary Get all users (admin only)
// @Tags Admin
// @Produce json
// @Success 200 {array} users.PublicUser
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/users [get]
func (h *AuthHandler) GetAllUsers(c *gin.Context) {
	users, err := h.svc.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, users)
}
